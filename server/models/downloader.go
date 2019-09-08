package models

import (
	"encoding/hex"
	"errors"
	"goload/server/models/configuration"
	"log"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/pivotal-golang/bytefmt"
)

type Downloader struct {
	hoster []Hoster
	config *configuration.Configuration
}

func NewDownloader(hoster []Hoster, config *configuration.Configuration) *Downloader {
	dl := &Downloader{hoster: hoster, config: config}
	return dl
}

func (dl *Downloader) DownloadPackage(pack *Package) error {
	savePath := dl.config.Dirs.DownloadDir + pack.Name
	mkdirErr := os.MkdirAll(savePath, 0755)
	if mkdirErr != nil {
		return errors.New("Error creating directory " + savePath)
	}
	for _, file := range pack.Files {
		hoster := dl.getHoster(file.Url)
		online, fileName, checksum, checksumType, size, metaInfo := hoster.getApiInfo(file)
		file.Filename = fileName
		file.Offline = !online
		file.Checksum = checksum
		file.ChecksumType = checksumType
		file.Size = size
		file.MetaInfo = metaInfo
	}
	pack.UpdateSize()
	BATCH_SIZE := 1
	i := 0
	for i < len(pack.Files) {
		requests := make([]*grab.Request, 0)
		requestMap := make(map[*grab.Request]*File)
		for b := 0; b < BATCH_SIZE && i+b < len(pack.Files); b++ {
			file := pack.Files[i+b]
			if file.Offline {
				file.Error = errors.New("Offline")
				log.Println("Offline: " + file.Url)
				file.Failed = true
				file.Progress = 100.0
				continue
			}
			hoster := dl.getHoster(file.Url)
			link, error := hoster.getDirectLink(file)
			if error != nil {
				log.Println("Get Direct Link failed " + error.Error())
				file.Failed = true
				file.Progress = 100.0
				file.Error = error
				continue
			}
			log.Println(link)
			req, requestError := grab.NewRequest(link)
			for _, c := range hoster.downloadCookies() {
				req.HTTPRequest.AddCookie(c)
			}
			if requestError != nil {
				file.Failed = true
				file.Progress = 100.0
				file.Error = requestError
				continue
			}
			b, _ := hex.DecodeString(file.Checksum)
			req.SetChecksum(file.ChecksumType, b)
			req.RemoveOnError = true
			req.BufferSize = 4096 * 1024
			req.Size = file.Size
			req.Filename = savePath + "/" + file.Filename
			requestMap[req] = file
			requests = append(requests, req)
		}
		i += BATCH_SIZE
		dl.downloadBatch(BATCH_SIZE, requests, requestMap, pack)

	}
	return nil
}

func (dl *Downloader) downloadBatch(batchSize int, requests []*grab.Request, requestMap map[*grab.Request]*File, pack *Package) {
	grabClient := grab.NewClient()
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
						requestMap[resp.Request].Progress = 100 * resp.Progress()
						requestMap[resp.Request].Failed = false
						requestMap[resp.Request].filePath = resp.Filename
						requestMap[resp.Request].ETE = 0
						log.Println("Average speed: " + bytefmt.ByteSize(uint64(resp.AverageBytesPerSecond())) + "/s")

					}

					// mark completed
					responses[i] = nil
					completed++
					pack.Update()
				}
			}
			for _, resp := range responses {
				if resp != nil {
					requestMap[resp.Request].DownloadSpeed = bytefmt.ByteSize(uint64(resp.AverageBytesPerSecond())) + "/s"
					requestMap[resp.Request].Progress = 100 * resp.Progress()
					requestMap[resp.Request].ETE = resp.ETA().Sub(time.Now())
					pack.UpdateProgress()
				}
			}
		}
	}
	pack.Finished = true
	t.Stop()
}

func (dl *Downloader) getHoster(url string) Hoster {
	for _, host := range dl.hoster {
		if host.supportsUrl(url) {
			return host
		}
	}
	return nil
}
