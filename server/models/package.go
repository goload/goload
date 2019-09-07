package models

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"goload/server/unrar"
	"time"
)

type Package struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Finished      bool      `json:"finished"`
	Files         []*File   `json:"files"`
	DLC           string    `json:"dlc"`
	Password      string    `json:"password"`
	Progress      float64   `json:"progress"`
	UnrarProgress float64   `json:"unrar_progress"`
	Extracting    bool      `json:"extracting"`
	DateAdded     time.Time `json:"date_added"`
	Size          uint64    `json:"size"`
}

type File struct {
	Url           string  `json:"url"`
	Finished      bool    `json:"finished"`
	Offline       bool    `json:"offline"`
	Checksum      string  `json:"-"`
	ChecksumType  string  `json:"-"`
	Progress      float64 `json:"progress"`
	UnrarProgress float64 `json:"unrar_progress"`
	Extracting    bool    `json:"extracting"`
	filePath      string
	Filename      string        `json:"filename"`
	Size          uint64        `json:"size"`
	DownloadSpeed string        `json:"download_speed"`
	Failed        bool          `json:"failed"`
	ETE           time.Duration `json:"ete"`
	Error         error
}

func (pack *Package) Download(downloader *Downloader) {
	log.Println("Downloading package " + pack.Name + " with " + strconv.Itoa(len(pack.Files)) + " files")
	pack.Finished = false
	downloader.DownloadPackage(pack)
	pack.Finished = true
}

func (pack *Package) UpdateProgress() {
	progess := 0.0
	for _, file := range pack.Files {
		progess += file.Progress
	}
	pack.Progress = float64(progess) / float64(len(pack.Files))
}
func (pack *Package) UpdateSize() {
	var size uint64 = 0
	for _, file := range pack.Files {
		size += file.Size
	}
	pack.Size = size
}

func (pack *Package) updateUnrarProgress(files []*File) {
	unrarProgress := float64(0.0)
	for _, file := range files {
		unrarProgress += file.UnrarProgress
	}
	pack.UnrarProgress = unrarProgress / float64(len(files))
}

func (pack *Package) Update() {
	pack.UpdateProgress()
	pack.UpdateSize()

}

func (pack *Package) Retry() {
	pack.Finished = false
	//TODO only temporary solution

}

func (pack *Package) Unrar(path string) {
	pack.Extracting = true
	unrarFiles := pack.getExtractableFiles()
	for _, file := range unrarFiles {
		c := unrar.Unrar(file.filePath, path+pack.Name+"/", pack.Password)
		file.Extracting = true
		for i := range c {
			if i.Error != nil {
				log.Println(i.Error)
				continue
			}
			file.UnrarProgress = i.Progess
			pack.updateUnrarProgress(unrarFiles)
		}
		log.Println("Extrated " + file.filePath + " successfully")
		file.Extracting = false
	}
	pack.Extracting = false
}

func (pack *Package) getExtractableFiles() []*File {
	r := regexp.MustCompile(`.*part0*1\.rar`)
	var files = make([]*File, 0)
	for _, file := range pack.Files {
		if file.filePath == "" {
			continue
		}
		if r.MatchString(file.Filename) || !strings.Contains(file.Filename, `part`) {
			files = append(files, file)
		}
	}
	return files
}
