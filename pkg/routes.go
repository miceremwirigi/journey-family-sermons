package routes

import (
	"github.com/gofiber/fiber/v2"
	downloads "github.com/miceremwirigi/journey-family-sermons/m/pkg/apis/download"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/apis/videos"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	videoRoutes := app.Group("/videos")

	videoHandler := videos.Handler{}
	videoHandler.RegisterVideoRoutes(db, videoRoutes)

	downloadsRoutes := app.Group("/download")

	downloadsHandler := downloads.Handler{}
	downloadsHandler.RegisterDownloadRoutes(downloadsRoutes)
}
