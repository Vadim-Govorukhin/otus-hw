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
	infoLog                = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // for info messages
	ErrUnsupportedFileName = errors.New("unsupported file name")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

const invalidFilenameChars = "="

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadDir error: %w", err)
	}

	for _, f := range files {
		fileName := f.Name()
		// infoLog.Println("Handling file: ", fileName)
		if strings.ContainsAny(fileName, invalidFilenameChars) {
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

		val := getEnvValue(buf)
		env[fileName] = EnvValue{Value: val, NeedRemove: NeedRemove}
	}
	return env, nil
}

func getEnvValue(buf []byte) string {
	s := string(buf)
	s = strings.Split(s, "\n")[0]                       // Keep the first line
	s = strings.ReplaceAll(s, string(rune(0x00)), "\n") // replace terminal nulls
	s = strings.TrimRight(s, "\t ")
	return s
}
