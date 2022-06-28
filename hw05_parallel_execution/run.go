package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func handleDoneTasks(curDoneTaskNum, n int, doneFlag bool, done *chan struct{}) (int, bool) {
	curDoneTaskNum++
	if !doneFlag && (curDoneTaskNum == n) {
		fmt.Println("[doneTask goroutine] none of tasks is running, close chanel 'done'")
		doneFlag = true
		close(*done)
	}
	return curDoneTaskNum, doneFlag
}

func handleErrors(curErrorNum, m int, doneFlag bool, done *chan struct{}) (int, bool) {
	curErrorNum++
	if !doneFlag && (curErrorNum == m) {
		fmt.Println("[handleError goroutine] limit m, close chanel 'done'")
		doneFlag = true
		close(*done) // but continue receive errors from still working workers
	}
	return curErrorNum, doneFlag
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
			for {
				t, ok := <-task
				if !ok {
					fmt.Printf("\t[goroutine %d] end by chanel 'task' \n", i)
					return
				}
				fmt.Printf("[goroutine %d] run task\n", i)
				err := t()
				if err != nil {
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
		var curErrorNum int
		for range errors {
			mu.Lock() // remove race
			curErrorNum, doneFlag = handleErrors(curErrorNum, m, doneFlag, &done)
			mu.Unlock()
		}
	}()

	// goroutine for handling falling goroutines
	go func() {
		var curDoneTaskNum int
		for range closeTask {
			mu.Lock() // remove race
			curDoneTaskNum, doneFlag = handleDoneTasks(curDoneTaskNum, n, doneFlag, &done)
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
