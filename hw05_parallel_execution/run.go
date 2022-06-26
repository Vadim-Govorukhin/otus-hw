package hw05parallelexecution

import (
	"errors"
	"fmt"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	var overLimit chan struct{}
	defer close(overLimit)

	var task chan func() error
	defer close(task)

	var errors chan error
	defer close(errors)

	fmt.Println("[main] rum goroutines")
	for i := 0; i < n; i++ {
		go func(i int) {
			fmt.Printf("[goroutine %d] start\n", i)
			for {
				select {
				case _, ok := <-overLimit:
					if !ok {
						fmt.Printf("[goroutine %d] end with ErrErrorsLimitExceeded\n", i)
						return
					}
					fmt.Printf("[goroutine %d] wtf\n", i)
				case t := <-task:
					fmt.Printf("[goroutine %d] run task\n", i)
					err := t()
					if err != nil {
						fmt.Printf("[goroutine %d] end with error: %s\n", i, err)
						errors <- err
						return // Надо ли?
					}
				}
			}
		}(i) // i для отлажевания
	}

	go func(m int) {
		fmt.Println("[overlimit goroutine] starts")
		var curErrorNum int
		for _ = range errors {
			fmt.Println("[overlimit goroutine] receive 1 error")
			curErrorNum++
			if curErrorNum == m {
				fmt.Println("[overlimit goroutine] limit m, close chanel")
				close(overLimit)
				return
			}
		}
		fmt.Println("[overlimit goroutine] not limit m, return")
	}(m)

	fmt.Println("[main] start sending tasks")
	for _, t := range tasks {
		select {
		case task <- t:
			continue
		case _, ok := <-overLimit:
			if !ok {
				return ErrErrorsLimitExceeded
			}
		}
	}
	fmt.Println("[main] all tasks is done")
	return nil
}
