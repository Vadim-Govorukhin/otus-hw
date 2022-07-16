package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

// CheckArgs - check given arguments and return new limit and error.
func CheckArgs(fileSize, offset, limit int64) (int64, error) {
	if fileSize < offset {
		return 0, ErrOffsetExceedsFileSize
	}
	if fileSize == offset {
		return 0, nil
	}

	switch {
	case limit > fileSize-offset:
		limit = fileSize - offset
	case limit == 0:
		limit = fileSize
	}

	return limit, nil
}

// ProgressBar - return simple progress bar string.
func ProgressBar(currRead, limit int64) string {
	if limit == 0 {
		return fmt.Sprintf("Completed %.2f%%", 100.)
	}
	percent := float64(currRead) / float64(limit)
	return fmt.Sprintf("Completed %.2f%%", percent*100)
}

func PrepareBufferLimit(limit int64) int64 {
	bufLimit := limit / 2
	if limit > 1024 { // Need to benchmark
		bufLimit = 512 // Need to benchmark
	}

	return bufLimit
}

func makeCopy(reader io.Reader, outputFile io.Writer, limit int64) error {
	bufSize := PrepareBufferLimit(limit)
	buffer := make([]byte, bufSize)

	var err error
	var currRead int64
	defer func() { infoLog.Println(ProgressBar(currRead, limit)) }()

	for ; currRead < limit; currRead += bufSize {
		if bufSize > limit-currRead {
			bufSize = limit - currRead
			buffer = make([]byte, bufSize)
		}
		infoLog.Println(ProgressBar(currRead, limit))

		_, err = reader.Read(buffer)
		if err == io.EOF {
			outputFile.Write(buffer)
			return nil
		}
		if err != nil {
			errorLog.Println(err)
			return err
		}
		outputFile.Write(buffer)
	}
	return nil
}

// Copy - copy fromPath file to toPath file with given offset and limit.
func Copy(fromPath, toPath string, offset, limit int64) error {
	// Get file size
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		errorLog.Println(err)
		return ErrUnsupportedFile
	}
	fileSize := fileInfo.Size()

	// Check input arguments and change limit if needed
	limit, err = CheckArgs(fileSize, offset, limit)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Input args is checked, limit = %v", limit)

	outputFile, err := os.Create(toPath)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Output file is created")
	defer outputFile.Close()

	inputFile, err := os.Open(fromPath)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	defer inputFile.Close()

	_, err = inputFile.Seek(offset, 0) // discard the following offset bytes
	if err != nil {
		errorLog.Println(err)
		return err
	}

	err = makeCopy(inputFile, outputFile, limit)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Wrote data to new file %s", outputFile.Name())

	return nil
}
