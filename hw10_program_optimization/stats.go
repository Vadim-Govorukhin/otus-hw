package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

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
	nWorkers := 5
	var tasks = make(chan string)

	var usersCh = make(chan string)

	var errors = make(chan error)
	defer close(errors)

	var wg sync.WaitGroup
	// Workers.
	wg.Add(nWorkers)
	//var p fastjson.Parser
	for i := 0; i < nWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			defer infoLog.Printf("[goroutine %v] end\n", i)
			var err error
			for task := range tasks {
				infoLog.Printf("[goroutine %v] take task '%s'\n", i, task)
				var user User
				if err = workForWorker(task, &user); err != nil {
					errorLog.Printf("[goroutine %v] exit task '%s' with error: %s", i, task, err)
					errors <- err
					break
				}
				usersCh <- user.Email
				/*
					email, err := workForWorkerTest(task, p)
					if err != nil {
						errorLog.Printf("[goroutine %v] exit task '%s' with error: %s", i, task, err)
						errors <- err
						break
					}
					usersCh <- email
				*/
				infoLog.Printf("[goroutine %v] send user", i)
			}

		}(i)
	}

	go func() {
		defer func() {
			close(tasks)
			wg.Wait()
			infoLog.Println("[sender] all goroutines stopped, close usersCh")
			close(usersCh)
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
			infoLog.Println("[sender] sending task")
			tasks <- scanner.Text()
		}
		infoLog.Println("[sender] send all tasks, close tasks")
	}()

	result := make(DomainStat)
	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err

	}
	for {
		select {
		case err := <-errors:
			return nil, err
		case userEmail, ok := <-usersCh:
			if !ok {
				return result, nil
			}
			if matched := re.Match([]byte(userEmail)); matched {
				result[strings.ToLower(strings.SplitN(userEmail, "@", 2)[1])] += 1
			}
		}
	}
}

func workForWorker(task string, user *User) error {
	return json.Unmarshal([]byte(task), user)
}

/*
func workForWorkerTest(task string, p fastjson.Parser) (string, error) {
	v, err := p.Parse(task)
	return string(v.GetStringBytes("Email")), err
}
*/
