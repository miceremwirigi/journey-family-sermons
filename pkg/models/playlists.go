package models

import (
	"time"
)

type YoutubePlaylist struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	Title        string         `json:"title"`
	Authors      []string       `gorm:"serializer:json" json:"authors"`
	Url          string         `json:"url"`
	ItemCount    int            `json:"item_count"`
	ThumbnailUrl string         `json:"thumbnail_url"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Videos       []YoutubeVideo `gorm:"many2many:playlist_videos;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
