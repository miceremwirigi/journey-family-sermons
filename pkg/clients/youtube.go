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
