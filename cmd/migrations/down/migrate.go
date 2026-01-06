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
	_ = db.Migrator().DropTable(models.PlaylistVideo{}) // Join table for the playlist to video many to many

	log.Println("Migration successful")
}
