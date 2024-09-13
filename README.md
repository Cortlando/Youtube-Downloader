# Youtube Downloader

Youtube Downloader is a program I am developing to automatically backup videos from my personal youtube playlist daily. The program downloads the videos from Youtube using [yt-dlp](https://github.com/yt-dlp/yt-dlp), and uploads them to my dropbox account. It keeps track of which videos have been downloaded with a SQLite DB. 


## Features

- Downloads videos from a playlist on Youtube
- Uploads the videos to dropbox
- Keeps track of which videos have been downloaded already



## Environment Variables

To run this project, you will need to add the following environment variables to your .env file. To get the `APP_KEY`, and `APP_SECRET` you have to create a dropbox app [here for details](https://www.dropbox.com/developers/reference/getting-started#app%20console). To the the `REFRESH_TOKEN` follow this [guide](https://www.dropboxforum.com/t5/Dropbox-API-Support-Feedback/Get-refresh-token-from-access-token/td-p/596739) up to step 5

`URL`: The url of the youtube playlist you want to download

`APP_KEY`: 

`REFRESH_TOKEN`:

`APP_SECRET`:

`DB_PATH_ARGS`: "./db/downloadedfiles.db?_journal=WAL&_timeout=5000"

`DB_PATH`: "./db/downloadedfiles.db"
## Run Locally

Clone the project

```bash
  git clone https://github.com/Cortlando/Youtube-Downloader.git
```

Go to the project directory

Install dependencies

```bash
  go mod download
```

Build the program

```bash
  CGO_ENABLED=1 GOOS=linux go build -o /youtube-downloader
```

Run the executable 

- If getting error about yt-dlp, you'll have to install it manually [link](https://github.com/yt-dlp/yt-dlp)
- If getting this error "exec: "gcc": executable file not found in %PATH% when trying go build" check this [link](https://stackoverflow.com/questions/43580131/exec-gcc-executable-file-not-found-in-path-when-trying-go-build)
## Todo
- Add volume to docker container for persistent storage
- Add the ability to check playlists from soundcloud
- Improve the error handling
- Have the program email me a log of which files were downloaded and uploaded successfully
- Add other storage providers for uploading videos(google drive, S3, etc...)
- Add Tests
