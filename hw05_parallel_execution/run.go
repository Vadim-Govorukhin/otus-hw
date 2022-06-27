package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	overLimit := make(chan struct{})
	//defer close(overLimit)

	task := make(chan func() error)
	defer close(task)

	errors := make(chan error)
	defer close(errors)

	var wg sync.WaitGroup
	defer wg.Wait()

	var ignoreErrors bool
	if m <= 0 {
		ignoreErrors = true
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case <-overLimit:
					fmt.Printf("\t[goroutine %d] end by chanel 'overlimit' \n", i)
					return
				case t, ok := <-task:
					if !ok {
						fmt.Printf("\t[goroutine %d] end by chanel 'task' \n", i)
						return
					}
					fmt.Printf("[goroutine %d] run task\n", i)
					err := t()
					if err != nil {
						fmt.Printf("[goroutine %d] end with error: %s\n", i, err)
						errors <- err
						fmt.Printf("[goroutine %d] send error: %s\n", i, err)
						if !ignoreErrors {
							return
						}
					}
					fmt.Printf("\t[goroutine %d] end task\n", i)
				}
			}
		}(i) // i для отлаживания
	}

	go func(m, n int) {
		var curErrorNum int
		for range errors {
			fmt.Println("[overlimit goroutine] receive 1 error")
			curErrorNum++
			if curErrorNum == m {
				fmt.Println("[overlimit goroutine] limit m, close chanel")
				close(overLimit)
			}
		}
		fmt.Println("[overlimit goroutine] not limit m, return")
	}(m, n)

	fmt.Println("[main] start sending tasks")
	var sendTaski int
	for {
		select {
		case <-overLimit:
			return ErrErrorsLimitExceeded
		case task <- tasks[sendTaski]:
			fmt.Println("[main] send task")
			sendTaski++
			if sendTaski == len(tasks) {
				close(overLimit)
				fmt.Println("[main] all tasks is done")
				return nil
			}
		}
	}
}
