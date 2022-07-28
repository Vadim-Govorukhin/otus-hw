package hw09structvalidator

import (
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

func (v ValidationErrors) Is(tgt error) bool {
	return v.Error() == tgt.Error()
}

func Validate(v interface{}) error {
	infoLog.Printf("start validate struct %+v\n", v)
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	var validationErrors ValidationErrors
	for i := 0; i < val.NumField(); i++ {
		f := t.Field(i)

		infoLog.Printf("search tag 'validate' of field '%v'", f.Name)
		if tag, ok := f.Tag.Lookup("validate"); ok {
			fv := val.Field(i)
			var err error
			if f.Type.Kind() == reflect.Struct {
				err = Validate(fv.Interface())
			} else {
				infoLog.Printf("\tcheck field '%v': value '%v' and tags '%s'", f.Name, fv, tag)
				err = validateField(fv, tag)
			}
			validationErrors = append(validationErrors, ValidationError{f.Name, err})
		}
	}
	infoLog.Println("List of errors", validationErrors)
	return validationErrors
}

func validateField(fv reflect.Value, tag string) error {
	var err error
	tagStruct, err := tags.ParseTags(tag, fv.Type().String())
	if err != nil {
		errorLog.Printf("parsing error %e", err)
		return err
	}

	if fv.Kind() == reflect.Slice {
		for i := 0; i < fv.Len(); i++ {
			err = tagStruct.ValidateValue(fv.Index(i))
			if err != nil {
				break
			}
		}
	} else {
		err = tagStruct.ValidateValue(fv)
	}
	infoLog.Printf("\tend check with error %v\n", err)
	return err
}
