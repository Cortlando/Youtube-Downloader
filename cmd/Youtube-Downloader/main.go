package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	// "github.com/cortlando/youtube-downloader/internal/sqlite"

	"github.com/joho/godotenv"
	"github.com/lrstanley/go-ytdlp"

	_ "github.com/mattn/go-sqlite3"
)

var DB_PATH string = "./downloadedfiles.db?_journal=WAL&_timeout=5000"

type extractedVideoInfo struct {
	Title      string
	Video_ID   string
	WebpageURL string
	// UploadDate string
}

type Env struct {
	videos YoutubeVideoModel
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
		fmt.Printf("Title: %s, ID: %s, URL: %s\n", v.Title, v.Video_ID, v.WebpageURL)
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
	db, err := sql.Open("sqlite3", DB_PATH)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	env := &Env{
		videos: YoutubeVideoModel{DB: db},
	}

	// env.getVideosFromDB()

	// ytdlp.Install(context.TODO(), nil)

	fmt.Printf("AAAAAAAAAAAA")
	var extractedVideosFromPlaylist []extractedVideoInfo = getVideosfromYoutubePlaylist()
	printVideos(extractedVideosFromPlaylist)

	env.initializeDB()
	env.videos.createYoutubeVideoTableIfNotExist()
	env.videos.testInsertIntoTable()
	videolist, err := env.videos.testSelectFromTable()

	for _, v := range videolist {
		fmt.Print(v)
	}
	fmt.Println("Done")
}

func (env *Env) initializeDB() {
	env.videos.CheckIfDBFIleExists()
}

func (env *Env) initializeYoutubeTable() {
	err := env.videos.createYoutubeVideoTableIfNotExist()

	if err != nil {
		log.Fatal(err)
	}
}

func (env *Env) testquery() {
	err := env.videos.testInsertIntoTable()

	if err != nil {
		log.Fatal(err)
	}
}

// func (env *Env) getVideosFromDB() {

// }
