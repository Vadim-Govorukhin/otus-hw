package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	easyjson "github.com/mailru/easyjson"
)

var (
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // for info message
)

//easyjson:json
type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	nWorkers := 10
	var tasks = make(chan string)
	var errors = make(chan error)

	//result := DomainStruct{stat: make(DomainStat)}
	re := regexp.MustCompile("@(?P<Domain>\\w+\\." + domain + ")")
	//re := regexp.MustCompile(`"Email":".+@(?P<Domain>\w+\.` + domain + ")")
	resultSl := make([]DomainStat, nWorkers)
	// Workers.
	var wg sync.WaitGroup
	wg.Add(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			defer infoLog.Printf("[goroutine %v] end\n", i)
			resultSl[i] = make(DomainStat, 100)
			for task := range tasks {
				infoLog.Printf("[goroutine %v] take task '%s'\n", i, task)
				var user User
				if err := workForWorker(task, &user, re, &resultSl[i]); err != nil {
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

		var err error
		infoLog.Println("[sender] start sending tasks")
		for scanner.Scan() {
			err = scanner.Err()
			if err != nil {
				errors <- err
				return
			}
			tasks <- scanner.Text()
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

func workForWorker(task string, user *User, re *regexp.Regexp, result *DomainStat) error {
	// if email := re.FindStringSubmatch(task); email != nil {
	// 	//(*result)[strings.ToLower(email[1])] += 1
	// }

	err := easyjson.Unmarshal([]byte(task), user)
	if err != nil {
		return err
	}

	if email := re.FindStringSubmatch(user.Email); email != nil {
		(*result)[strings.ToLower(email[1])] += 1
	}

	return nil
}
