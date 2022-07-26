package hw09structvalidator

import (
	"errors"
	"log"
	"os"
	"reflect"
)

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

var (
	ErrUnsupportedTag       = errors.New("unsupported tag")
	ErrUnsupportedTypeField = errors.New("unsupported type field")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	// IMPLEMENT ME
	return v[0].Err.Error()
}

func Validate(v interface{}) error {
	infoLog.Printf("start validate type %[1]T with val %[1]v\n", v)
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	for i := 0; i < val.NumField(); i++ {
		f := t.Field(i)
		fv := val.Field(i)

		tag := f.Tag.Get("validate")
		if tag != "" {
			tagStruct, err := parseTags(tag, fv.Type().String())
			if err != nil {
				errorLog.Printf("parsing error %e", err)
				return err
			}
			err = validateField(fv, tagStruct)
			if err != nil {
				errorLog.Printf("validating error %e", err)
				return err
			}
		}

	}
	return nil
}

func validateField(v reflect.Value, tagStruct Tagger) error {
	infoLog.Printf("validate value %v with tags %+v\n", v, tagStruct)

	switch v.Kind() {
	case reflect.Invalid:
		errorLog.Printf("%s is invalid\n", v)
	case reflect.Int:
		infoLog.Printf("validate int %v\n", v.Int())
	case reflect.String:
		infoLog.Printf("validate string %s\n", v.String())
	case reflect.Slice:
		infoLog.Println("validate slice")
		for i := 0; i < v.Len(); i++ {
			Validate(v.Index(i))
		}
	}
	return nil
}
