package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"os"

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

func (m YoutubeVideoModel) CheckIfDBFIleExists() {
	if _, err := os.Stat(DB_PATH); err == nil {
		fmt.Println("DB File exists")
	} else if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create("downloadedfiles.db")
		defer file.Close()

		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("This shouldn't happen")
	}

}

func (m YoutubeVideoModel) createYoutubeVideoTableIfNotExist() error {

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

func (m YoutubeVideoModel) testInsertIntoTable() error {
	_, err := m.DB.Exec(`INSERT INTO youtubevideos(title, video_id, webpage_url) values('b','b','b')`)

	return err
}

func (m YoutubeVideoModel) testSelectFromTable() ([]YoutubeVideo, error) {
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

func parseDownloadDate(s string) (string, error) {
	parsedTime, err := time.Parse(rfc3339Milli, s)

	if err != nil {
		return "", err
	}

	return parsedTime.UTC().String(), nil
}
