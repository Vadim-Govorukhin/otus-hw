package tags

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strings"
)

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

var (
	ErrUnsupportedTypeField = errors.New("unsupported type field")
	ErrUnsupportedTag       = errors.New("unsupported tag")
	ErrInvaildByTag         = errors.New("field is invalid by tag")
)

type Tagger interface {
	FillField(string) error
	ValidateValue(reflect.Value) error
}

func ParseTags(tag string, typeField string) (Tagger, error) {
	tagStruct := chooseTagStruct(typeField)
	if tagStruct == nil {
		return nil, ErrUnsupportedTypeField
	}

	slTag := strings.Split(tag, "|")
	var err error
	for _, t := range slTag {
		err = tagStruct.FillField(t)
		if err != nil {
			ErrorLog.Printf("parsing error %e", err)
			return nil, ErrUnsupportedTag
		}
	}
	return tagStruct, nil
}

func chooseTagStruct(typeField string) (tagStruct Tagger) {
	switch typeField {
	case "string", "[]string":
		tagStruct = &StringTags{}
	case "int", "[]int":
		tagStruct = &IntTags{}
	default:
		switch reflect.TypeOf(typeField).String() {
		case "string", "[]string":
			tagStruct = &StringTags{}
		case "int", "[]int":
			tagStruct = &IntTags{}
		}
	}
	return
}
