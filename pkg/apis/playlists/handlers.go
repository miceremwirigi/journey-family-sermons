package playlists

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	app *fiber.App
	db *gorm.DB
}
