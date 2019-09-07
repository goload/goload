package models

import (
	"bytes"
	"errors"
	"goload/server/models/configuration"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const API_KEY string = "lhF2IeeprweDfu9ccWlxXVVypA5nA3EL"
const API_URL string = "http://uploaded.net/api/filemultiple"
const URL_PATTERN = `https?://(?:www\.)?(uploaded\.(to|net)|ul\.to)(/file/|/?\?id=|.*?&id=|/)(?P<ID>\w+)`
const LOGIN_URL = "http://uploaded.net/io/login"

type Uploaded struct {
	config      *configuration.Configuration
	loginCookie *http.Cookie
}

func NewUploaded(config *configuration.Configuration) *Uploaded {
	ul := &Uploaded{config: config}
	ul.login()
	return ul
}

func (ul *Uploaded) getDirectLink(file *File) (string, error) {
	if ul.loginCookie == nil || ul.loginCookie.Expires.Before(time.Now()) {
		log.Println("Cookie expired, logging in")
		err := ul.login()
		if err != nil {
			return "", err
		}
	}
	htmlRequest, htmlRequestErr := http.NewRequest("GET", file.Url, nil)
	if htmlRequestErr != nil {
		log.Println("htmlRequestErr")
	}
	htmlRequest.AddCookie(ul.loginCookie)
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		if len(via) > 0 {
			req.AddCookie(ul.loginCookie)
		}
		return nil
	}}
	htmlResp, htmlErr := client.Do(htmlRequest)
	if htmlErr != nil {
		return "", errors.New("File " + file.Url + " not found")
	}
	defer htmlResp.Body.Close()
	dddata, _ := ioutil.ReadAll(htmlResp.Body)
	htmlString := string(dddata)
	link, linkError := extractDirectLink(htmlString)
	if linkError != nil {
		log.Println(htmlString)
		return "", errors.New("Link " + file.Url + " not found")
	}
	return link, nil

}

func (ul *Uploaded) supportsUrl(url string) bool {
	re := regexp.MustCompile(URL_PATTERN)
	matches := re.FindAllStringSubmatch(url, -1)
	return matches != nil
}

func (ul *Uploaded) getApiInfo(file *File) (online bool, filename string, checksum string, checksumType string, size uint64) {
	re := regexp.MustCompile(URL_PATTERN)
	n1 := re.SubexpNames()
	matches := re.FindAllStringSubmatch(file.Url, -1)
	if matches == nil {
		return false, "", "", "", 0
	}
	r2 := matches[0]
	md := map[string]string{}
	for i, n := range r2 {
		md[n1[i]] = n
	}
	resp, err := http.PostForm(API_URL,
		url.Values{"apikey": {API_KEY}, "id_0": {md["ID"]}})
	defer resp.Body.Close()
	if err != nil {
		return false, "", "", "", 0
	}
	body, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(body)
	results := strings.Split(stringBody, ",")
	if results[0] != "online" {
		return
	}
	online = true
	fileSize, _ := strconv.Atoi(results[2])
	size = uint64(fileSize)
	filename = results[4][:len(results[4])-1]
	checksum = results[3]
	checksumType = "sha1"
	return
}

func extractDirectLink(htmlString string) (string, error) {
	find := `<form method="post" action="`
	index := strings.Index(htmlString, find)
	if index == -1 {
		return "", errors.New("File link not found")
	}
	substring := htmlString[(index + len(find)):len(htmlString)]
	quoteIndex := strings.Index(substring, `"`)
	return substring[:quoteIndex], nil
}

func (ul *Uploaded) login() error {
	data := url.Values{}
	data.Set("id", ul.config.Account.Username)
	data.Add("pw", ul.config.Account.Password)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", LOGIN_URL, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	login, _ := client.Do(r)
	defer login.Body.Close()
	for _, element := range login.Cookies() {
		if element.Name == "login" {
			ul.loginCookie = element
			return nil
		}

	}
	return errors.New("Login failed")
}
