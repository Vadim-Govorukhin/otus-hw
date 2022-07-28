package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

type users [100_000]User

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	nWorkers := 5
	var tasks = make(chan string)

	var usersCh = make(chan User, 2*nWorkers)

	var errors = make(chan error)
	defer close(errors)

	var wg sync.WaitGroup
	// Workers.
	wg.Add(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			defer infoLog.Printf("[goroutine %v] end\n", i)
			infoLog.Printf("[goroutine %v] start\n", i)
			var err error
			for task := range tasks {
				infoLog.Printf("[goroutine %v] take task '%s'\n", i, task)
				var user User
				if err = json.Unmarshal([]byte(task), &user); err != nil {
					errorLog.Printf("[goroutine %v] exit with error: %s", i, err)
					errors <- err
					break
				}
				usersCh <- user
				infoLog.Printf("[goroutine %v] send user", i)
			}

		}(i)
	}

	go func() {
		content, err := ioutil.ReadAll(r)
		if err != nil {
			errors <- err
			return
		}
		infoLog.Println("[sender] start sending tasks")
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			infoLog.Println("[sender] sending task")
			tasks <- line
		}
		infoLog.Println("[sender] send all tasks, close tasks")
		close(tasks)
		wg.Wait()
		infoLog.Println("[sender] all goroutines stopped, close usersCh")
		close(usersCh)
	}()

	var u users
	var i int
	for {
		select {
		case err := <-errors:
			errorLog.Printf("[main] get error %s", err)
			close(tasks)
			return nil, fmt.Errorf("get users error: %w", err)
		case user, ok := <-usersCh:
			if !ok {
				infoLog.Print("[main] received all users")
				return countDomains(u, domain)
			}
			u[i] = user
			i++
		}
	}
	/*
			u, err := getUsers(r)
		if err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}
		return countDomains(u, domain)
	*/
}

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
