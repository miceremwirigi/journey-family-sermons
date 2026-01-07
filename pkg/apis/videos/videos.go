package videos

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/miceremwirigi/journey-family-sermons/m/cmd/config"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/apis"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/clients"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/models"
	serializer "github.com/miceremwirigi/journey-family-sermons/m/pkg/serializers"
	"gorm.io/gorm"
)

// Get Sermons from JouneyFamily uploads
func (h *Handler) GetVideosList(c *fiber.Ctx) error {
	var mainPlaylist models.YoutubePlaylist
	conf := config.LoadConfig()
	runes := []rune(conf.JourneyFamilyChannelID)
	if len(runes) > 1 {
		runes[1] = 'U'
	}
	mainPlaylistID := string(runes)

	// Fetch only videos linked to the Main Channel ID
	err := h.db.Preload("Videos").First(&mainPlaylist, "id = ?", mainPlaylistID).Error
	if err != nil {
		return c.Status(404).JSON(apis.ErrorDataResponse(err.Error(), "Failed to retreive sermon video list from DB", 404))
	}

	return c.JSON(mainPlaylist.Videos)
}

// Fetches from youtube all links to videos of the JourneyFamilyChannel and saves to db
func (h Handler) FetchVideoData(c *fiber.Ctx) error {
	conf := config.LoadConfig()

	// Find uploads playlist id from channel id
	runes := []rune(conf.JourneyFamilyChannelID)
	if len(runes) > 1 {
		runes[1] = 'U'
	}
	mainPlaylistID := string(runes)
	var mainPlaylist models.YoutubePlaylist
	h.db.FirstOrCreate(&mainPlaylist, models.YoutubePlaylist{ID: mainPlaylistID, Title: "All Sermons"})

	apiResponse := *clients.FetchChannelVideosList(conf.JourneyFamilyChannelID)

	for _, item := range apiResponse.Items {
		video := serializer.MapYouTubeItemToModel(item)

		// Save video
		err := saveVideoIfNotExists(h.db, video)
		if err != nil {
			return c.Status(500).JSON(apis.ErrorDataResponse(err.Error(), "Failed to save video info: "+video.Title, 500))
		}

		// Link video to the "All Sermons" playlist
		err = h.db.Model(&mainPlaylist).Association("Videos").Append(&video)
		if err != nil {
			return c.Status(500).JSON(apis.ErrorDataResponse(err.Error(), "Failed to link video associations: "+video.Title, 500))
		}
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func saveVideoIfNotExists(db *gorm.DB, video models.YoutubeVideo) error {
	result := db.Where(models.YoutubeVideo{ID: video.ID}).Attrs(video).FirstOrCreate(&video)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (h Handler) RepairVideoData(c *fiber.Ctx) error {
	conf := config.LoadConfig()
	videosList := []models.YoutubeVideo{}
	apiResponse := *clients.FetchChannelVideosList(conf.JourneyFamilyChannelID)

	for _, item := range apiResponse.Items {
		video := serializer.MapYouTubeItemToModel(item)
		videosList = append(videosList, video)

		// add video to database video unless it already exists
		err := saveVideoEvenIfItExists(h.db, video)
		if err != nil {
			msg := fmt.Sprintf("Failed to create or update video %s", video.ID)
			panic(msg)
		}

	}

	fmt.Printf("Successfully updated videos list %d", len(videosList))

	return nil
}

func saveVideoEvenIfItExists(db *gorm.DB, video models.YoutubeVideo) error {
	result := db.Where(models.YoutubeVideo{ID: video.ID}).Assign(models.YoutubeVideo{ID: video.ID}).FirstOrCreate(&video)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
