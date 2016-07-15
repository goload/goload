package hoster

import (
	"github.com/cavaliercoder/grab"
	"regexp"
)

type LinkInfo struct {
	Online bool
	Filename string
	Checksum []byte
	Size uint64
}

type Hoster interface {
	LinkInfo(link string) LinkInfo
	PremiumRequest(link string) grab.Request
	LinkPattern() regexp.Regexp
}