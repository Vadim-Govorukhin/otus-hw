package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // for info messages
var ErrUnsupportedFileName = errors.New("unsupported file name")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

const invalidFileNames = "="

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environment := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadDir error: %w", err)
	}

	for _, f := range files {
		fileName := f.Name()
		infoLog.Println("Handling file", fileName)
		if strings.ContainsAny(fileName, invalidFileNames) {
			return nil, ErrUnsupportedFileName
		}

		var NeedRemove bool
		if f.Size() == 0 {
			NeedRemove = true
		}

		buf, err := ioutil.ReadFile(filepath.Join(dir, fileName))
		if err != nil {
			return nil, fmt.Errorf("ioutil.ReadFile error: %w", err)
		}
		// remove '\t' and ' ' from Value .
		buf = bytes.Replace(buf, []byte(`0x00`), []byte("\n"), 1) // replace terminal nulls
		buf = bytes.Split(buf, []byte("\n"))[0]                   // Keep the first line
		buf = bytes.TrimRight(buf, "\t ")

		environment[fileName] = EnvValue{Value: string(buf), NeedRemove: NeedRemove}
	}
	return environment, nil
}
