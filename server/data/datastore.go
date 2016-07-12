package data

import (
	"goload/server/models"
	"github.com/boltdb/bolt"
	"encoding/json"
)

type Datastore struct {
	packages map[string]*models.Package
	db *bolt.DB
}

var PACKAGES []byte = []byte("packages")

func NewDatastore(db *bolt.DB) *Datastore {
	db.Update(func(tx *bolt.Tx)error {
		tx.CreateBucketIfNotExists(PACKAGES)
		return nil
	})
	return &Datastore{db:db,packages:make(map[string]*models.Package)}
}

func (ds *Datastore) LoadData()  {
	ds.db.View(func(tx *bolt.Tx) error{
		b:= tx.Bucket(PACKAGES)
		b.ForEach(func(k,v []byte) error {
			pack := &models.Package{}
			json.Unmarshal(v,pack)
			ds.packages[pack.Id] = pack
			return nil
		})
		return nil
	})
}

func (ds *Datastore) SaveData()  {
	ds.db.Update(func(tx *bolt.Tx) error{
		b:= tx.Bucket(PACKAGES)
		for k, value := range ds.packages {
			packJson,_ := json.Marshal(value)
			b.Put([]byte(k),packJson)
		}
		return nil
	})
}

func (ds *Datastore) AddPackage(pack *models.Package) {
	ds.packages[pack.Id] = pack
	
}

func (ds *Datastore) RemovePackage(id string) bool {
	_, exists := ds.packages[id]
	delete(ds.packages, id)
	ds.db.Update(func(tx *bolt.Tx) error{
		b:= tx.Bucket(PACKAGES)
		b.Delete([]byte(id))
		return nil
	})
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