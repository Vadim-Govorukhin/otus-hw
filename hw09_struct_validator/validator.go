package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/Vadim-Govorukhin/otus-hw/hw09_struct_validator/tags"
)

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

var (
	ErrUnsupportedTypeField = errors.New("unsupported type field")
	ErrUnsupportedTag       = errors.New("unsupported tag")
	ErrInvaildByTag         = errors.New("field is invalid by tag")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := ""
	for _, val := range v {
		res += fmt.Sprintf("%+v", val)
	}
	return res
}

func Validate(v interface{}) error {
	infoLog.Printf("start validate struct %+v\n", v)
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	var validationErrors ValidationErrors
	for i := 0; i < val.NumField(); i++ {
		f := t.Field(i)
		fv := val.Field(i)

		if tag, ok := f.Tag.Lookup("validate"); ok {
			infoLog.Printf("check field '%v': value '%v' and tags '%s'", f.Name, fv, tag)
			tagStruct, err := tags.ParseTags(tag, fv.Type().String())
			if err != nil {
				errorLog.Printf("parsing error %e", err)
				return err
			}

			infoLog.Printf("\tvalidate value '%v' with tags %+v\n", fv, tagStruct)
			err = tagStruct.IsValid(fv)
			infoLog.Printf("\tend check with error %v\n", err)

			validationErrors = append(validationErrors, ValidationError{f.Name, err})
		}

	}
	infoLog.Println("List of errors", validationErrors)
	return validationErrors
}
