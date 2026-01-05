package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3/client"
	"github.com/miceremwirigi/journey-family-sermons/m/cmd/config"
	serializer "github.com/miceremwirigi/journey-family-sermons/m/pkg/serializers"
)

func FetchChannelVideosList(channelId string) *serializer.YouTubeResponse {
	conf := config.LoadConfig()

	var uploadsPlaylist string
	runes := []rune(channelId)
	if len(runes) > 1 {
		runes[1] = 'U'
	}
	uploadsPlaylist = string(runes)

	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=%s&channelType=any&type=video&maxResults=50&key=%s", uploadsPlaylist, conf.YoutubeDataApiKey)
	log.Println(url)

	client := client.New()
	client.AddHeader("Accept", "application/json")
	client.SetTimeout(30 * time.Second)

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	items := serializer.YouTubeResponse{}
	err = json.Unmarshal(resp.Body(), &items)
	if err != nil {
		panic(err)
	}

	log.Println(items)
	return &items
}

// FetchPlaylistMetadata gets the Title and total item count
func FetchPlaylistMetadata(playlistId string) *serializer.YouTubeResponse {
	conf := config.LoadConfig()
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?part=snippet,contentDetails&id=%s&key=%s", playlistId, conf.YoutubeDataApiKey)

	return executeYouTubeRequest(url, "")
}

// FetchPlaylistItems gets the videos to find contributors
// func FetchPlaylistItems(playlistId string) *serializer.YouTubeResponse {
// 	conf := config.LoadConfig()
// 	// We need 'snippet' to see videoOwnerChannelTitle
// 	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=%s&maxResults=50&key=%s", playlistId, conf.YoutubeDataApiKey)

// 	return executeYouTubeRequest(url, "")
// }

// Helper updated to accept a pageToken
func executeYouTubeRequest(baseUrl string, pageToken string) *serializer.YouTubeResponse {
	url := baseUrl
	if pageToken != "" {
		url = fmt.Sprintf("%s&pageToken=%s", baseUrl, pageToken)
	}

	c := client.New()
	c.AddHeader("Accept", "application/json")
	resp, err := c.Get(url)
	if err != nil {
		log.Printf("Request error: %v", err)
		return nil
	}
	defer resp.Close()

	var result serializer.YouTubeResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Printf("Unmarshal error: %v", err)
		return nil
	}
	return &result
}

// SyncSinglePlaylist gets a playlist's video metadata
// by looping through every page of the playlist
// using ne nextPageToken provided in every fetch.
// It ensures all pages are fetched.
func SyncSinglePlaylist(playlistId string) []serializer.YouTubeItem {
	conf := config.LoadConfig()
	baseUrl := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=%s&maxResults=50&key=%s", playlistId, conf.YoutubeDataApiKey)

	var allItems []serializer.YouTubeItem
	nextPageToken := ""

	for {
		resp := executeYouTubeRequest(baseUrl, nextPageToken)
		if resp == nil {
			break
		}

		// Append items from current page to our master list
		allItems = append(allItems, resp.Items...)

		// Move to the next page or exit if we're done
		nextPageToken = resp.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return allItems
}
