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

	_ = db.AutoMigrate(models.YoutubeVideo{})
	_ = db.AutoMigrate(models.ThumbnailData{})

	log.Println("Migration successful")
}
