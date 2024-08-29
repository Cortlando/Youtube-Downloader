package youtube

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lrstanley/go-ytdlp"
)

type ExtractedVideoInfo struct {
	Title      string
	Video_ID   string
	WebpageURL string
	// UploadDate string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadEnvUrlVar() string {
	url := os.Getenv("URL")

	return url
}

func GetVideosfromYoutubePlaylist() []ExtractedVideoInfo {
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

	var extractedVideosFromPlaylist []ExtractedVideoInfo

	for i, _ := range videosInPlaylist {
		var video = ExtractedVideoInfo{
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

func PrintVideos(videolist []ExtractedVideoInfo) {
	for _, v := range videolist {
		fmt.Printf("Title: %s, ID: %s, URL: %s\n", v.Title, v.Video_ID, v.WebpageURL)
	}
}
