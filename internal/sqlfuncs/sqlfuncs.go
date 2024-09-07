package sqlfuncs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"os"

	"github.com/cortlando/youtube-downloader/internal/youtube"
	_ "github.com/mattn/go-sqlite3"
)

const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

type YoutubeVideo struct {
	title         string
	Video_ID      string
	WebpageURL    string
	download_date string
}

type YoutubeVideoModel struct {
	DB *sql.DB
}

// const (
// 	DB_PATH string = "./downloadedfiles.db?_journal=WAL&_timeout=5000"
// )

// func (m YoutubeVideoModel) CheckIfDBFIleExists() {
// 	if _, err := os.Stat(DB_PATH); err == nil {
// 		fmt.Println("DB File exists")
// 	} else if errors.Is(err, os.ErrNotExist) {
// 		file, err := os.Create("./db/downloadedfiles.db")

// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer file.Close()
// 	} else {
// 		fmt.Println("This shouldn't happen")
// 	}

// }

func CheckIfDBFIleExists() error {

	if _, err := os.Stat(os.Getenv("DB_PATH")); err == nil {
		fmt.Println("DB File exists")
	} else if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("./db", 0777)
		if err != nil {
			return err
		}

		file, err := os.Create(os.Getenv("DB_PATH"))

		if err != nil {
			fmt.Print(err.Error())
			return err
		}
		defer file.Close()

	} else {
		fmt.Println("This shouldn't happen")
	}

	return nil

}

func (m YoutubeVideoModel) CreateYoutubeVideoTableIfNotExist() error {

	// return err
	stmt := `CREATE TABLE IF NOT EXISTS youtubevideos(
		Title TEXT,
		video_id TEXT PRIMARY KEY,
		 webpage_url TEXT,
		 download_date TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ'))
		 ) strict`
	i, err := m.DB.Exec(stmt)
	println(i)
	return err
}

func (m YoutubeVideoModel) TestInsertIntoTable() error {
	_, err := m.DB.Exec(`INSERT INTO youtubevideos(title, video_id, webpage_url) values('b','b','b')`)

	return err
}

func (m YoutubeVideoModel) InsertYoutubeVideosIntoTable(uploadedVideos []youtube.ExtractedVideoInfo) {
	stmt, err := m.DB.PrepareContext(context.TODO(), "INSERT INTO youtubevideos(title, video_id, webpage_url) values (?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	t0 := time.Now()

	for _, v := range uploadedVideos {
		_, err := stmt.ExecContext(context.TODO(), v.Title, v.Video_ID, v.WebpageURL)

		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Printf("\nInserts took: %v\n", time.Since(t0))

}

// func (m YoutubeVideoModel) InsertYoutubeVideosIntoTable(uploadedVideos []youtube.ExtractedVideoInfo) {
// 	stmt, err := m.DB.PrepareContext(context.TODO(), "INSERT INTO youtubevideos(title, video_id, webpage_url) values (?, ?, ?)")
// 	wg := sync.WaitGroup{}
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	t0 := time.Now()

// 	for _, v := range uploadedVideos {
// 		wg.Add(1)

// 		go func() {
// 			defer wg.Done()
// 			_, err := stmt.ExecContext(context.TODO(), v.Title, v.Video_ID, v.WebpageURL)

// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 		}()

// 		// _, err := stmt.ExecContext(context.TODO(), v.Title, v.Video_ID, v.WebpageURL)

// 		// if err != nil {
// 		// 	log.Fatal(err)
// 		// }

// 	}

// 	fmt.Printf("\nInserts took: %v\n", time.Since(t0))
// }

func (m YoutubeVideoModel) TestSelectFromTable() ([]YoutubeVideo, error) {
	rows, err := m.DB.Query("SELECT * FROM  youtubevideos")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var videos []YoutubeVideo

	for rows.Next() {
		var video YoutubeVideo
		if err := rows.Scan(&video.title, &video.Video_ID, &video.WebpageURL, &video.download_date); err != nil {
			return videos, err
		}

		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return videos, err
	}

	return videos, nil

}

func (m YoutubeVideoModel) GetAllYoutubeVideoIDs() ([]string, error) {
	rows, err := m.DB.Query("SELECT video_id FROM  youtubevideos")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var videoIDs []string

	for rows.Next() {
		var id string

		err := rows.Scan(&id)

		if err != nil {
			return videoIDs, err
		}

		videoIDs = append(videoIDs, id)
	}

	return videoIDs, err

}

func parseDownloadDate(s string) (string, error) {
	parsedTime, err := time.Parse(rfc3339Milli, s)

	if err != nil {
		return "", err
	}

	return parsedTime.UTC().String(), nil
}
