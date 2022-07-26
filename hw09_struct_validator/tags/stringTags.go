package tags

import (
	"strconv"
	"strings"
)

type StringTags struct {
	len   int
	regex string
	in    []string
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
	case "regex":
		T.regex = m[1]
	case "in":
		T.in = strings.Split(m[1], ",")
	default:
		return ErrUnsupportedTag
	}
	return nil
}

func (T *StringTags) IsValid(i interface{}) (bool, error) {
	s, ok := i.(string)
	if !ok {
		return false, ErrUnsupportedTypeField
	}

	// len
	if (T.len != 0) || (len(s) != T.len) {
		return false, ErrInvaildByTag
	}

	// regex
	if T.regex != "" {
		//
		//return false, ErrInvaildByTag
	}

	// in
	if len(T.in) != 0 {
		var flag bool
		for _, str := range T.in {
			if str == s {
				flag = true
				break
			}
		}
		if !flag {
			return false, ErrInvaildByTag
		}
	}

	return true, nil
}
