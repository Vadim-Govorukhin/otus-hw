package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	result := ""
	runes := []rune(s)
	r_tmp := "error"
	var flag bool
	r_s := ""

	for i, r := range runes {
		if flag {
			flag = false
			continue
		}

		if r == rune('\\') {
			r_s = string(r) + string(runes[i+1])
			flag = true
		} else {
			r_s = string(r)
		}

		switch {
		case unicode.IsDigit(r):
			num, _ := strconv.Atoi(r_s)
			if r_tmp == "error" {
				return "", ErrInvalidString
			}
			if num == 0 {
				result = strings.TrimSuffix(result, r_tmp)
			} else {
				for j := 1; j < num; j++ {
					result += r_tmp
				}
			}
			r_tmp = "error"

		case unicode.IsLetter(r) || unicode.IsSpace(r):
			r_tmp = r_s
			result += r_tmp
		default:
			return "", ErrInvalidString

		}
	}

	return result, nil
}
