package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ErrUnsupportedFileName = errors.New("unsupported file name")

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info messages
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error messages
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	Environment := make(map[string]EnvValue)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		errorLog.Println(err)
		return nil, err
	}

	invalidFileNames := "="
	for _, f := range files {
		fileName := f.Name()
		infoLog.Println("Handling file", fileName)
		if strings.ContainsAny(fileName, invalidFileNames) {
			errorLog.Println(ErrUnsupportedFileName.Error())
			return nil, ErrUnsupportedFileName
		}

		var NeedRemove bool
		if f.Size() == 0 {
			NeedRemove = true
		}

		buf, err := ioutil.ReadFile(filepath.Join(dir, fileName))
		if err != nil {
			errorLog.Println(err)
			return nil, err
		}
		buf = bytes.Replace(buf, []byte(`0x00`), []byte("\n"), 1) // replace terminal nulls
		buf = bytes.Split(buf, []byte("\n"))[0]                   // Keep the first line

		Environment[fileName] = EnvValue{Value: string(buf), NeedRemove: NeedRemove}
	}
	return Environment, nil
}
