package models

import (
"log"
"strconv"
"regexp"
"strings"
"goload/server/unrar"

	"time"
)

type Package struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Finished  bool `json:"finished"`
	Files     []*File  `json:"files"`
	Password  string `json:"password"`
	Progress  float64 `json:"progress"`
	DateAdded time.Time `json:"date_added"`
	Size uint64 `json:"size"`
}

type File struct {
	Url string `json:"url"`
	Finished bool `json:"finished"`
	Progress float64 `json:"progress"`
	filePath string
	Filename string `json:"filename"`
	Size uint64 `json:"size"`
	DownloadSpeed string `json:"download_speed"`
	Failed bool `json:"failed"`
	ETE time.Duration `json:"ete"`
	Error error
}


func (pack *Package) Download(downloader *Uploaded) {
	log.Println("Downloading package" + pack.Name + " with "+ strconv.Itoa(len(pack.Files)) + " files")
	pack.Finished = false
	downloader.DownloadPackage(pack)
	pack.Finished = true
}

func (pack *Package) UpdateProgress() {
	progess := 0.0
	for _,file:= range pack.Files{
		progess+= file.Progress
	}
	pack.Progress= float64(progess)/float64(len(pack.Files))
}
func (pack *Package) UpdateSize() {
	var size uint64 = 0
	for _,file:= range pack.Files{
		size+= file.Size
	}
	pack.Size= size
}

func (pack *Package) Update() {
	pack.UpdateProgress()
	pack.UpdateSize()

}

func (pack *Package) Unrar(path string) {
	r, _ := regexp.Compile(`.*part0*1\.rar`)
	for _,file :=range pack.Files{
		if file.filePath == ""{
			continue;
		} 
		if r.MatchString(file.filePath) || !strings.Contains(file.filePath,`part`) {
			unrar.Unrar(file.filePath,path+pack.Name+"/",pack.Password)
		}
	}
}