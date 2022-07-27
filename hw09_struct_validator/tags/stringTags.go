package tags

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type StringTags struct {
	len    int
	regexp string
	in     []string
}

func (T *StringTags) FillField(tag string) error {
	m := strings.Split(tag, ":")
	if len(m) > 2 {
		return ErrUnsupportedTag
	}

	switch m[0] {
	case "len":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			ErrorLog.Printf("parsing error %e", err)
			return err
		}
		T.len = i
	case "regexp":
		T.regexp = m[1]
	case "in":
		T.in = strings.Split(m[1], ",")
	default:
		ErrorLog.Printf("Unsupported tag name: %s\n", m[0])
		return ErrUnsupportedTag
	}
	return nil
}

func (T *StringTags) IsValid(i reflect.Value) error {
	val := i.String()

	// len
	if (T.len != 0) && (len(val) != T.len) {
		return ErrInvaildByTag
	}

	// regex
	InfoLog.Println("Start check regex")
	if T.regexp != "" {
		re, err := regexp.Compile(T.regexp)
		if err != nil {
			ErrorLog.Printf("regex error %e", err)
			return err
		}
		InfoLog.Println("Here 1")
		if ok := re.MatchString(T.regexp); !ok {
			return ErrInvaildByTag
		}
		InfoLog.Println("Here 2")
	}

	// in
	if len(T.in) != 0 {
		var flag bool
		for _, str := range T.in {
			if str == val {
				flag = true
				break
			}
		}
		if !flag {
			return ErrInvaildByTag
		}
	}

	return nil
}
