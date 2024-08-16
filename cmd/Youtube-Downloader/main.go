package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lrstanley/go-ytdlp"
)

type extractedVideoInfo struct {
	Title      string
	ID         string
	WebpageURL string
	// UploadDate string
}

func loadEnvVar() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadEnvUrlVar() string {
	url := os.Getenv("URL")

	return url
}

func extractString(s *string) string {
	if s == nil {
		return ""
	}
	return *s

}

func printVideos(videolist []extractedVideoInfo) {
	for _, v := range videolist {
		fmt.Printf("Title: %s, ID: %s, URL: %s\n", v.Title, v.ID, v.WebpageURL)
	}
}

func getVideosfromYoutubePlaylist() []extractedVideoInfo {
	youtubePlaylistUrl := loadEnvUrlVar()

	dl := ytdlp.New().
		PrintJSON().
		FlatPlaylist().
		// GetID().
		// GetDuration().
		// GetTitle().
		// GetURL().

		// NoProgress().
		// FormatSort("ba").
		ExtractAudio().
		// GetTitle().
		AudioQuality("0").
		// RecodeVideo("mp4").
		// NoPlaylist().
		// NoOverwrites().
		Continue()
		// Output("%(extractor)s - %(title)s.%(ext)s")

	playlist, err := dl.Run(context.TODO(), youtubePlaylistUrl)
	videosInPlaylist, err2 := playlist.GetExtractedInfo()

	var extractedVideosFromPlaylist []extractedVideoInfo

	for i, _ := range videosInPlaylist {
		var video = extractedVideoInfo{
			string(*videosInPlaylist[i].Title),
			string(videosInPlaylist[i].ID),
			string(*videosInPlaylist[i].WebpageURL),
		}

		// fmt.Print(video)
		extractedVideosFromPlaylist = append(extractedVideosFromPlaylist, video)
		// extracted_info = append(extracted_info, video)
	}
	if err != nil {
		panic(err)
	}

	if err2 != nil {
		panic(err)
	}

	return extractedVideosFromPlaylist

}

func main() {
	loadEnvVar()

	ytdlp.Install(context.TODO(), nil)

	// fmt.Printf("AAAAAAAAAAAA")
	var extractedVideosFromPlaylist []extractedVideoInfo = getVideosfromYoutubePlaylist()
	printVideos(extractedVideosFromPlaylist)

	// fmt.Println("Done")
}
