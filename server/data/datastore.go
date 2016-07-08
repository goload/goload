package data

import (
	"goload/server/models"
)

type Datastore struct {
	packages map[string]*models.Package
}

func NewDatastore() *Datastore {
	return &Datastore{packages:make(map[string]*models.Package)}
}

func (ds *Datastore) AddPackage(pack *models.Package) {
	ds.packages[pack.Id] = pack
	
}

func (ds *Datastore) RemovePackage(id string) bool {
	_, exists := ds.packages[id]
	delete(ds.packages, id)
	return exists

}

func (ds *Datastore) GetPackages() []*models.Package {
	values := make([]*models.Package, 0)
	for _, value := range ds.packages {
		values = append(values, value)
	}
	return values
}

func (ds *Datastore) GetPackage(id string) (pack *models.Package,exists bool) {
	pack,exists = ds.packages[id]
	return
}