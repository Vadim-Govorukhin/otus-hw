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

func (t *StringTags) FillField(tag string) error {
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
		t.len = i
	case "regexp":
		t.regexp = m[1]
	case "in":
		t.in = strings.Split(m[1], ",")
	default:
		errorLog.Printf("unsupported tag name: %s\n", m[0])
		return ErrUnsupportedTag
	}
	return nil
}

func (t *StringTags) ValidateValue(i reflect.Value) error {
	infoLog.Printf("\tvalidate value '%v' with tags %+v\n", i, *t)
	val := i.String()

	// len
	if (t.len != 0) && (len(val) != t.len) {
		return ErrInvaildByTag
	}

	// regex
	if t.regexp != "" {
		re, err := regexp.Compile(t.regexp)
		if err != nil {
			errorLog.Printf("regex error %s", err)
			return err
		}
		if ok := re.MatchString(val); !ok {
			return ErrInvaildByTag
		}
	}

	// in
	if len(t.in) != 0 {
		var flag bool
		for _, str := range t.in {
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
