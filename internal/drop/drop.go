package drop

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
)

const singleShotUploadSizeCutoff int64 = 32 * (1 << 20)

type DropboxModel struct {
	user users.Client
	file files.Client
}

type uploadChunk struct {
	data   []byte
	offset uint64
	close  bool
}

// TODO: Rework all this code becase returning users.New doesn't make sense,
// I could return a users.New and files.new and put them in the same struct??????
// Have to see if thats worth it or not
func InitDropbox() DropboxModel {
	token := os.Getenv("ACCESS_TOKEN")
	config := dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogDebug, // if needed, set the desired logging level. Default is off
	}
	dbUser := users.New(config)
	dbFiles := files.New(config)

	return DropboxModel{
		user: dbUser,
		file: dbFiles,
	}

}

func (d DropboxModel) GetAccount() error {
	if resp, err := d.user.GetCurrentAccount(); err != nil {
		return err
	} else {
		fmt.Printf("Name: %v", resp.Name)
	}

	return nil
}

func (d DropboxModel) UploadFile(filepath *string) error {
	//Opens file

	file, err := os.Open(*filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	//Gets information about the file
	fileInfo, err := file.Stat()

	if err != nil {
		return err
	}

	//Sets the info for downloading, argument is path of file on dropbox
	commitInfo := files.NewCommitInfo("/yuppers'")

	//Set to overwrite file on dropbox, if uploading something that already exists
	commitInfo.Mode.Tag = "overwrite"

	// ts := time.Now().UTC().Round(time.Second)
	// commitInfo.ClientModified = &ts

	//arg is path to file thats being uploaded
	fileUploadArg := files.NewUploadArg("/test.txt")

	fileUploadArg.Mode.Tag = "overwrite"

	if fileInfo.Size() > singleShotUploadSizeCutoff {
		// return nil
		// return d.uploadLargeFileConcurrent(file, fileInfo.Size(), commitInfo)
		return d.uploadLargeFile(fileInfo.Size(), file, commitInfo)
	} else {
		res, err := d.file.Upload(fileUploadArg, file)

		fmt.Print(res)

		return err

	}
}

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
