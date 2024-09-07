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
	Extension  string
	// UploadDate string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ytdlp.MustInstall(context.TODO(), nil)

	fmt.Println("Initizlized environment variable in youtube package")
}

func loadEnvUrlVar() string {
	url := os.Getenv("URL")

	return url
}

func GetVideosfromYoutubePlaylist() []ExtractedVideoInfo {
	fmt.Println("Running GetVideosfromYoutubePlaylist")
	youtubePlaylistUrl := loadEnvUrlVar()

	dl := ytdlp.New().
		PrintJSON().
		FlatPlaylist().
		ExtractAudio().
		AudioQuality("0").
		Continue()

	playlist, err := dl.Run(context.TODO(), youtubePlaylistUrl)
	fmt.Println(playlist)
	videosInPlaylist, _ := playlist.GetExtractedInfo()

	var extractedVideosFromPlaylist []ExtractedVideoInfo

	for i := range videosInPlaylist {

		var video = ExtractedVideoInfo{
			//Removes forward slashes because they cause uploading to break
			strings.ReplaceAll(*videosInPlaylist[i].Title, "/", " "),
			string(videosInPlaylist[i].ID),
			string(*videosInPlaylist[i].WebpageURL),
			string(*videosInPlaylist[i].URL),
		}

		fmt.Println(string(videosInPlaylist[i].Extension))
		// fmt.Print(video)
		extractedVideosFromPlaylist = append(extractedVideosFromPlaylist, video)
		// extracted_info = append(extracted_info, video)
	}

	// fmt.Print(playlist)

	if err != nil {
		fmt.Print("ERROR WENT OFF")
		panic(err)
	}

	// cleanVideoTitles(&extractedVideosFromPlaylist)
	// fmt.Print(dl)
	// fmt.Print(youtubePlaylistUrl)
	return extractedVideosFromPlaylist

}

func PrintVideos(videolist []ExtractedVideoInfo) {
	for _, v := range videolist {
		fmt.Printf("Title: %s, ID: %s, URL: %s, Ext: %s", v.Title, v.Video_ID, v.WebpageURL, v.Extension)
	}
}

// TODO: Figure out what is causing the program to freeze
// I think its something to do with errCH
func DownloadYoutubeVideos(videosToDownload []ExtractedVideoInfo) (map[string]ExtractedVideoInfo, []error) {
	wg := sync.WaitGroup{}

	errCh := make(chan error, len(videosToDownload))
	videoCh := make(chan ExtractedVideoInfo, len(videosToDownload))
	// var downloadedVideos []ExtractedVideoInfo
	var downloadedVideos = make(map[string]ExtractedVideoInfo)
	var errorList []error

	for _, video := range videosToDownload {
		wg.Add(1)
		go downloadYoutubeVideo(video, &wg, errCh, videoCh)
	}

	go func() {
		wg.Wait()
		// fmt.Print(len(errCh))
		// fmt.Print(len(videoCh))
		close(videoCh)
		close(errCh)
	}()

	if len(errCh) > 0 {
		for e := range errCh {
			if e != nil {

				errorList = append(errorList, e)
			}

		}
	}

	//Returns a map with the id as key, since the downloaded video titles are there id
	//I'll be able to link the file to its corresponding struct
	for v := range videoCh {

		downloadedVideos[v.Video_ID] = v

	}

	return downloadedVideos, errorList
}

// TODO: Investigate why [] getting added to video titles
func downloadYoutubeVideo(video ExtractedVideoInfo, wg *sync.WaitGroup, errCh chan<- error, videoCh chan<- ExtractedVideoInfo) error {
	defer wg.Done()

	dl := ytdlp.New().
		FlatPlaylist().
		ExtractAudio().
		AudioQuality("0").
		Paths("./downloads").
		Output("%(id)s.%(ext)s").
		RestrictFilenames().
		Continue()

	result, err := dl.Run(context.TODO(), video.WebpageURL)

	errCh <- err

	//If the video downloaded, add it to the video channel
	if err == nil {
		videoCh <- video
	}

	fmt.Println(result.Stdout)

	return nil

}
