package serializer

import (
	"github.com/miceremwirigi/journey-family-sermons/m/pkg/models"
)

func MapYouTubeItemToModel(item YouTubeItem) models.YoutubeVideo {
	return models.YoutubeVideo{
		ID:          item.Snippet.ResourceId.VideoID,
		Title:       item.Snippet.Title,
		Description: item.Snippet.Description,
		PublishedAt: item.Snippet.PublishedAt,
		Thumbnails: models.ThumbnailData{
			DefaultURL: item.Snippet.Thumbnails.Default.URL,
			MediumURL:  item.Snippet.Thumbnails.Medium.URL,
			HighURL:    item.Snippet.Thumbnails.High.URL,
		},
	}
}
