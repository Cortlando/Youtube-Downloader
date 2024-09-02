package youtube

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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

func DownloadYoutubeVideos(videosToDownload []string) []error {
	wg := sync.WaitGroup{}

	errCh := make(chan error, 1)

	for _, video := range videosToDownload {
		wg.Add(1)
		go downloadYoutubeVideo(video, &wg, errCh)
	}

	go func() {
		wg.Wait()

		close(errCh)
	}()

	var errorList []error
	for e := range errCh {
		if e != nil {

			errorList = append(errorList, e)
		}

	}

	return errorList
}

func downloadYoutubeVideo(videoId string, wg *sync.WaitGroup, errCh chan<- error) error {

	defer wg.Done()
	dl := ytdlp.New().
		FlatPlaylist().
		ExtractAudio().
		AudioQuality("0").
		Paths("/downloads").
		Continue()

	var ytString strings.Builder

	ytString.WriteString("https://www.youtube.com/watch?v=")
	ytString.WriteString(videoId)

	result, err := dl.Run(context.TODO(), ytString.String())

	errCh <- err
	fmt.Print(result)

	return nil

}
