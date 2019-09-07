package models

type Hoster interface {
	supportsUrl(url string) bool
	getDirectLink(file *File) (string, error)
	getApiInfo(file *File) (online bool, filename string, checksum string, checksumType string, size uint64)
}
