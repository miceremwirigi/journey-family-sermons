package downloads

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ConverterResponse matches the successful "tunnel" response
type ConverterResponse struct {
	Status   string `json:"status"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

func (h *Handler) GetVideoDownloadLink(c *fiber.Ctx) error {
	videoID := c.Params("videoid")
	linkData, err := seekDownloadLink(videoID)
	if err != nil {
		return err
	}

	return c.JSON(linkData)
}

// Fetches a youtube video mp3 download link
func seekDownloadLink(videoID string) (*ConverterResponse, error) {
	iframeOrigin := "https://iframe.y2meta-uk.com"
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"

	client := &http.Client{Timeout: 30 * time.Second}

	// STEP 1: Fetch the Validation Key
	keyReq, _ := http.NewRequest("GET", "https://cnv.cx/v2/sanity/key", nil)
	keyReq.Header.Set("User-Agent", userAgent)
	keyReq.Header.Set("Origin", iframeOrigin)
	keyReq.Header.Set("Referer", iframeOrigin+"/")
	keyReq.Header.Set("Accept", "application/json, text/plain, */*")

	keyResp, err := client.Do(keyReq)
	if err != nil || keyResp.StatusCode != 200 {
		err = fmt.Errorf("failed to get key. status: %d. error %s", keyResp.StatusCode, err)
		return nil, err
	}
	defer keyResp.Body.Close()

	var keyData struct {
		Key string `json:"key"`
	}
	json.NewDecoder(keyResp.Body).Decode(&keyData)

	if keyData.Key == "" {
		err = fmt.Errorf("key was empty. kerver might be blocking us. error: %s", err)
		return nil, err
	}

	// STEP 2: Request the Conversion
	formData := url.Values{}
	formData.Set("link", "https://youtu.be/"+videoID)
	formData.Set("format", "mp4")
	formData.Set("audioBitrate", "128")
	formData.Set("videoQuality", "720")
	formData.Set("filenameStyle", "pretty")
	formData.Set("vCodec", "h264")

	convReq, _ := http.NewRequest("POST", "https://cnv.cx/v2/converter", strings.NewReader(formData.Encode()))

	convReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	convReq.Header.Set("key", keyData.Key)
	convReq.Header.Set("User-Agent", userAgent)
	convReq.Header.Set("Origin", iframeOrigin)
	convReq.Header.Set("Referer", iframeOrigin+"/")
	convReq.Header.Set("Accept", "*/*")

	convResp, err := client.Do(convReq)
	if err != nil {
		err = fmt.Errorf("Conversion Request Error: %s", err)
		return nil, err
	}
	defer convResp.Body.Close()

	// Parse the final result
	var result ConverterResponse
	bodyBytes, _ := io.ReadAll(convResp.Body)
	json.Unmarshal(bodyBytes, &result)

	if result.Status != "tunnel" && result.Status != "tunnel" || result.URL == "" {
		err = fmt.Errorf("Conversion failed. Server said: %s", string(bodyBytes))
		// tryFallback(videoID)
		return nil, err
	}

	return &result, nil
}

// Returns an alternative html page for downloading
func tryFallback(videoID string) {
	fmt.Println("Attempting fallback...")
	url := fmt.Sprintf("https://conv.mp3youtube.cc/download/%s", videoID)
	resp, _ := http.Get(url)
	if resp != nil {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Fallback Response:", string(body))
	}
}
