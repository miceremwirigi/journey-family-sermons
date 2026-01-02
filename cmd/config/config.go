package config

import "os"

type Config struct {
	DatabaseUrl string
	TestDatabaseUrl string
	YoutubeDataApiKey string
	JourneyFamilyChannelID string
	JourneyFamilyVideoListUrl string
	Environment string
}

func LoadConfig() *Config {
	var c Config
	c.YoutubeDataApiKey, _ =	os.LookupEnv("database_url")
	c.YoutubeDataApiKey, _ =	os.LookupEnv("test_database_url")
	c.YoutubeDataApiKey, _ =	os.LookupEnv("youtube_data_api_key")
	c.JourneyFamilyChannelID, _ = os.LookupEnv("journey_family_channel_id")
	c.JourneyFamilyVideoListUrl, _ =	os.LookupEnv("journey_family_video_list_url")
	c.Environment, _ =	os.LookupEnv("env")
	return &c
}