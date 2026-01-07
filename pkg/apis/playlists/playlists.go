package playlists

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/apis"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/clients"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/models"
	serializer "github.com/miceremwirigi/journey-family-sermons/m/pkg/serializers"
)

type addPlaylistRequest struct {
	url string
}

// Retreives all playlists from db
func (h *Handler) GetAllPlaylists(c *fiber.Ctx) error {
	playlists := &[]models.YoutubePlaylist{}

	err := h.db.Model(models.YoutubePlaylist{}).Find(playlists).Error
	if err != nil {
		return c.Status(404).JSON(apis.ErrorDataResponse(err.Error(), "Failed to retreive video list from  DB", 404))
	}
	return c.JSON(playlists)
}

// Get one whole playlist with info from db
func (h *Handler) GetPlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlist := &models.YoutubePlaylist{}

	err := h.db.Preload("Videos").Model(&models.YoutubePlaylist{}).First(playlist, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(apis.ErrorDataResponse(err.Error(), "Could not find playlist with id "+id, 404))
	}
	return c.JSON(playlist)
}

// Get all videos of a playlist without exta playlist info
func (h *Handler) GetPlaylistVideos(c *fiber.Ctx) error {
	playlistID := c.Params("id")
	var playlist models.YoutubePlaylist

	// Use Preload to fetch the associated videos through the junction table
	err := h.db.Model(&models.YoutubePlaylist{}).Preload("Videos").First(&playlist, "id = ?", playlistID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Playlist not found"})
	}

	// Returns only the videos belonging to THIS playlist
	return c.JSON(playlist.Videos)
}

// Add a playlist to be tracked in db
func (h *Handler) AddPlaylist(c *fiber.Ctx) error {
	var body struct {
		Url string `json:"url"`
	}
	c.BodyParser(&body)

	// Extract ID
	u, _ := url.Parse(body.Url)
	pID := u.Query().Get("list")

	// 1. Get Playlist Info
	metaResp := clients.FetchPlaylistMetadata(pID)
	if len(metaResp.Items) == 0 {
		return c.Status(404).JSON(apis.ErrorDataResponse("Zero items", "Playlist not found", 404))
	}
	playlistItem := metaResp.Items[0]

	// 2. Fetch every single item in the playlist to populate collaborators
	allItems := clients.SyncSinglePlaylist(pID)
	authorMap := make(map[string]bool)

	// Always include the main owner
	authorMap[playlistItem.Snippet.ChannelTitle] = true

	// Check every video for a different owner
	for _, item := range allItems {
		if item.Snippet.VideoOwnerChannelTitle != "" {
			authorMap[item.Snippet.VideoOwnerChannelTitle] = true
		}
	}

	// Convert map to slice
	authors := make([]string, 0, len(authorMap))
	for name := range authorMap {
		authors = append(authors, name)
	}

	// 3. Save to DB
	playlist := models.YoutubePlaylist{
		ID:           pID,
		Title:        playlistItem.Snippet.Title,
		Authors:      authors,
		Url:          body.Url,
		ItemCount:    playlistItem.ContentDetails.ItemCount,
		ThumbnailUrl: playlistItem.Snippet.Thumbnails.High.URL,
	}

	h.db.Where(models.YoutubePlaylist{ID: pID}).Attrs(playlist).FirstOrCreate(&playlist)
	return c.JSON(playlist)
}

func (h *Handler) SyncSinglePlaylist(c *fiber.Ctx) error {
	playlistID := c.Params("id")

	// Get the playlist from DB
	var playlist models.YoutubePlaylist
	if err := h.db.First(&playlist, "id = ?", playlistID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Playlist not found"})
	}

	// Fetch from YouTube API
	apiItems := clients.SyncSinglePlaylist(playlistID)

	for _, item := range apiItems {
		video := serializer.MapYouTubeItemToModel(item)

		// Save/Update the video record
		h.db.Where(models.YoutubeVideo{ID: video.ID}).Attrs(video).FirstOrCreate(&video)

		// 4. Create the relationship in the junction table
		// This ensures the video is "pinned" to this specific playlist
		h.db.Model(&playlist).Association("Videos").Append(&video)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Playlist videos updated"})
}

func (h *Handler) SyncAllPlaylists(c *fiber.Ctx) error {
	// Get the playlist from DB
	var playlists []models.YoutubePlaylist
	if err := h.db.Find(&playlists).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Playlist not found"})
	}

	for _, playlist := range playlists {
		// Fetch from YouTube API
		apiItems := clients.SyncSinglePlaylist(playlist.ID)

		for _, item := range apiItems {
			video := serializer.MapYouTubeItemToModel(item)

			// Save/Update the video record
			h.db.Where(models.YoutubeVideo{ID: video.ID}).Attrs(video).FirstOrCreate(&video)

			// 4. Create the relationship in the junction table
			// This ensures the video is "pinned" to this specific playlist
			h.db.Model(&playlist).Association("Videos").Append(&video)
		}
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Playlist videos updated"})
}

func (h *Handler) DeletePlaylist(c *fiber.Ctx) error {
	id := c.Params("id")
	playlist := &models.YoutubePlaylist{}

	err := h.db.First(playlist, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(apis.ErrorDataResponse(err.Error(), "Could not find playlistwith id "+id, 404))
	}

	err = h.db.Model(playlist).Association("Videos").Clear()
	if err != nil {
		return c.Status(500).JSON(apis.ErrorDataResponse(err.Error(), "Failed to clear playlist associations", 500))
	}

	err = h.db.Delete(playlist).Error
	if err != nil {
		return c.Status(500).JSON(apis.ErrorDataResponse(err.Error(), "Failed to delete playlist", 500))
	}
	return c.Status(200).JSON(map[string]interface{}{"message": "successfully deleted playlist"})
}
