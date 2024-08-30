# syntax=docker/dockerfile:1

FROM golang:1.22.6

WORKDIR /app

RUN apt-get -y update

RUN apt-get install -y ffmpeg

COPY yt-dlp_linux ./

COPY test.txt ./

COPY test2.txt ./

COPY .env ./

# RUN ./yt-dlp_linux

# RUN pip3 install yt-dlp

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/Youtube-Downloader/main.go ./

# COPY internal/drop/drop.go ./

# COPY internal/sqlfuncs/sqlfuncs.go ./

# COPY internal/youtube/youtube.go ./

ADD internal/drop internal/drop

ADD internal/sqlfuncs internal/sqlfuncs

ADD internal/youtube internal/youtube

RUN go mod tidy

RUN CGO_ENABLED=1 GOOS=linux go build -o /youtube-downloader

CMD ["/youtube-downloader"]
