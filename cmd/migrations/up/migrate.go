package main

import (
	"log"

	"github.com/miceremwirigi/journey-family-sermons/m/cmd/config"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/databases"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/models"
)

func main() {
	log.Println("Running database up migrations ...")
	conf := config.LoadConfig()
	db := databases.StartDatabase(conf.Environment)

	err := db.AutoMigrate(models.YoutubeVideo{}); 
	if err != nil {
		log.Fatal(err.Error())
	}
	err = db.AutoMigrate(models.ThumbnailData{}); 
	if err != nil {
		log.Fatal(err.Error())
	}
	err = db.AutoMigrate(models.YoutubePlaylist{}); 
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Migration successful")
}
