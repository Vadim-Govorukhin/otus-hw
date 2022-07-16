package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
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

func makeCopy(reader io.ReadSeeker, outputFile io.Writer, limit, offset int64) error {
	if limit == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	var mu sync.Mutex
	task := make(chan Task)

	bufSize := PrepareBufferLimit(limit)
	iterNum := int(math.Ceil(float64(limit) / float64(bufSize)))
	chunk := make([]chan []byte, iterNum)
	for i := 0; i < iterNum; i++ {
		chunk[i] = make(chan []byte, 1)
	}

	var chunkNum int
	go func() {
		var currRead int64
		for ; currRead < limit; currRead += bufSize {
			if bufSize > limit-currRead {
				bufSize = limit - currRead
			}
			infoLog.Printf("[sender] send task %v to read %v bytes with offset %v", chunkNum, bufSize, currRead)
			task <- Task{chunkNum, currRead, bufSize}
			chunkNum++
		}
		infoLog.Print("[sender] close channel task")
		close(task)
	}()

	nWorkers := 5
	go func(nWorkers int) {
		for n := 0; n < nWorkers; n++ {
			wg.Add(1)
			go func(reader io.ReadSeeker, task chan Task, chunk []chan []byte, n int) {
				infoLog.Printf("[goroutine %v] start", n)
				defer wg.Done()
				defer infoLog.Printf("[goroutine %v] end", n)

				var buffer []byte
				for t := range task {
					infoLog.Printf("[goroutine %v] task %v to read %v bytes with offset %v", n, t.chunkNum, t.bufSize, t.currRead)
					buffer = make([]byte, t.bufSize)

					mu.Lock()
					if _, err := reader.Seek(t.currRead+offset, io.SeekStart); err != nil {
						errorLog.Println(err)
						return
					}
					_, err := reader.Read(buffer)
					mu.Unlock()

					if err == io.EOF {
						infoLog.Printf("[goroutine %v] sending data from task %v", n, t.chunkNum)
						chunk[t.chunkNum] <- buffer
						close(chunk[t.chunkNum])
						return
					}
					if err != nil {
						errorLog.Println(err)
						return
					}

					infoLog.Printf("[goroutine %v] sending data from task %v", n, t.chunkNum)
					chunk[t.chunkNum] <- buffer
					close(chunk[t.chunkNum])
				}

			}(reader, task, chunk, n)
		}
		wg.Wait()
		infoLog.Print("[manager goroutines] close channel chunk")
	}(nWorkers)

	for i := 0; i < iterNum; {
		data := <-chunk[i]
		i++
		infoLog.Println(ProgressBar(int64(i), int64(iterNum)))
		infoLog.Printf("[main] write data from chunkNum %v of %v", i, iterNum)
		outputFile.Write(data)
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

	if _, err = inputFile.Seek(offset, 0); err != nil {
		errorLog.Println(err)
		return err
	}

	err = makeCopy(inputFile, outputFile, limit, offset)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Wrote data to new file %s", outputFile.Name())

	return nil
}
