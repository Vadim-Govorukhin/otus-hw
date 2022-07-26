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
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
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
	InfoLog.Printf("start validate type %[1]T with val %[1]v\n", v)
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	var validationErrors ValidationErrors
	for i := 0; i < val.NumField(); i++ {
		f := t.Field(i)
		fv := val.Field(i)

		tag := f.Tag.Get("validate")
		if tag != "" {
			tagStruct, err := tags.ParseTags(tag, fv.Type().String())
			if err != nil {
				ErrorLog.Printf("parsing error %e", err)
				return err
			}
			validationErr := validateField(f.Name, fv, tagStruct)
			validationErrors = append(validationErrors, validationErr)
		}

	}
	InfoLog.Println("List of errors", validationErrors)
	return nil
}

func validateField(name string, v reflect.Value, tagStruct tags.Tagger) ValidationError {
	InfoLog.Printf("validate value %v with tags %+v\n", v, tagStruct)
	isValid, err := tagStruct.IsValid(v)
	InfoLog.Printf("value is %v with errors %+e\n", isValid, err)
	return ValidationError{name, err}

	/*
		switch v.Kind() {
		case reflect.Invalid:
			ErrorLog.Printf("%s is invalid\n", v)
		case reflect.Int:
			InfoLog.Printf("validate int %v\n", v.Int())
		case reflect.String:
			InfoLog.Printf("validate string %s\n", v.String())
		case reflect.Slice:
			InfoLog.Println("validate slice")
			for i := 0; i < v.Len(); i++ {
				Validate(v.Index(i))
			}
		}
	*/
}
