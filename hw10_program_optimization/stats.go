package hw10programoptimization

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // for info message

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
	var user User
	re := regexp.MustCompile("@(?P<Domain>\\w+\\." + domain + ")")

	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if email := re.FindStringSubmatch(user.Email); email != nil {
			result[strings.ToLower(email[1])]++
		}
	}

	return result, nil
}
