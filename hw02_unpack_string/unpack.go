package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (result string, err error) {
	posRunes := []rune(s)
	var specialFlag bool
	var repeatnum int
	var curRuneIsDig bool
	var nextRuneIsDig bool

	for i, r := range posRunes {
		curRuneIsDig = unicode.IsDigit(r)
		nextRuneIsDig = getNextRuneIsDig(i, posRunes)

		if i == 0 && curRuneIsDig ||
			(curRuneIsDig && nextRuneIsDig && !specialFlag) {
			return "", ErrInvalidString
		}
		repeatnum = 1

		if curRuneIsDig || specialFlag {
			specialFlag = false
			continue
		}

		if r == rune('\\') {
			specialFlag = true
			if i != len(posRunes)-1 {
				r = handlespecial(r, posRunes[i+1])
			} else {
				return "", ErrInvalidString
			}
			i += 1
			nextRuneIsDig = getNextRuneIsDig(i, posRunes)
		}

		if nextRuneIsDig {
			nextr := posRunes[i+1]
			repeatnum, err = strconv.Atoi(string(nextr))
			if err != nil {
				return "", err
			}
		}

		for j := 0; j < repeatnum; j++ {
			result += string(r)
		}
	}
	return result, nil
}

func handlespecial(special rune, r rune) rune {
	// in this cases special is "\"
	return r
}

func getNextRuneIsDig(i int, posRunes []rune) (nextRuneIsDig bool) {
	if i != len(posRunes)-1 {
		nextRuneIsDig = unicode.IsDigit(posRunes[i+1])
	} else {
		nextRuneIsDig = false
	}
	return
}
