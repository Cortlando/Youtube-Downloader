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

	if len(vidsToDownload) == 0 {
		fmt.Println("There are no new videos to download")
		return
	}

	fmt.Println(vidsToDownload)

	downloadedVideos, errorList := youtube.DownloadYoutubeVideos(vidsToDownload)

	for _, e := range errorList {
		fmt.Println(e)
	}

	for _, e := range downloadedVideos {
		fmt.Println(e)
	}

	if len(downloadedVideos) == 0 {
		fmt.Println("No videos were downloaded")
		return
	}

	uploadedVideos, errorList2 := env.drop.UploadFiles(downloadedVideos)

	for _, e := range errorList2 {
		fmt.Println(e)
	}
	for _, e := range uploadedVideos {
		fmt.Println(e)
	}

	env.youtubevideomodel.InsertYoutubeVideosIntoTable(uploadedVideos)

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
func comparePlaylistAndDB(yt []youtube.ExtractedVideoInfo, db []string) []youtube.ExtractedVideoInfo {
	var videosToDownload []youtube.ExtractedVideoInfo

	for _, v := range yt {

		// if i > len(db) {
		// 	break
		// }

		if !slices.Contains(db, v.Video_ID) {
			// videosToDownload = append(videosToDownload, v.Video_ID)
			videosToDownload = append(videosToDownload, youtube.ExtractedVideoInfo{
				Title:      v.Title,
				Video_ID:   v.Video_ID,
				WebpageURL: v.WebpageURL,
			})
		}
	}

	return videosToDownload

}
