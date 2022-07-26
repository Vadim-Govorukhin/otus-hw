package hw09structvalidator

import (
	"strconv"
	"strings"
)

type Tagger interface {
	FillField(string) error
}

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
			errorLog.Printf("parsing error %e", err)
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

type IntTags struct {
	min int
	max int
	in  []int
}

func (T *IntTags) FillField(tag string) error {
	m := strings.Split(tag, ":")
	if len(m) > 2 {
		return ErrUnsupportedTag
	}

	switch m[0] {
	case "min":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			errorLog.Printf("parsing error %e", err)
			return err
		}
		T.min = i
	case "max":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			errorLog.Printf("parsing error %e", err)
			return err
		}
		T.max = i
	case "in":
		var arr []int
		var err error
		for i, s := range strings.Split(m[1], ",") {
			arr[i], err = strconv.Atoi(s)
			if err != nil {
				errorLog.Printf("parsing error %e", err)
				return err
			}
		}
		T.in = arr
	default:
		return ErrUnsupportedTag
	}
	return nil
}

func parseTags(tag string, typeField string) (Tagger, error) {
	infoLog.Printf("parse tags %s of field type %s", tag, typeField)

	var tagStruct Tagger
	switch typeField {
	case "string":
		tagStruct = &StringTags{}
	case "int":
		tagStruct = &IntTags{}
	default:
		return nil, ErrUnsupportedTypeField
	}

	slTag := strings.Split(tag, "|")
	var err error
	for _, t := range slTag {
		err = tagStruct.FillField(t)
		if err != nil {
			errorLog.Printf("parsing error %e", err)
			return nil, ErrUnsupportedTag
		}
	}
	return tagStruct, nil
}
