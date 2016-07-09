package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"goload/server/controllers"
	"goload/server/models"
	"goload/server/data"
	"time"
	"goload/server/models/configuration"
	"os"
	"os/signal"
)


var Version = "dev"

func main() {
	config, error := configuration.NewConfigurationFromFileName("config.json")
	if(error!= nil) {
		log.Fatal(error)
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

	ul := models.NewUploaded(config)
	database := data.NewDatastore()
	packageController := controllers.NewPackageController(database,ul)
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
		versionJson := []byte(`{"version":"`+Version+`"}`)
		w.Write(versionJson)
	})
	go LoopPackages(database, ul,config);
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for range c {
			log.Println("Shutting Down")
			os.Exit(0)
		}
	}()
	log.Fatal(http.ListenAndServe(":3000", router))

}

func LoopPackages(database *data.Datastore, ul *models.Uploaded,config *configuration.Configuration) {
	log.Println("Starting Download loop")
	for {
		for _, pack := range database.GetPackages() {
			if (!pack.Finished) {
				pack.Download(ul)
				go pack.Unrar(config.Dirs.ExtractDir)
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}