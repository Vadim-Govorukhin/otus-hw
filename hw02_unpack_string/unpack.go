package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (result string, err error) {
	var resultBuilder strings.Builder

	var repeatnum int
	var specialFlag bool
	var curRuneIsDig bool
	var nextRuneIsDig bool

	posRunes := []rune(s)
	for i, r := range posRunes {
		curRuneIsDig = unicode.IsDigit(r)
		nextRuneIsDig = getNextRuneIsDig(i, posRunes)

		if (i == 0 && curRuneIsDig) ||
			(curRuneIsDig && nextRuneIsDig && !specialFlag) {
			return "", ErrInvalidString
		}

		if curRuneIsDig || specialFlag {
			specialFlag = false
			continue
		}

		repeatnum = 1
		if r == rune('\\') {
			specialFlag = true
			if i != len(posRunes)-1 {
				r, err = handleSpecial(posRunes[i+1])
				if err != nil {
					return "", err
				}
			} else {
				return "", ErrInvalidString
			}
			i++
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
			resultBuilder.WriteRune(r)
		}
	}
	result = resultBuilder.String()
	return result, nil
}

func handleSpecial(r rune) (rune, error) {
	// in this cases special is always "\"
	switch r {
	case 'n':
		return '\n', nil
	case 't':
		return '\t', nil
	}
	if !(unicode.IsDigit(r) || r == rune('\\')) {
		return 0, ErrInvalidString
	}
	return r, nil
}

func getNextRuneIsDig(i int, posRunes []rune) bool {
	if i != len(posRunes)-1 {
		return unicode.IsDigit(posRunes[i+1])
	}
	return false
}
