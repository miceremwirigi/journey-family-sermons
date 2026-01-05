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

func (h *Handler) GetMp4DownloadLink(c *fiber.Ctx) error {
	videoID := c.Params("videoid")
	linkData, err := seekDownloadLinkmp4(videoID)
	if err != nil {
		return err
	}

	return c.JSON(linkData)
}

func (h *Handler) GetMp3DownloadLink(c *fiber.Ctx) error {
	videoID := c.Params("videoid")
	linkData, err := seekDownloadLinkmp3(videoID)
	if err != nil {
		return err
	}

	return c.JSON(linkData)
}

// Fetches a youtube video mp4 download link
func seekDownloadLinkmp4(videoID string) (*ConverterResponse, error) {
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

// Fetches a youtube video mp3 download link
func seekDownloadLinkmp3(videoID string) (*ConverterResponse, error) {
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
        return nil, fmt.Errorf("failed to get key. status: %d", keyResp.StatusCode)
    }
    defer keyResp.Body.Close()

    var keyData struct {
        Key string `json:"key"`
    }
    json.NewDecoder(keyResp.Body).Decode(&keyData)

    // STEP 2: Request the Conversion (Modified for MP3)
    formData := url.Values{}
    formData.Set("link", "https://youtu.be/"+videoID)
    formData.Set("format", "mp3")           
    formData.Set("audioBitrate", "320")     //  320 for highest quality
    formData.Set("filenameStyle", "pretty")
    
    // Note: videoQuality and vCodec are usually ignored by servers when format is mp3, 
    // but you can remove them to keep the request clean.

    convReq, _ := http.NewRequest("POST", "https://cnv.cx/v2/converter", strings.NewReader(formData.Encode()))

    convReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    convReq.Header.Set("key", keyData.Key)
    convReq.Header.Set("User-Agent", userAgent)
    convReq.Header.Set("Origin", iframeOrigin)
    convReq.Header.Set("Referer", iframeOrigin+"/")

    convResp, err := client.Do(convReq)
    if err != nil {
        return nil, err
    }
    defer convResp.Body.Close()

    var result ConverterResponse
    bodyBytes, _ := io.ReadAll(convResp.Body)
    json.Unmarshal(bodyBytes, &result)

    if (result.Status != "tunnel" && result.Status != "success") || result.URL == "" {
        return nil, fmt.Errorf("conversion failed: %s", string(bodyBytes))
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
