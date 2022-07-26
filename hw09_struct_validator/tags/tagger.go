package tags

import (
	"errors"
	"log"
	"os"
	"strings"
)

var (
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

var (
	ErrUnsupportedTypeField = errors.New("unsupported type field")
	ErrUnsupportedTag       = errors.New("unsupported tag")
	ErrInvaildByTag         = errors.New("field is invalid by tag")
)

type Tagger interface {
	FillField(string) error
	IsValid(interface{}) (bool, error)
}

func ParseTags(tag string, typeField string) (Tagger, error) {
	InfoLog.Printf("parse tags %s of field type %s", tag, typeField)

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
			ErrorLog.Printf("parsing error %e", err)
			return nil, ErrUnsupportedTag
		}
	}
	return tagStruct, nil
}
