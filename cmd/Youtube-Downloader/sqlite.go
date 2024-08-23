package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"os"

	_ "github.com/mattn/go-sqlite3"
)

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

	// tx, err := m.DB.BeginTx(context.TODO(), nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer tx.Rollback()
	// stmt := `CREATE TABLE IF NOT EXISTS youtubevideos(
	// 	Title TEXT,
	// 	video_id TEXT PRIMARY KEY,
	// 	 webpage_url TEXT,
	// 	 download_date TEXT
	// 	 )`
	// i, err := tx.ExecContext(context.TODO(), stmt)

	// println(i)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tx.Commit()

	// return err
	stmt := `CREATE TABLE IF NOT EXISTS youtubevideos(
		Title TEXT,
		video_id TEXT PRIMARY KEY,
		 webpage_url TEXT,
		 download_date TEXT
		 )`
	i, err := m.DB.Exec(stmt)
	println(i)
	return err
}

func (m YoutubeVideoModel) testInsertIntoTable() error {
	_, err := m.DB.Exec(`INSERT INTO youtubevideos(title, video_id, webpage_url, download_date) values('a','a','a','a')`)

	return err
}

// func (m YoutubeVideoModel) All() ([]extractedVideoInfo, error) {
// 	rows, err := m.DB.Query("SELECT title, id, webpageurl FROM YOUTUBEVIDEOS")

// 	return rows, err
// }

// func (m YoutubeVideoModel) ConnectDB() error {
// 	db, err := sql.Open("sqlite3", "./test.db")
// 	defer db.Close()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err = db.Ping(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return db.Ping()

// }
