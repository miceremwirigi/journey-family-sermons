package serializer

import "time"

type YouTubeResponse struct {
	Kind          string        `json:"kind"`
	Etag          string        `json:"etag"`
	NextPageToken string        `json:"nextPageToken"`
	PageInfo      PageInfo      `json:"pageInfo"`
	Items         []YouTubeItem `json:"items"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type YouTubeItem struct {
	Kind           string         `json:"kind"`
	Etag           string         `json:"etag"`
	ID             string         `json:"id"`
	Snippet        Snippet        `json:"snippet"`
	ContentDetails ContentDetails `json:"contentDetails"`
}

type ContentDetails struct {
	ItemCount int `json:"itemCount"`
}

type Snippet struct {
	Title                  string       `json:"title"`
	Description            string       `json:"description"`
	PublishedAt            time.Time    `json:"publishedAt"`
	ChannelId              string       `json:"channelId"`
	Thumbnails             ThumbnailSet `json:"thumbnails"`
	ChannelTitle           string       `json:"channelTitle"`
	PlaylistId             string       `json:"playlistId"`
	Position               int          `json:"position"`
	ResourceId             ResourceId   `json:"resourceId"`
	VideoOwnerChannelTitle string       `json:"videoOwnerChannelTitle"`
	VideoOwnerChannelId    string       `json:"videoOwnerChannelId"`
}

type ResourceId struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type ThumbnailSet struct {
	Default  Thumbnail `json:"default"`
	Medium   Thumbnail `json:"medium"`
	High     Thumbnail `json:"high"`
	Standard Thumbnail `json:"standard"`
	Maxres   Thumbnail `json:"maxres"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
