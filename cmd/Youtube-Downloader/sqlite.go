package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"os"

	_ "github.com/mattn/go-sqlite3"
)

const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

type YoutubeVideo struct {
	Title         string
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
	_, err := m.DB.Exec(`INSERT INTO youtubevideos(title, video_id, webpage_url) values('a','a','a')`)

	return err
}
