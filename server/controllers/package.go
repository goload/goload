package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"goload/server/models"
	"github.com/nu7hatch/gouuid"
	"goload/server/data"
	"time"
)

type PackageController struct {
	database *data.Datastore

}

func NewPackageController(database *data.Datastore) *PackageController {
	return &PackageController{database:database}
}

// GetUser retrieves an individual user resource
func (uc PackageController) ListPackages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	packagesJson, _ := json.Marshal(uc.database.GetPackages())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(packagesJson)

}

func (uc PackageController) CreatePackage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	pack := &models.Package{}
	error := json.NewDecoder(r.Body).Decode(pack)
	if error != nil || pack.Name == "" {
		w.WriteHeader(400)
		return
	}
	u4, _ := uuid.NewV4()
	pack.Id = u4.String()
	pack.DateAdded = time.Now()
	uc.database.AddPackage(pack)
	pjson, _ := json.Marshal(pack)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(pjson)
}

func (uc PackageController) RemovePackage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	id := p.ByName("id")
	success := uc.database.RemovePackage(id)
	if !success {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
}

func (uc PackageController) RetryPackage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	id := p.ByName("id")
	pack,exists := uc.database.GetPackage(id)
	if !exists || !pack.Finished {
		w.WriteHeader(404)
		return
	}
	pack.Retry()
	w.WriteHeader(200)
}


