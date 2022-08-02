package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"sync"
)

func GetDomainStatRegexp(r io.Reader, domain string) (DomainStat, error) {
	nWorkers := 10
	tasks := make(chan string)
	usersCh := make(chan string)
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
			resultSl[i] = make(DomainStat, 100)
			for task := range tasks {
				infoLog.Printf("[goroutine %v] take task '%s'\n", i, task)
				workForWorkerRegexp(task, re, &resultSl[i])
			}
		}(i)
	}

	go func() {
		defer func() {
			infoLog.Println("[sender] send all tasks, close tasks")
			close(tasks)
			wg.Wait()
			infoLog.Println("[sender] all goroutines stopped, close usersCh")
			close(usersCh)
			close(errors)
		}()

		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)

		infoLog.Println("[sender] start sending tasks")
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				errors <- err
				return
			}
			tasks <- scanner.Text()
		}
	}()

	for err := range errors {
		return nil, err
	}

	res := resultSl[0]
	for i := 1; i < nWorkers; i++ {
		for key, val := range resultSl[i] {
			res[key] += val
		}
	}
	return res, nil
}

func workForWorkerRegexp(task string, re *regexp.Regexp, result *DomainStat) {
	if email := re.FindStringSubmatch(task); email != nil {
		(*result)[strings.ToLower(email[1])]++
	}
}
