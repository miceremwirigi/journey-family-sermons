package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type YoutubeVideo struct {
	ID          string `gorm:"primaryKey"`
	Title       string `gorm:"index"`
	Description string
	PublishedAt time.Time     `gorm:"index"`
	Thumbnails  ThumbnailData `gorm:"type:jsonb"` // Stores thumbnails as a JSON column in DB
}

type ThumbnailData struct {
	DefaultURL string    `gorm:"primaryKey"`
	MediumURL  string    	
	HighURL    string    	
	Standard   string 	
	Maxres     string 	
}

// Value implements the driver.Valuer interface for GORM to save as JSON
func (t ThumbnailData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan implements the sql.Scanner interface for GORM to read from JSON
func (t ThumbnailData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
