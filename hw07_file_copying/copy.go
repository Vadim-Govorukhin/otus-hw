package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Логгер для записи информационных сообщений.
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Логгер для записи сообщений об ошибках.
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

// GetFileSize - get size of input file.
func GetFileSize(fromPath string) (int64, error) {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		errorLog.Println(err)
		return 0, ErrUnsupportedFile
	}
	return fileInfo.Size(), nil
}

// CheckArgs - check given arguments.
func CheckArgs(fileSize, offset int64, limit *int64) error {
	if *limit > fileSize-offset {
		*limit = fileSize - offset
	}
	if *limit == 0 {
		*limit = fileSize
	}

	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

// ProgressBar - return simple progress bar string.
func ProgressBar(i, limit, readFileSize int64) string {
	minSize := readFileSize
	if limit < readFileSize {
		minSize = limit
	}
	percent := float64(i) / float64(minSize)
	return fmt.Sprintf("Completed %.2f%%", percent*100)
}

func PrepareBuffer(limit int64) ([]byte, int64) {
	bufLimit := limit / 2
	if limit > 512 { /////
		bufLimit = 64 /////
	}

	data := make([]byte, bufLimit)
	return data, bufLimit
}

func makeCopy(reader io.Reader, outputFile io.Writer, limit int64, progressBar func(int64) string) error {
	data, bufLimit := PrepareBuffer(limit)
	var err error
	for n := int64(0); n < limit; n += bufLimit {
		if bufLimit > limit-n {
			bufLimit = limit - n
			data = make([]byte, bufLimit)
		}
		infoLog.Println(progressBar(n))

		_, err = reader.Read(data)
		if err == io.EOF {
			outputFile.Write(data)
			break
		}
		if err != nil {
			errorLog.Println(err)
			return err
		}
		outputFile.Write(data)
	}
	return nil
}

// Copy - copy fromPath file to toPath file with given offset and limit.
func Copy(fromPath, toPath string, offset, limit int64) error {
	fileSize, err := GetFileSize(fromPath)
	if err != nil {
		errorLog.Println(err)
		return err
	}

	err = CheckArgs(fileSize, offset, &limit)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Input args checked")

	outputFile, err := os.Create(toPath)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Output file created")
	defer outputFile.Close()

	inputFile, err := os.Open(fromPath)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile) // creates a new reader
	_, err = reader.Discard(int(offset)) // discard the following offset bytes
	if err != nil {
		errorLog.Println(err)
		return err
	}

	progressBar := func(i int64) string { return ProgressBar(i, limit, fileSize-offset) }

	err = makeCopy(reader, outputFile, limit, progressBar)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Wrote data to new file")

	return nil
}
