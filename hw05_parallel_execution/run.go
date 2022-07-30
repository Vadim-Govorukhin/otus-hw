package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func handleDoneWorkers(doneWorkersNum, n int, doneFlag bool, done *chan struct{}) (int, bool) {
	doneWorkersNum++
	if !doneFlag && (doneWorkersNum == n) {
		fmt.Println("[doneworkers goroutine] none of workers is running, close chanel 'done'")
		doneFlag = true
		close(*done)
	}
	return doneWorkersNum, doneFlag
}

func handleErrors(errorsNum, m int, doneFlag bool, done *chan struct{}) (int, bool) {
	errorsNum++
	if !doneFlag && (errorsNum == m) {
		fmt.Println("[handleError goroutine] limit m, close chanel 'done'")
		doneFlag = true
		close(*done) // but continue receive errors from still working workers
	}
	return errorsNum, doneFlag
}

func Run(tasks []Task, n, m int) error {
	fmt.Println("====================== NEW =====================")
	done := make(chan struct{}) // channel for stop work

	closeTask := make(chan struct{}) // channel for handle fallen workers
	defer close(closeTask)
	errors := make(chan error) // channel for handle errors
	defer close(errors)

	var wg sync.WaitGroup
	defer wg.Wait()

	task := make(chan func() error) // channel for giving tasks to workers
	defer close(task)

	var ignoreErrors bool
	if m <= 0 {
		// workers continue working in case error
		ignoreErrors = true
	}

	// create n workers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer func() { closeTask <- struct{}{} }()
			for t := range task {
				fmt.Printf("[goroutine %d] run task\n", i)
				if err := t(); err != nil {
					errors <- err
					fmt.Printf("[goroutine %d] send error: %s\n", i, err)
					if !ignoreErrors {
						fmt.Printf("\t[goroutine %d] end by error \n", i)
						return
					}
				}
			}
		}(i)
	}

	// goroutine for handling errors
	var doneFlag bool // flag of close channel 'done' or not
	var mu sync.Mutex
	go func() {
		var errorsNum int
		for range errors {
			mu.Lock() // remove race
			errorsNum, doneFlag = handleErrors(errorsNum, m, doneFlag, &done)
			mu.Unlock()
		}
	}()

	// goroutine for handling falling goroutines
	go func() {
		var doneWorkersNum int
		for range closeTask {
			mu.Lock() // remove race
			doneWorkersNum, doneFlag = handleDoneWorkers(doneWorkersNum, n, doneFlag, &done)
			mu.Unlock()
		}
	}()

	for i := 0; i < len(tasks); {
		select {
		case <-done:
			return ErrErrorsLimitExceeded
		case task <- tasks[i]:
			fmt.Println("[main] send task")
			i++
		}
	}

	mu.Lock() // remove race
	doneFlag = true
	mu.Unlock()
	close(done)
	fmt.Println("[main] all tasks is done")
	return nil
}
