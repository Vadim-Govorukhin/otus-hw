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
	chunkNum, bufSize int64
}

func (t Task) String() string {
	return fmt.Sprintf("task %v to read %v bytes with offset %v",
		t.chunkNum, t.bufSize, t.chunkNum*t.bufSize)
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
		limit = fileSize - offset
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

func PrepareBufferLimit(limit int64) (int64, int) {
	bufSize := limit / 2
	if limit > 1024 { // Need to benchmark
		bufSize = 512 // Need to benchmark
	}
	iterNum := int(math.Ceil(float64(limit) / float64(bufSize)))
	return bufSize, iterNum
}

func makeCopyAsync(reader io.ReadSeeker, outputFile io.Writer, limit, offset int64) error {
	var mu sync.Mutex
	task := make(chan Task)
	errorsChan := make(chan error)
	defer close(errorsChan)

	bufSize, iterNum := PrepareBufferLimit(limit)
	chunks := make([]chan []byte, iterNum) // create a slice of channels for each reading iteration
	for i := 0; i < iterNum; i++ {
		chunks[i] = make(chan []byte, 1)
	}

	go func() {
		var t Task
		for i := 0; i < iterNum; i++ {
			t = Task{int64(i), bufSize}
			infoLog.Printf("[sender] send %s", t)
			task <- t
		}
		infoLog.Print("[sender] close channel task")
		close(task)
	}()

	nWorkers := 5
	for n := 0; n < nWorkers; n++ {
		go func(n int, reader io.ReadSeeker, task chan Task, chunk []chan []byte) {
			infoLog.Printf("[goroutine %v] start", n)
			defer infoLog.Printf("[goroutine %v] end", n)

			var buffer []byte
			for t := range task {
				infoLog.Printf("[goroutine %v] run %s", n, t)
				buffer = make([]byte, t.bufSize)

				mu.Lock()
				if _, err := reader.Seek(t.chunkNum*t.bufSize+offset, io.SeekStart); err != nil {
					mu.Unlock()
					errorLog.Printf("[goroutine %v] task %v stopped with error %e ",
						n, t.chunkNum, err)
					errorsChan <- err
					continue
				}

				n, err := reader.Read(buffer)
				mu.Unlock()
				if (err != nil) && !errors.Is(err, io.EOF) {
					errorLog.Printf("[goroutine %v] task %v stopped with error %e ",
						n, t.chunkNum, err)
					errorsChan <- err
					continue
				}

				infoLog.Printf("[goroutine %v] sending data from task %v", n, t.chunkNum)
				chunk[t.chunkNum] <- buffer[:n]
				close(chunk[t.chunkNum])
			}
		}(n, reader, task, chunks)
	}

	var err error
	for i := 0; i < iterNum; {
		select {
		case data := <-chunks[i]:
			i++
			infoLog.Println(ProgressBar(int64(i), int64(iterNum)))
			infoLog.Printf("[main] write data from chunkNum %v of %v", i, iterNum)
			outputFile.Write(data)
		case e := <-errorsChan:
			err = e
		}
	}
	return err
}

func makeCopySync(reader io.Reader, outputFile io.Writer, limit int64) error {
	bufSize, _ := PrepareBufferLimit(limit)
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
		if errors.Is(err, io.EOF) {
			outputFile.Write(buffer)
			return nil
		}
		if err != nil {
			errorLog.Println(err)
			return err
		}
		_, err = outputFile.Write(buffer)
		if err != nil {
			errorLog.Println(err)
			return err
		}
	}
	return nil
}

// Copy - copy fromPath file to toPath file with given offset and limit.
func Copy(fromPath, toPath string, offset, limit int64, isAsync bool) error {
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
	if limit == 0 {
		return nil
	}

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

	if isAsync {
		err = makeCopyAsync(inputFile, outputFile, limit, offset)
	} else {
		err = makeCopySync(inputFile, outputFile, limit)
	}
	if err != nil {
		errorLog.Println(err)
		return err
	}
	infoLog.Printf("Wrote data to new file %s", outputFile.Name())

	return nil
}
