package helpers

import (
	"io"
	"log"
	"os"
)

// IOProgress represents I/O progress information.
type IOProgress struct {
	CompletedBytes  int64
	TotalBytes      int64
	PercentComplete int
}

// IOProgressFunc is a callback invoked to indicate progress.
type IOProgressFunc func(progress IOProgress)

// ProgressForReader wraps the specified io.Reader in an io.Reader that invokes progressFunc to report progress.
func ProgressForReader(innerReader io.Reader, progressFunc IOProgressFunc, totalBytes int64) io.Reader {
	progressReader := &progressReader{
		InnerReader:     innerReader,
		ProgressFunc:    progressFunc,
		TotalBytes:      totalBytes,
		progressChannel: make(chan IOProgress, 10),
		stopChannel:     make(chan bool, 1),
	}
	progressReader.startProgressPump()

	return progressReader
}

// ProgressForFileReader wraps the specified os.File in an io.Reader that invokes progressFunc to report progress.
func ProgressForFileReader(file *os.File, progressFunc IOProgressFunc) (io.Reader, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	progressReader := ProgressForReader(file, progressFunc, fileInfo.Size())

	return progressReader, nil
}

// progressReader wraps io.Reader, and invoking a callback to indicate progress.
type progressReader struct {
	InnerReader     io.Reader
	ProgressFunc    IOProgressFunc
	CompletedBytes  int64
	TotalBytes      int64
	PercentComplete int
	progressChannel chan IOProgress
	stopChannel     chan bool
}

var _ io.Reader = &progressReader{}

// Read reads up to len(buffer) bytes into buffer.
//
// It returns the number of bytes read and any error encountered.
func (reader *progressReader) Read(buffer []byte) (bytesRead int, err error) {
	bytesRead, err = reader.InnerReader.Read(buffer)
	if err != nil {
		reader.stopProgressPump(true)

		return
	}

	reader.CompletedBytes = reader.CompletedBytes + int64(bytesRead)

	// Handle the last x bytes, if required.
	if reader.CompletedBytes == reader.TotalBytes {
		reader.PercentComplete = 100
		reader.notifyProgress()

		return
	}

	// Don't bother reporting changes less than 10%.
	percentComplete := int(float64(100) * (float64(reader.CompletedBytes) / float64(reader.TotalBytes)))
	if percentComplete-reader.PercentComplete >= 10 {
		reader.PercentComplete = percentComplete
		reader.notifyProgress()
	}

	return
}

// Notify the
func (reader *progressReader) notifyProgress() {
	reader.progressChannel <- IOProgress{
		CompletedBytes:  reader.CompletedBytes,
		TotalBytes:      reader.TotalBytes,
		PercentComplete: reader.PercentComplete,
	}
}

// Start the reader's progress notification pump.
func (reader *progressReader) startProgressPump() {
	go reader.progressPump()
}

// Notify the progress pump that it should terminate.
func (reader *progressReader) stopProgressPump(dueToError bool) {
	reader.stopChannel <- dueToError
}

// Read from the progress channel, raising notifications as required.
func (reader *progressReader) progressPump() {
	log.Println("Progress pump started.")

Loop:
	for {
		select {
		case stoppedDueToError := <-reader.stopChannel:
			if stoppedDueToError {
				log.Println("Progress pump stopped due to error.")
			}

			break Loop

		case progress := <-reader.progressChannel:
			reader.ProgressFunc(progress)

			if progress.PercentComplete == 100 {
				reader.stopProgressPump(false)
			}
		}
	}

	log.Println("Progress pump terminated.")
}
