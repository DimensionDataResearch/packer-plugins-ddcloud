package helpers

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/andelf/go-curl"
)

// MinProgressPercentChange is the minimum change in percentage completion to report.
var MinProgressPercentChange = 5

// UploadProgressFunc is a call-back function that receives progress information.
//
// Return true to continue, or false to cancel the upload.
type UploadProgressFunc func(fileName string, percentComplete int, currentBytes int64, totalBytes int64) bool

// Upload handles the process of uploading files to CloudControl.
type Upload struct {
	// The name of the FTPS host to which files will be uploaded.
	Host string

	// The FTP username.
	Username string

	// The FTP password.
	Password string

	// An optional callback that receives progress information.
	ProgressFunc UploadProgressFunc
}

// File uploads a file to Cloud Control.
func (upload *Upload) File(file *os.File) error {
	ensureCurlInitialized()

	var requestExecutionError error

	request, err := upload.createFileRequest(file, &requestExecutionError)
	if err != nil {
		return err
	}
	defer request.Cleanup()

	err = request.Perform()
	if err != nil {
		return err
	}
	if requestExecutionError != nil {
		return fmt.Errorf("Request execution error: %s", requestExecutionError)
	}

	return nil
}

// Create a file upload request.
//
// requestExecutionError will receive an error if the request's read callback encounters an error.
func (upload *Upload) createFileRequest(sourceFile *os.File, requestExecutionError *error) (request *curl.CURL, err error) {
	var sourceFileInfo os.FileInfo
	sourceFileInfo, err = sourceFile.Stat()
	if err != nil {
		return
	}

	targetFileName := path.Base(sourceFile.Name())
	request = curl.EasyInit()
	err = request.Setopt(curl.OPT_URL, fmt.Sprintf("ftp://%s/%s", upload.Host, targetFileName))
	if err != nil {
		return
	}

	err = request.Setopt(curl.OPT_UPLOAD, true)
	if err != nil {
		return
	}

	err = request.Setopt(curl.OPT_USE_SSL, true)
	if err != nil {
		return
	}

	err = request.Setopt(curl.OPT_USERNAME, upload.Username)
	if err != nil {
		return
	}

	err = request.Setopt(curl.OPT_PASSWORD, upload.Password)
	if err != nil {
		return
	}

	err = request.Setopt(curl.OPT_INFILESIZE, sourceFileInfo.Size())
	if err != nil {
		return
	}

	err = request.Setopt(curl.OPT_READDATA, sourceFile)
	if err != nil {
		return
	}

	err = request.Setopt(curl.OPT_READFUNCTION, func(data []byte, userdata interface{}) int {
		source := userdata.(*os.File)
		bytesRead, err := source.Read(data)
		if err != nil && err.Error() != "EOF" {
			*requestExecutionError = err

			return 0
		}

		return bytesRead // When this is 0, we're done.
	})
	if err != nil {
		return nil, err
	}

	// Progress.
	progressFunc := upload.ProgressFunc
	if progressFunc == nil {
		return
	}

	err = request.Setopt(curl.OPT_NOPROGRESS, false)
	if err != nil {
		return
	}

	var currentPercentComplete int
	err = request.Setopt(curl.OPT_PROGRESSFUNCTION, func(dltotal, dlnow, ultotal, ulnow float64, userdata interface{}) bool {
		if ultotal <= 0 {
			return true
		}

		// Only show a change of MinProgressPercentChange or more.
		percentComplete := int((ulnow / ultotal) * 100)
		if percentComplete-currentPercentComplete < MinProgressPercentChange {
			return true
		}
		currentPercentComplete = percentComplete

		currentBytes := int64(ulnow)
		totalBytes := int64(ultotal)

		return progressFunc(targetFileName, percentComplete, currentBytes, totalBytes)
	})
	if err != nil {
		return
	}

	return
}

// One-time curl setup.
var curlInitializer sync.Once

func ensureCurlInitialized() {
	curlInitializer.Do(func() {
		err := curl.GlobalInit(curl.GLOBAL_DEFAULT)
		if err != nil {
			panic(err)
		}
	})
}
