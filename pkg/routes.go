package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/apis/videos"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	videoRoutes := app.Group("/videos")

	videoHandler := videos.Handler{}
	videoHandler.RegisterVideoRoutes(db, videoRoutes)
}
