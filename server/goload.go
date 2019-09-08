package main

import (
	"goload/server/controllers"
	"goload/server/data"
	"goload/server/models"
	"goload/server/models/configuration"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

var Version = "dev"

func main() {
	config, error := configuration.NewConfigurationFromFileName("config.json")
	if error != nil {
		log.Fatal(error)
	}
	db, dberr := bolt.Open("goload_database.db", 0600, nil)
	if dberr != nil {
		log.Fatal(dberr)
	}

	router := httprouter.New()
	router.ServeFiles("/fonts/*filepath", http.Dir("./public/fonts"))
	router.ServeFiles("/public/*filepath", http.Dir("./public"))
	router.GET("/settings", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl := template.Must(template.ParseFiles("public/index.html"))

		tmpl.Execute(w, nil)
	})
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl := template.Must(template.ParseFiles("public/index.html"))

		tmpl.Execute(w, nil)
	})

	var hoster []models.Hoster
	hoster = append(hoster, models.NewUploaded(config))
	hoster = append(hoster, models.NewShareOnline(config))
	dl := models.NewDownloader(hoster, config)
	database := data.NewDatastore(db)
	database.LoadData()
	packageController := controllers.NewPackageController(database)
	router.DELETE("/api/packages/:id", packageController.RemovePackage)
	router.POST("/api/packages", packageController.CreatePackage)
	router.GET("/api/packages", packageController.ListPackages)
	router.GET("/api/packages/:id/retry", packageController.RetryPackage)
	configController := controllers.NewConfiguartionController(config)
	router.PUT("/api/config/dirs", configController.UpdateDirs)
	router.PUT("/api/config/account", configController.UpdateAccount)
	router.GET("/api/config", configController.GetConfiguration)
	router.GET("/api/version", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		versionJson := []byte(`{"version":"` + Version + `"}`)
		w.Write(versionJson)
	})
	go LoopPackages(database, dl, config)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Shutting Down")
			database.SaveData()
			os.Exit(0)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Hour * 1)
			database.SaveData()
		}
	}()
	log.Fatal(http.ListenAndServe(":3000", router))

}

func LoopPackages(database *data.Datastore, dl *models.Downloader, config *configuration.Configuration) {
	log.Println("Starting Download loop")
	for {
		for _, pack := range database.GetPackages() {
			if !pack.Finished {
				pack.Download(dl)
				go pack.Unrar(config.Dirs.ExtractDir)
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}
