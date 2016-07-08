package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"goload/server/models/configuration"
)

type ConfigurationController struct {
	config *configuration.Configuration
}

func NewConfiguartionController(config *configuration.Configuration) *ConfigurationController {
	return &ConfigurationController{config:config}
}

// GetUser retrieves an individual user resource
func (cc *ConfigurationController) GetConfiguration(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var config configuration.Configuration  = *cc.config
	config.Account.Password = ""
	configuration, _ := json.Marshal(config)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(configuration)
}

func (cc *ConfigurationController) UpdateDirs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	dirs := configuration.Dirs{}
	error := json.NewDecoder(r.Body).Decode(&dirs)
	if(error!= nil) {
		w.WriteHeader(400)
		return
	}
	cc.config.Dirs = dirs
	cc.config.Sanitize()
	cc.config.LogConfig()
	cc.config.Save()
	w.WriteHeader(200)
}

func (cc *ConfigurationController) UpdateAccount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	account := configuration.Account{}
	error := json.NewDecoder(r.Body).Decode(&account)
	if(error!= nil || account.Password == "") {
		w.WriteHeader(400)
		return
	}
	cc.config.Account = account
	cc.config.Sanitize()
	cc.config.LogConfig()
	cc.config.Save()
	w.WriteHeader(200)
}


