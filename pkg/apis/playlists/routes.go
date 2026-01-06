package playlists

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (h *Handler) RegisterPlaylistRoutes(db *gorm.DB, r fiber.Router) {
	h.db = db

	r.Get("/", h.GetAllPlaylists)
	r.Get("/info/:id", h.GetPlaylist)
	r.Get("/videos/:id", h.GetPlaylistVideos)
	r.Post("/sync/:id", h.SyncSinglePlaylist)
	r.Post("/sync", h.SyncAllPlaylists)
	r.Post("/add", h.AddPlaylist)
	r.Delete("/:id", h.DeletePlaylist)
}
