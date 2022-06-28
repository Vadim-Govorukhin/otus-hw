package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	fmt.Println("====================== NEW =====================")
	done := make(chan struct{})

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
	runningGoroutines := int32(n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer func() { atomic.AddInt32(&runningGoroutines, -1) }()
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
	var flag bool
	go func() {
		var curErrorNum int
		for range errors {

			fmt.Println("[done goroutine] error +1")
			curErrorNum++
			if !flag && ((curErrorNum == m) || (runningGoroutines == 0)) {
				fmt.Println("[done goroutine] limit m, close chanel")
				flag = true
				close(done) //but continue receive errors from still working workers
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
	flag = true
	fmt.Println("[main] all tasks is done")
	return nil
}
