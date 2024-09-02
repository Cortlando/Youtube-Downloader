package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"

	"github.com/cortlando/youtube-downloader/internal/drop"
	"github.com/cortlando/youtube-downloader/internal/sqlfuncs"
	"github.com/cortlando/youtube-downloader/internal/youtube"

	_ "github.com/mattn/go-sqlite3"
)

type Env struct {
	youtubevideomodel sqlfuncs.YoutubeVideoModel
	drop              drop.DropboxModel
}

func main() {
	// loadEnvVar()
	db, err := sql.Open("sqlite3", sqlfuncs.DB_PATH)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	dropboxUser := drop.InitDropbox()

	env := &Env{
		youtubevideomodel: sqlfuncs.YoutubeVideoModel{DB: db},
		drop:              dropboxUser,
	}

	var extractedVideosFromPlaylist []youtube.ExtractedVideoInfo = youtube.GetVideosfromYoutubePlaylist()
	youtube.PrintVideos(extractedVideosFromPlaylist)

	env.initializeDB()
	env.youtubevideomodel.CreateYoutubeVideoTableIfNotExist()
	// env.youtubevideomodel.TestInsertIntoTable()

	vidsInDB, err := env.youtubevideomodel.GetAllYoutubeVideoIDs()

	if err != nil {
		log.Fatal(err)
	}

	vidsToDownload := comparePlaylistAndDB(extractedVideosFromPlaylist, vidsInDB)

	// err = env.drop.GetAccount()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Print(vidsToDownload)

	// var path string = "./test2.txt"
	// err = env.drop.UploadFile(&path)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	youtube.DownloadYoutubeVideos(vidsToDownload)
}

func (env *Env) initializeDB() {
	env.youtubevideomodel.CheckIfDBFIleExists()
}

func (env *Env) initializeYoutubeTable() {
	err := env.youtubevideomodel.CreateYoutubeVideoTableIfNotExist()

	if err != nil {
		log.Fatal(err)
	}
}

func (env *Env) testquery() {
	err := env.youtubevideomodel.TestInsertIntoTable()

	if err != nil {
		log.Fatal(err)
	}
}

// Checks videos currently in playlist to videos in db
// returns slice that contains all videos in playlist that are not in db
func comparePlaylistAndDB(yt []youtube.ExtractedVideoInfo, db []string) []string {
	var videosToDownload []string

	for _, v := range yt {

		// if i > len(db) {
		// 	break
		// }

		if !slices.Contains(db, v.Video_ID) {
			videosToDownload = append(videosToDownload, v.Video_ID)
		}
	}

	return videosToDownload

}
