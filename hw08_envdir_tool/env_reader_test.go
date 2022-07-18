package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var testFiles = []struct {
	fileName string
	body     interface{}
}{
	{fileName: "NUM", body: 123},
	{fileName: "STR", body: "string"},
	{fileName: "NIL", body: "to delete"},
	{fileName: "NIL", body: nil},
	{fileName: "STR2", body: "string\nSECON STRING"},
	{fileName: "HELLO", body: "\"hello\""},
	{fileName: "F.TXT", body: "i am txt"},
}

func shouldGetwd(t *testing.T) string {
	t.Helper()
	curDir, err := os.Getwd()
	require.NoError(t, err)
	return curDir
}

func prepareTestDir(t *testing.T, testDirName string) (string, error) {
	t.Helper()
	os.Mkdir(testDirName, 0o750)
	testDirPath := filepath.Join(shouldGetwd(t), testDirName)

	for _, testCase := range testFiles {
		newFile, err := os.Create(filepath.Join(testDirPath, testCase.fileName))
		if err != nil {
			return "", err
		}

		var data []byte
		switch v := testCase.body.(type) {
		case int:
			data = []byte(fmt.Sprintf("%d", v))
		case string:
			data = []byte(v)
		}
		if len(data) != 0 {
			newFile.Write(data)
		}
		newFile.Close()
	}
	return testDirPath, nil
}

func setupTest(t *testing.T) {
	infoLog.Printf("====== start test %s =====\n", t.Name())
}
func TestReadDir(t *testing.T) {
	resultTestFiles := make(Environment)
	for _, t := range testFiles {
		var NeedRemove bool
		if (t.body == "to delete") || (t.body == nil) { // TODO
			NeedRemove = true
		}

		var Value string
		switch v := t.body.(type) { // TODO
		case int:
			Value = fmt.Sprintf("%d", v)
		case string:
			Value = strings.Split(v, "\n")[0]
		}

		resultTestFiles[t.fileName] = EnvValue{Value: Value, NeedRemove: NeedRemove}
	}

	t.Run("preparing test dir", func(t *testing.T) {
		testDirName := "testDir"
		defer os.RemoveAll(testDirName)
		_, err := prepareTestDir(t, testDirName)
		require.NoError(t, err)
	})

	t.Run("invalid filename", func(t *testing.T) {
		setupTest(t)

		testDirName := "testDir"
		defer os.RemoveAll(testDirName)
		os.Mkdir(testDirName, 0o750)
		testDirPath := filepath.Join(shouldGetwd(t), testDirName)
		newFile, _ := os.Create(filepath.Join(testDirPath, "VAR="))
		newFile.Close()

		_, err := ReadDir(testDirPath)
		require.ErrorIs(t, err, ErrUnsupportedFileName)
	})

	t.Run("test case", func(t *testing.T) {
		setupTest(t)

		testDirName := "testDir"
		defer os.RemoveAll(testDirName)
		testDirPath, _ := prepareTestDir(t, testDirName)

		Environment, err := ReadDir(testDirPath)
		require.NoError(t, err)
		require.Equal(t, resultTestFiles, Environment)

	})
}
