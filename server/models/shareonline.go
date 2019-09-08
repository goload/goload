package models

import (
	"errors"
	"fmt"
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

type ShareOnline struct {
	config       *configuration.Configuration
	apiUrl       string
	urlPattern   string
	loginUrl     string
	loginCookies []*http.Cookie
}

func NewShareOnline(config *configuration.Configuration) *ShareOnline {
	so := &ShareOnline{config: config}
	so.loginCookies = make([]*http.Cookie, 0)
	so.apiUrl = "http://api.share-online.biz/linkcheck.php"
	so.urlPattern = `https?://(?:www\.)?(share-online\.biz|egoshare\.com)/(download\.php\?id=|dl/)(?P<ID>\w+)`
	so.loginUrl = "https://api.share-online.biz/cgi-bin"
	return so
}

func (so *ShareOnline) getDirectLink(file *File) (string, error) {
	if len(so.loginCookies) == 0 || so.loginCookies[0].Expires.Before(time.Now()) {
		log.Println("Cookie expired, logging in")
		err := so.login()
		if err != nil {
			return "", err
		}
	}
	resp, err := http.Get(fmt.Sprintf("https://api.share-online.biz/cgi-bin?q=linkdata&lid=%s&username=%s&password=%s", file.MetaInfo["fileid"], so.config.Account.Username, so.config.Account.Password))
	body, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(body)
	log.Println(stringBody)
	for _, line := range strings.Split(strings.TrimSuffix(stringBody, "\n"), "\n") {
		split := strings.Split(line, ":")
		if split[0] == "URL" {
			return strings.ReplaceAll(strings.TrimSpace(line[4:]), "http://", "https://"), nil
		}
	}
	return "", err

}

func (so *ShareOnline) supportsUrl(url string) bool {
	re := regexp.MustCompile(so.urlPattern)
	matches := re.FindAllStringSubmatch(url, -1)
	return matches != nil
}

func (so *ShareOnline) getApiInfo(file *File) (online bool, filename string, checksum string, checksumType string, size uint64, metaInfo map[string]string) {
	re := regexp.MustCompile(so.urlPattern)
	n1 := re.SubexpNames()
	matches := re.FindAllStringSubmatch(file.Url, -1)
	if matches == nil {
		return false, "", "", "", 0, nil
	}
	r2 := matches[0]
	md := map[string]string{}
	for i, n := range r2 {
		md[n1[i]] = n
	}
	resp, err := http.Get(fmt.Sprintf("%s?md5=1&links=%s", so.apiUrl, md["ID"]))
	defer resp.Body.Close()
	if err != nil {
		return false, "", "", "", 0, nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(body)
	log.Println(stringBody)
	results := strings.Split(stringBody, ";")
	if results[1] != "OK" {
		return
	}
	online = true
	fileSize, _ := strconv.Atoi(strings.TrimSpace(results[3]))
	size = uint64(fileSize)
	log.Println(fileSize)
	filename = results[2]
	log.Println(filename)
	checksum = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(results[4])), "\n\n", "")
	log.Println(checksum)
	metaInfo = make(map[string]string)
	metaInfo["fileid"] = results[0]
	log.Println(results[0])
	checksumType = "md5"
	return
}

func (so *ShareOnline) downloadCookies() []*http.Cookie {
	return so.loginCookies
}

func (so *ShareOnline) extractDirectLink(htmlString string) (string, error) {
	find := `<form method="post" action="`
	index := strings.Index(htmlString, find)
	if index == -1 {
		return "", errors.New("File link not found")
	}
	substring := htmlString[(index + len(find)):len(htmlString)]
	quoteIndex := strings.Index(substring, `"`)
	return substring[:quoteIndex], nil
}

func (so *ShareOnline) login() error {
	data := url.Values{}
	data.Set("username", so.config.Account.Username)
	data.Add("password", so.config.Account.Password)
	data.Add("q", "userdetails")
	data.Add("aux", "traffic")
	login, _ := http.Get(fmt.Sprintf("%s?q=userdetails&aux=traffic&username=%s&password=%s", so.loginUrl, so.config.Account.Username, so.config.Account.Password))
	defer login.Body.Close()
	body, _ := ioutil.ReadAll(login.Body)
	stringBody := string(body)
	cookieVal := ""
	expireTime := time.Now()
	for _, line := range strings.Split(strings.TrimSuffix(stringBody, "\n"), "\n") {
		split := strings.Split(line, "=")
		if split[0] == "a" {
			cookieVal = split[1]
		}
		if split[0] == "expire_date" {
			i, _ := strconv.ParseInt(split[1], 10, 64)
			expireTime = time.Unix(i, 0)
		}

	}
	if cookieVal != "" {
		so.loginCookies = make([]*http.Cookie, 0)
		so.loginCookies = append(so.loginCookies, &http.Cookie{Name: "a", Value: cookieVal, Expires: expireTime})
		so.loginCookies = append(so.loginCookies, &http.Cookie{Name: "version", Value: "v4", Expires: expireTime})
		return nil
	}
	return errors.New("Login failed")
}
