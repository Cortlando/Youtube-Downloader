package drop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cortlando/youtube-downloader/internal/youtube"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
	"golang.org/x/sync/errgroup"
)

// const singleShotUploadSizeCutoff int64 = 32 * (1 << 20)

const singleShotUploadSizeCutoff int64 = 157286400

type DropboxModel struct {
	user users.Client
	file files.Client
}

type uploadChunk struct {
	data   []byte
	offset uint64
	close  bool
}

type authResponse struct {
	Access_token string
	Token_type   string
	Expires_in   int32
}

// TODO: Rework all this code becase returning users.New doesn't make sense,
// I could return a users.New and files.new and put them in the same struct??????
// Have to see if thats worth it or not
func InitDropbox() DropboxModel {
	fmt.Printf("File size cut off is: %d", singleShotUploadSizeCutoff)

	token := GetNewToken()
	config := dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogDebug, // if needed, set the desired logging level. Default is off
	}
	dbUser := users.New(config)
	dbFiles := files.New(config)

	fmt.Println(dbUser)
	fmt.Println(dbFiles)

	return DropboxModel{
		user: dbUser,
		file: dbFiles,
	}

}

func GetNewToken() string {
	// API endpoint for token refresh
	authUrl := "https://api.dropbox.com/oauth2/token"

	// Create URL-encoded form data
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", os.Getenv("REFRESH_TOKEN"))
	formData.Set("client_id", os.Getenv("APP_KEY"))
	formData.Set("client_secret", os.Getenv("APP_SECRET"))

	// Create a new POST request
	req, err := http.NewRequest("POST", authUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %s", err.Error())

	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing request: %s", err.Error())

	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body %s", err.Error())

	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	var authRes authResponse

	err = json.Unmarshal(body, &authRes)

	if err != nil {
		fmt.Println("Error parsing JSON: ", err.Error())
		log.Fatalf("Error parsing JSON: %s", err.Error())

	}

	fmt.Println("Access Token:", authRes.Access_token)

	return authRes.Access_token
}

func (d DropboxModel) GetAccount() error {
	if resp, err := d.user.GetCurrentAccount(); err != nil {
		println("broke")
		return err
	} else {
		fmt.Printf("Name: %v", resp.Name)
	}

	return nil
}

// Figure out a way to track which videos get uploaded successfully, and which dont
// TODO:Refactor this code so that it opens the download directory, and then gets the file names that way
func (d DropboxModel) UploadFiles(downloadedVideos map[string]youtube.ExtractedVideoInfo) ([]youtube.ExtractedVideoInfo, []error) {

	fmt.Print("Starting UploadFiles \n")
	// wg := sync.WaitGroup{}
	g := errgroup.Group{}
	g.SetLimit(3)

	errCh := make(chan error, len(downloadedVideos))
	videoCh := make(chan youtube.ExtractedVideoInfo, len(downloadedVideos))

	var uploadedVideos []youtube.ExtractedVideoInfo
	var errorList []error

	//Use this to loop over the downloaded videos, instead of passing in an array
	c, err := os.ReadDir("./downloads")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range c {
		// wg.Add(1)

		info, _ := file.Info()

		titleNoExtension := strings.Split(info.Name(), ".")

		video := downloadedVideos[titleNoExtension[0]]

		// fmt.Println(info.Name(), downloadedVideos[i])

		pathToFile := fmt.Sprintf("./downloads/%s", file.Name())
		uploadPath := fmt.Sprintf("/%s.%s", video.Title, titleNoExtension[1])

		g.Go(func() error {
			fmt.Print("Starting UploadFile Singuler:")
			fmt.Print(pathToFile)
			fmt.Print("\n")
			file, err := os.Open(pathToFile)

			if err != nil {
				log.Fatal(err)
			}
			//Removes file when done uploading/trying to upload it
			defer os.Remove(pathToFile)
			defer file.Close()

			//Gets information about the file
			fileInfo, err := file.Stat()

			if err != nil {
				log.Fatal(err)
			}

			//Sets the info for downloading, argument is path of file on dropbox
			commitInfo := files.NewCommitInfo(uploadPath)

			//Set to overwrite file on dropbox, if uploading something that already exists
			commitInfo.Mode.Tag = "overwrite"

			//argument is path of file on dropbox
			fileUploadArg := files.NewUploadArg(uploadPath)

			fileUploadArg.Mode.Tag = "overwrite"

			if fileInfo.Size() > singleShotUploadSizeCutoff {
				err := d.uploadLargeFile(fileInfo.Size(), file, commitInfo)
				errCh <- err

				if err == nil {
					videoCh <- video
				}
				return err
			} else {
				res, err := d.file.Upload(fileUploadArg, file)

				//TODO:Delete the file after uploading it
				errCh <- err

				if err == nil {
					videoCh <- video
				}

				fmt.Println(res)

				fmt.Println("Finished uploading")

				return err

			}

		})

		// go d.UploadFile(&pathToFile, &uploadPath, &wg, errCh, videoCh, downloadedVideos[titleNoExtension[0]])
	}

	go func() {
		// wg.Wait()
		g.Wait()
		fmt.Println(len(errCh))
		fmt.Println(len(videoCh))
		close(videoCh)
		close(errCh)
	}()

	for e := range errCh {
		if e != nil {

			errorList = append(errorList, e)
		}

	}

	//Returns a slice of videos that downloaded successfully
	//I can pass this to the db and dropbox functions
	for v := range videoCh {

		uploadedVideos = append(uploadedVideos, v)

	}

	return uploadedVideos, errorList
}

// func (d DropboxModel) UploadFile(pathToFile *string, uploadPath *string, wg *sync.WaitGroup, errCh chan<- error, videoCh chan<- youtube.ExtractedVideoInfo, video youtube.ExtractedVideoInfo) error {
// 	//Opens file
// 	defer wg.Done()
// 	fmt.Print("Starting UploadFile Singuler:")
// 	fmt.Print(*pathToFile)
// 	fmt.Print("\n")
// 	file, err := os.Open(*pathToFile)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer file.Close()

// 	//Gets information about the file
// 	fileInfo, err := file.Stat()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	//Sets the info for downloading, argument is path of file on dropbox
// 	commitInfo := files.NewCommitInfo(*uploadPath)

// 	//Set to overwrite file on dropbox, if uploading something that already exists
// 	commitInfo.Mode.Tag = "overwrite"

// 	// ts := time.Now().UTC().Round(time.Second)
// 	// commitInfo.ClientModified = &ts

// 	//argument is path of file on dropbox
// 	fileUploadArg := files.NewUploadArg(*uploadPath)

// 	fileUploadArg.Mode.Tag = "overwrite"

// 	if fileInfo.Size() > singleShotUploadSizeCutoff {
// 		// return nil
// 		// return d.uploadLargeFileConcurrent(file, fileInfo.Size(), commitInfo)
// 		return d.uploadLargeFile(fileInfo.Size(), file, commitInfo)
// 	} else {
// 		res, err := d.file.Upload(fileUploadArg, file)

// 		errCh <- err

// 		if err == nil {
// 			videoCh <- video
// 		}

// 		fmt.Print(res)

// 		fmt.Print("Finished uploading")

// 		return err

// 	}
// }

// func (d DropboxModel) UploadFile(pathToFile *string, uploadPath *string, errCh chan<- error, videoCh chan<- youtube.ExtractedVideoInfo, video youtube.ExtractedVideoInfo) error {
// 	//Opens file

// 	fmt.Print("Starting UploadFile Singuler:")
// 	fmt.Print(*pathToFile)
// 	fmt.Print("\n")
// 	file, err := os.Open(*pathToFile)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer file.Close()

// 	//Gets information about the file
// 	fileInfo, err := file.Stat()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	//Sets the info for downloading, argument is path of file on dropbox
// 	commitInfo := files.NewCommitInfo(*uploadPath)

// 	//Set to overwrite file on dropbox, if uploading something that already exists
// 	commitInfo.Mode.Tag = "overwrite"

// 	// ts := time.Now().UTC().Round(time.Second)
// 	// commitInfo.ClientModified = &ts

// 	//argument is path of file on dropbox
// 	fileUploadArg := files.NewUploadArg(*uploadPath)

// 	fileUploadArg.Mode.Tag = "overwrite"

// 	if fileInfo.Size() > singleShotUploadSizeCutoff {
// 		// return nil
// 		// return d.uploadLargeFileConcurrent(file, fileInfo.Size(), commitInfo)
// 		return d.uploadLargeFile(fileInfo.Size(), file, commitInfo)
// 	} else {
// 		res, err := d.file.Upload(fileUploadArg, file)

// 		errCh <- err

// 		if err == nil {
// 			videoCh <- video
// 		}

// 		fmt.Print(res)

// 		fmt.Print("Finished uploading")

// 		return err

// 	}
// }

func (d DropboxModel) uploadLargeFile(sizeOfFile int64, file io.Reader, commitInfo *files.CommitInfo) error {
	//Size of data chucks being sent
	chunkSize := int64(4194304)
	//Sets upload session arguments
	uploadSessionStartArgs := files.NewUploadSessionStartArg()

	uploadSessionStartArgs.SessionType = &files.UploadSessionType{}

	uploadSessionStartArgs.SessionType.Tag = files.UploadSessionTypeSequential

	res, err := d.file.UploadSessionStart(uploadSessionStartArgs, nil)

	if err != nil {
		return err
	}

	var written = int64(0)

	// cursor := files.NewUploadSessionCursor(res.SessionId,)

	for written < sizeOfFile {
		//Reading 4mb sized chunks from the file
		fmt.Printf("Bytes Written so far: %d \n", written)
		data, err := io.ReadAll(&io.LimitedReader{R: file, N: chunkSize})

		if err != nil {
			return err
		}

		expectedLen := chunkSize

		if written+chunkSize > sizeOfFile {
			expectedLen = sizeOfFile - written
		}

		if len(data) != int(expectedLen) {
			return fmt.Errorf("failed to read %d bytes from source", expectedLen)
		}

		chunk := uploadChunk{
			data:   data,
			offset: uint64(written),
			//If reaching EOF, set close to true
			close: written+chunkSize >= sizeOfFile,
		}

		//////////////////////////////////////////////////////

		cursor := files.NewUploadSessionCursor(res.SessionId, chunk.offset)
		args := files.NewUploadSessionAppendArg(cursor)
		args.Close = chunk.close

		if err := d.uploadOneChunk(args, chunk.data); err != nil {
			return err
		}

		written += int64(len(data))

	}

	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)

	_, err = d.file.UploadSessionFinish(args, nil)

	if err != nil {
		return err
	}

	return nil
}

func (d DropboxModel) uploadOneChunk(args *files.UploadSessionAppendArg, data []byte) error {
	for {
		//Uploads another chunk
		err := d.file.UploadSessionAppendV2(args, bytes.NewReader(data))

		//If error is rate limit error, wait and try again
		if err != nil {
			switch errt := err.(type) {
			//Auth coming from auth dropbox package
			case auth.RateLimitAPIError:
				time.Sleep(time.Second * time.Duration(errt.RateLimitError.RetryAfter))
				continue
			default:
				return err
			}
		}
		return nil
	}

}
