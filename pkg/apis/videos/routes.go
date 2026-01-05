package videos

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (h *Handler) RegisterVideoRoutes(db *gorm.DB, r fiber.Router) {
	h.db = db
	r.Get("/", h.GetVideosList)
	r.Post("/sync", h.FetchVideoData)
	r.Get("/repair", h.FetchVideoData)
}
