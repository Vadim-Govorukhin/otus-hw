package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedFileName = errors.New("unsupported file name")
)

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                 // for info message
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // for error message
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func (env EnvValue) String() string {
	return fmt.Sprintf("\tValue: %s,\n\t NeedRemove: %t", env.Value, env.NeedRemove)
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

	var invalidFileNames = "="
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

		Environment[fileName] = EnvValue{Value: string(buf), NeedRemove: NeedRemove}

	}

	return Environment, nil
}
