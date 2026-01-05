package downloads

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterDownloadRoutes(r fiber.Router) {
	r.Get("/:videoid", h.GetVideoDownloadLink)
}
