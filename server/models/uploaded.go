package models

import (
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"os"
	"strconv"
	"bytes"
	"strings"
	"time"
	"github.com/pivotal-golang/bytefmt"
	"errors"
	"github.com/cavaliercoder/grab"
	"goload/server/models/configuration"
)

type Uploaded struct {
	config *configuration.Configuration
	loginCookie *http.Cookie
}

type WriteCounter struct {
	FileSize uint64 // Total # of bytes transferred
	Total    uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	//log.Println("Progess: "+ strconv.FormatFloat(float64(wc.Total)/float64(wc.FileSize)*100.0,'f', -1, 64))
	return n, nil
}

func NewUploaded(config *configuration.Configuration) *Uploaded {
	ul := &Uploaded{config:config}
	ul.login()
	return ul
}

func (ul *Uploaded) DownloadPackage(pack *Package) (error) {
	savePath := ul.config.Dirs.DownloadDir + pack.Name
	mkdirErr := os.MkdirAll(savePath, 0755)
	if mkdirErr != nil {
		return errors.New("Error creating directory " + savePath)
	}
	for _, file := range pack.Files {
		link, error := ul.getDirectLink(file)
		if (error != nil) {
			file.Failed = true
			file.Error = error
			continue
		}
		fileName, size := getHeadParams(link)
		file.Filename = fileName
		file.Size = size
	}
	pack.UpdateSize()
	BATCH_SIZE := 1
	i := 0
	for i < len(pack.Files) {
		requests := make([]*grab.Request,0)
		requestMap :=make(map[*grab.Request]*File)
		for b:=0; b < BATCH_SIZE && i+b <len(pack.Files); b++ {
			file := pack.Files[i+b]
			link, error := ul.getDirectLink(file)
			if (error != nil) {
				log.Println(error)
				file.Failed = true
				file.Progress = 100.0
				file.Error = error
				continue
			}
			req, requestError := grab.NewRequest(link)
			if (requestError != nil) {
				file.Failed = true
				file.Progress = 100.0
				file.Error = requestError
				continue
			}
			req.BufferSize = 4096*1024
			req.Size = file.Size
			req.Filename = savePath+ "/" + file.Filename
			requestMap[req] = file
			requests = append(requests,req)
		}
		i+= BATCH_SIZE
		ul.downloadBatch(BATCH_SIZE,requests,requestMap,pack);

	}
	return nil
}

func (ul *Uploaded) downloadBatch(batchSize int, requests []*grab.Request,requestMap map[*grab.Request]*File, pack *Package){
	grabClient:= grab.NewClient()
	t := time.NewTicker(200 * time.Millisecond)
	respch := grabClient.DoBatch(batchSize, requests...)
	completed := 0
	responses := make([]*grab.Response, 0)
	for completed < len(requests) {
		select {
		case resp := <-respch:
			if resp != nil {
				responses = append(responses, resp)
				log.Printf("Started downloading %s %d / %d bytes (%d%%)\n", resp.Filename, resp.BytesTransferred(), resp.Size, int(100*resp.Progress()))

			}
		case <-t.C:
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					// print final result
					if resp.Error != nil {
						log.Printf("Error downloading %s: %v\n", resp.Filename, resp.Error)
						requestMap[resp.Request].Failed = true
						requestMap[resp.Request].Progress = 100.0
					} else {
						log.Printf("Finished %s %d / %d bytes (%d%%)\n", resp.Filename, resp.BytesTransferred(), resp.Size, int(100*resp.Progress()))
						requestMap[resp.Request].Finished = true
						requestMap[resp.Request].Progress = 100*resp.Progress()

						requestMap[resp.Request].Failed = false
						requestMap[resp.Request].filePath = resp.Filename
						log.Println("Average speed: " +bytefmt.ByteSize(uint64(resp.AverageBytesPerSecond())) + "/s")

					}

					// mark completed
					responses[i] = nil
					completed++
					pack.Update()
				}
			}
			for _, resp := range responses {
				if resp != nil {
					//TODO Speed + ETA
					requestMap[resp.Request].DownloadSpeed = bytefmt.ByteSize(uint64(resp.AverageBytesPerSecond())) + "/s"
					requestMap[resp.Request].Progress = 100*resp.Progress()
					requestMap[resp.Request].ETE = resp.ETA().Sub(time.Now())
					pack.UpdateProgress()
				}
			}
		}
	}
	pack.Finished = true
	t.Stop()
}

func (ul *Uploaded)getDirectLink(file *File) (string, error) {
	if (ul.loginCookie == nil || ul.loginCookie.Expires.Before(time.Now())) {
		log.Println("Cookie expired, logging in")
		err := ul.login()
		if (err != nil) {
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
	dddata, _ := ioutil.ReadAll(htmlResp.Body);
	htmlString := string(dddata)
	link, linkError := extractDirectLink(htmlString)
	if linkError != nil {
		return "", errors.New("File " + file.Url + " not found")
	}
	return link, nil

}

func getHeadParams(link string) (string, uint64) {
	header, _ := http.Head(link)
	length, _ := strconv.Atoi(header.Header.Get("Content-Length"))
	fileNameBase := header.Header.Get("Content-Disposition")
	searchString := `attachment; filename="`
	fileName := fileNameBase[strings.Index(fileNameBase, searchString) + len(searchString):len(fileNameBase) - 1]
	return fileName, uint64(length)
}

func extractDirectLink(htmlString string) (string, error) {
	find := `<form method="post" action="`
	index := strings.Index(htmlString, `<form method="post" action="`)
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
	r, _ := http.NewRequest("POST", "http://uploaded.net/io/login", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	login, _ := client.Do(r)
	defer login.Body.Close()
	for _, element := range login.Cookies() {
		if (element.Name == "login") {
			ul.loginCookie = element
			return nil
		}

	}
	return errors.New("Login failed")
}


