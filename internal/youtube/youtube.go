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

	ytdlp.MustInstall(context.TODO(), nil)

	fmt.Print("Initizlized environment variable in youtube package")
}

func loadEnvUrlVar() string {
	url := os.Getenv("URL")

	return url
}

func GetVideosfromYoutubePlaylist() []ExtractedVideoInfo {
	youtubePlaylistUrl := loadEnvUrlVar()
	fmt.Print("1")
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
	// fmt.Print("2")
	playlist, err := dl.Run(context.TODO(), youtubePlaylistUrl)
	// fmt.Print("3")
	videosInPlaylist, _ := playlist.GetExtractedInfo()

	// fmt.Print(videosInPlaylist)
	var extractedVideosFromPlaylist []ExtractedVideoInfo

	for i := range videosInPlaylist {
		var video = ExtractedVideoInfo{
			string(*videosInPlaylist[i].Title),
			string(videosInPlaylist[i].ID),
			string(*videosInPlaylist[i].WebpageURL),
		}

		// fmt.Print(video)
		extractedVideosFromPlaylist = append(extractedVideosFromPlaylist, video)
		// extracted_info = append(extracted_info, video)
	}

	fmt.Print(playlist)

	if err != nil {
		fmt.Print("ERROR WENT OFF")
		panic(err)
	}

	fmt.Print(dl)
	fmt.Print(youtubePlaylistUrl)
	return extractedVideosFromPlaylist

}

func PrintVideos(videolist []ExtractedVideoInfo) {
	for _, v := range videolist {
		fmt.Printf("Title: %s, ID: %s, URL: %s\n", v.Title, v.Video_ID, v.WebpageURL)
	}
}
