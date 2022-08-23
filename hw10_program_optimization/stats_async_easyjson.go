package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	easyjson "github.com/mailru/easyjson"
)

func GetDomainStatEasyjson(r io.Reader, domain string) (DomainStat, error) {
	nWorkers := 10
	tasks := make(chan []byte)
	errors := make(chan error)

	re := regexp.MustCompile("@(?P<Domain>\\w+\\." + domain + ")")

	// Workers.
	resultSl := make([]DomainStat, nWorkers)
	var wg sync.WaitGroup
	wg.Add(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			defer infoLog.Printf("[goroutine %v] end\n", i)
			resultSl[i] = make(DomainStat, 10)
			var user User
			for task := range tasks {
				infoLog.Printf("[goroutine %v] take task '%s'\n", i, task)
				if err := workForWorkerEasyjson(task, &user, re, &resultSl[i]); err != nil {
					errors <- fmt.Errorf("[goroutine %v] exit task '%s' with error: %w", i, task, err)
					break
				}
			}
		}(i)
	}

	go func() {
		defer func() {
			infoLog.Println("[sender] send all tasks, close tasks")
			close(tasks)
			wg.Wait()
			infoLog.Println("[sender] all goroutines stopped, close errors")
			close(errors)
		}()

		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)

		infoLog.Println("[sender] start sending tasks")
		var currLine []byte
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				errors <- err
				return
			}
			currLine = scanner.Bytes()
			line := make([]byte, len(currLine))
			copy(line, currLine)
			tasks <- line
		}
	}()

	for err := range errors {
		return nil, err
	}

	result := resultSl[0]
	for i := 1; i < nWorkers; i++ {
		for key, val := range resultSl[i] {
			result[key] += val
		}
	}

	return result, nil
}

func workForWorkerEasyjson(task []byte, user *User, re *regexp.Regexp, result *DomainStat) error {
	if err := easyjson.Unmarshal(task, user); err != nil {
		return err
	}

	if email := re.FindStringSubmatch(user.Email); email != nil {
		(*result)[strings.ToLower(email[1])]++
	}

	return nil
}
