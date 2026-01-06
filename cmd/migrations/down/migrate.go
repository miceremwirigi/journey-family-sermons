package main

import (
	"log"

	"github.com/miceremwirigi/journey-family-sermons/m/cmd/config"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/databases"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/models"
)

func main() {
	log.Println("Running database down migrations ...")
	conf := config.LoadConfig()
	db := databases.StartDatabase(conf.Environment)

	_ = db.Migrator().DropTable(models.YoutubeVideo{})
	_ = db.Migrator().DropTable(models.ThumbnailData{})
	_ = db.Migrator().DropTable(models.YoutubePlaylist{})

	log.Println("Migration successful")
}
