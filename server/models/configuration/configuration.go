package configuration

import (
	"encoding/json"
	"os"
	"errors"
	"log"
	"io"
)

type Configuration struct {
	Dirs    Dirs `json:"dirs"`
	Account Account `json:"account"`
	filename string
}

type Dirs struct {
	DownloadDir string `json:"downloadDir"`
	ExtractDir  string `json:"extractDir"`
}
type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewConfigurationFromFileName(filename string) (*Configuration, error) {
	file, error := os.Open(filename)
	if (error != nil) {
		return nil, errors.New("Failed to open config file")
	}
	defer file.Close()
	return newConfigurationFromFile(file,filename)
}

func newConfigurationFromFile(file io.Reader,filename string)  (*Configuration, error) {
	decoder := json.NewDecoder(file)
	config := Configuration{}
	error := decoder.Decode(&config)
	if (error != nil) {
		return nil, errors.New("Could not read config, maybe malformed?")
	}
	if config.Dirs.DownloadDir == "" {
		return nil, errors.New("downloadsDir not found in config")
	}
	if config.Dirs.ExtractDir == "" {
		return nil, errors.New("extractedDir not found in config")
	}
	if config.Account.Username == "" {
		return nil, errors.New("username not found in config")
	}
	if config.Account.Password == "" {
		return nil, errors.New("password not found in config")
	}
	config.filename = filename
	config.Sanitize()
	config.LogConfig()
	return &config, nil

}

func (config *Configuration ) Sanitize()  {
		config.Dirs.DownloadDir = sanitizePath(config.Dirs.DownloadDir)
		config.Dirs.ExtractDir = sanitizePath(config.Dirs.ExtractDir)
}
func sanitizePath(path string) string {
	if(path[len(path)-1] != '/') {
		return path+"/"
	}
	return path
}

func (config *Configuration ) LogConfig()  {
	log.Println("Config:")
	log.Println("Downloading to " + config.Dirs.DownloadDir + ", extracting to " + config.Dirs.ExtractDir)
	log.Println("Using account " + config.Account.Username)
}

func (config *Configuration) Save() {
	f, err := os.Create(config.filename)
	if err != nil {
		log.Println("Saving config failed")
	}
	defer f.Close()
	data, _ := json.MarshalIndent(config,"","   ")
	_,werror := f.Write(data)
	f.Sync()
	if(werror != nil) {
		log.Println("saved!?"  + werror.Error())
	}
}