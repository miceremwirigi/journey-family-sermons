package videos

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/miceremwirigi/journey-family-sermons/m/cmd/config"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/apis"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/clients"
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/models"
	serializer "github.com/miceremwirigi/journey-family-sermons/m/pkg/serializers"
	"gorm.io/gorm"
)

func (h *Handler) GetVideosList(c *fiber.Ctx) error {
	// channelId := c.Params("channelId", "UCxs87BsBSzw8FC7IlvfQ5kA")
	// fmt.Println(channelId)
	log.Println("fetching from db")
	videosList := []models.YoutubeVideo{}
	err := h.db.Model(models.YoutubeVideo{}).Find(&videosList).Error
	if err != nil {
		return c.Status(404).JSON(apis.ErrorDataResponse(err.Error(),"Fauld to retraive video list fom DB", 404))
	}

	log.Println("Fetch successful")

	return c.JSON(videosList)
}

func (h Handler) FetchVideoData(c *fiber.Ctx) error {
	conf := config.LoadConfig()
	videosList := []models.YoutubeVideo{}
	apiResponse := *clients.FetchChannelVideosList(conf.JourneyFamilyChannelID)

	for _, item := range apiResponse.Items {
		video := serializer.MapYouTubeItemToModel(item)
		videosList = append(videosList, video)

		// add video to database video unless it already exists 
		err := saveVideoIfNotExists(h.db, video)
		if err != nil {
			msg := fmt.Sprintf("Failed to create video %s", video.ID)
			panic(msg)
		}

	}

	fmt.Printf("Successfully updated videos list %d", len(videosList))

	return nil
}

func saveVideoIfNotExists(db *gorm.DB, video models.YoutubeVideo) error {
	result := db.Where(models.YoutubeVideo{ID: video.ID}).FirstOrCreate(&video)
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
