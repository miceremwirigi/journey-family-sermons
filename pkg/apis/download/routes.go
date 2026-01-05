package downloads

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterDownloadRoutes(r fiber.Router) {
	r.Get("/mp3/:videoid", h.GetMp3DownloadLink)
	r.Get("/mp4/:videoid", h.GetMp4DownloadLink)
}
