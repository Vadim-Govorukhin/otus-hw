package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	fmt.Println("====================== NEW =====================")
	done := make(chan struct{})

	closeTask := make(chan struct{})
	defer close(closeTask)
	errors := make(chan error)
	defer close(errors)

	var wg sync.WaitGroup
	defer wg.Wait()

	task := make(chan func() error)
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
						return
					}
				}
			}
		}(i) // i для отлаживания
	}

	// goroutine for catching errors
	var doneFlag bool
	go func() {
		var curErrorNum int
		for range errors {
			curErrorNum++
			if !doneFlag && (curErrorNum == m) {
				fmt.Println("[done goroutine] limit m, close chanel")
				doneFlag = true
				close(done) //but continue receive errors from still working workers
			}
		}
	}()

	// goroutine for catching stopped goroutines
	go func() {
		var curDoneTaskNum int
		for range closeTask {
			curDoneTaskNum++
			if !doneFlag && (curDoneTaskNum == n) {
				fmt.Println("[done goroutine] none of tasks is running, close chanel")
				doneFlag = true
				close(done)
				return
			}
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

	close(done)
	doneFlag = true
	fmt.Println("[main] all tasks is done")
	return nil
}
