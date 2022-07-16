package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

type Task struct {
	chunkNum          int
	currRead, bufSize int64
}

type copyChunc struct {
	chunkNum int
	data     []byte
}

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
	wg := sync.WaitGroup{}
	chunk := make(chan copyChunc)
	task := make(chan Task)

	nWorkers := 5
	go func(nWorkers int) {
		for n := 0; n < nWorkers; n++ {
			wg.Add(1)
			go func(reader io.Reader, task chan Task, n int) {
				infoLog.Printf("[goroutine %v] start", n)
				defer wg.Done()
				defer infoLog.Printf("[goroutine %v] end", n)

				var bufReader *bufio.Reader
				var buffer []byte

				for t := range task {
					infoLog.Printf("[goroutine %v] take task %v", n, t.chunkNum)
					bufReader = bufio.NewReader(reader) // creates a new reader
					if _, err := bufReader.Discard(int(t.currRead)); err != nil {
						errorLog.Println(err)
						return
					}

					buffer = make([]byte, t.bufSize)
					_, err := bufReader.Read(buffer)
					if err == io.EOF {
						chunk <- copyChunc{t.chunkNum, buffer}
						return
					}
					if err != nil {
						errorLog.Println(err)
						return
					}

					infoLog.Printf("[goroutine %v] sending data from task %v", n, t.chunkNum)
					chunk <- copyChunc{t.chunkNum, buffer}
				}

			}(reader, task, n)
		}
		wg.Wait()
		infoLog.Print("[sender tasks] close channel chunk")
		close(chunk)
	}(nWorkers)

	bufSize := PrepareBufferLimit(limit)
	var chunkNum int
	go func() {
		var currRead int64
		for ; currRead < limit; currRead += bufSize {
			infoLog.Printf("[sender tasks] send task %v to chanel", chunkNum)
			if bufSize > limit-currRead {
				bufSize = limit - currRead
			}
			task <- Task{chunkNum, currRead, bufSize}
			chunkNum++
		}
		infoLog.Print("[sender tasks] close channel task")
		close(task)
	}()

	var seekChunk int = 1
	chunkStore := make(map[int][]byte)
	go func() {
		for ch := range chunk {
			infoLog.Printf("[main] reveive data from task %v", ch.chunkNum)
			chunkStore[ch.chunkNum] = ch.data
		}
	}()

	for {
		if val, ok := chunkStore[seekChunk]; ok {
			infoLog.Println(ProgressBar(int64(seekChunk), int64(chunkNum)))
			infoLog.Printf("[main] write data from chunkNum %v of %v", seekChunk, chunkNum)
			outputFile.Write(val)
			if seekChunk == chunkNum {
				return nil
			}
			seekChunk++
		}
	}

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
