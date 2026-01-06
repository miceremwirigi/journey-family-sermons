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

// In pkg/serializers/youtube.go
func MapPlaylistResponseToModel(pItem YouTubeItem, collaborators []string, fullUrl string) models.YoutubePlaylist {
	return models.YoutubePlaylist{
		ID:           pItem.ID,
		Title:        pItem.Snippet.Title,
		Authors:      collaborators,
		Url:          fullUrl,
		ItemCount:    pItem.ContentDetails.ItemCount, // Note: You need to add ContentDetails to your YouTubeItem struct
		ThumbnailUrl: pItem.Snippet.Thumbnails.High.URL,
	}
}
