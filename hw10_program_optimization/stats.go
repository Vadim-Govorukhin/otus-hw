package hw10programoptimization

import (
	"bufio"
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

type DomainStruct struct {
	stat DomainStat
	mu   sync.Mutex
}

func (d *DomainStruct) Add(key string) {
	d.mu.Lock()
	d.stat[key] += 1
	d.mu.Unlock()
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	nWorkers := 10
	var tasks = make(chan string)
	var usersCh = make(chan string)
	var errors = make(chan error)

	result := DomainStruct{stat: make(DomainStat)}
	//re, err := regexp.Compile("@(?P<Domain>\\w+\\." + domain + ")")
	re2, err := regexp.Compile("\"Email\":\".+@(?P<Domain>\\w+\\." + domain + ")")
	if err != nil {
		return nil, err
	}

	// Workers.
	var wg sync.WaitGroup
	wg.Add(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			defer infoLog.Printf("[goroutine %v] end\n", i)
			for task := range tasks {
				infoLog.Printf("[goroutine %v] take task '%s'\n", i, task)
				var user User
				if err := workForWorker(task, &user, re2, &result); err != nil {
					errorLog.Printf("[goroutine %v] exit task '%s' with error: %s", i, task, err)
					errors <- err
					break
				}
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
			infoLog.Println("[sender] sending task")
			tasks <- scanner.Text()
			//tasks <- scanner.Bytes()
		}
		infoLog.Println("[sender] send all tasks, close tasks")
	}()

	for err := range errors {
		return nil, err
	}

	return result.stat, nil
}

func workForWorker(task string, user *User, re *regexp.Regexp, result *DomainStruct) error {
	if email := re.FindStringSubmatch(task); email != nil {
		result.Add(strings.ToLower(email[1]))
	}
	/*
		err := json.Unmarshal([]byte(task), user)
		if err != nil {
			return err
		}

		if email := re.FindStringSubmatch(user.Email); email != nil {
			result.Add(strings.ToLower(email[1]))
		}

	*/
	return nil
}
