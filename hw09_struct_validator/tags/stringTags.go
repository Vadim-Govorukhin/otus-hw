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
		return ErrTagInvalidSyntax
	}

	switch m[0] {
	case "len":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			errorLog.Printf("strconv.Atoi error %s", err)
			return ErrTagInvalidSyntax
		}
		T.len = i
	case "regexp":
		T.regexp = m[1]
	case "in":
		T.in = strings.Split(m[1], ",")
	default:
		errorLog.Printf("unsupported tag name: %s\n", m[0])
		return ErrUnsupportedTag
	}
	return nil
}

func (T *StringTags) ValidateValue(i reflect.Value) error {
	infoLog.Printf("\tvalidate value '%v' with tags %+v\n", i, *T)
	val := i.String()

	// len
	if (T.len != 0) && (len(val) != T.len) {
		return ErrInvaildByTag
	}

	// regex
	if T.regexp != "" {
		re, err := regexp.Compile(T.regexp)
		if err != nil {
			errorLog.Printf("regex error %s", err)
			return err
		}
		if ok := re.MatchString(val); !ok {
			return ErrInvaildByTag
		}
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
