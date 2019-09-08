package models

import "net/http"

type Hoster interface {
	downloadCookies() []*http.Cookie
	supportsUrl(url string) bool
	getDirectLink(file *File) (string, error)
	getApiInfo(file *File) (online bool, filename string, checksum string, checksumType string, size uint64, metaInfo map[string]string)
}
