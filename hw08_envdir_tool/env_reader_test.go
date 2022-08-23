package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testFiles = []struct {
	fileName string
	body     interface{}
	expected EnvValue
}{
	{fileName: "NUM", body: 123, expected: EnvValue{"123", false}},
	{fileName: "STR0", body: "", expected: EnvValue{"", true}},
	{fileName: "STRT", body: "string\t", expected: EnvValue{"string", false}},
	{fileName: "TSTR", body: "\tstring", expected: EnvValue{"\tstring", false}},
	{fileName: "STRS", body: "string ", expected: EnvValue{"string", false}},
	{fileName: "STR2", body: "string\nSECON STRING", expected: EnvValue{"string", false}},
	{fileName: "STR3", body: "\nSECON STRING", expected: EnvValue{"", false}},
	{fileName: "HELLO", body: "\"hello\"", expected: EnvValue{"\"hello\"", false}},
	{fileName: "F.TXT", body: "i am txt", expected: EnvValue{"i am txt", false}},
	{fileName: "NIL", body: "to delete", expected: EnvValue{"", true}},
	{fileName: "NIL", body: nil, expected: EnvValue{"", true}},
	{fileName: "NotNIL", body: nil, expected: EnvValue{"", true}},
	{fileName: "NotNIL", body: "not delete", expected: EnvValue{"not delete", false}},
	{
		fileName: "NIL2", body: "   foo" + string(rune(0x00)) + "with new line",
		expected: EnvValue{"   foo\nwith new line", false},
	},
}

func prepareTestDir(t *testing.T, testDirName string) (string, error) {
	t.Helper()
	dir, err := os.MkdirTemp("", testDirName)
	require.NoError(t, err)

	for _, testCase := range testFiles {
		newFile, err := os.Create(filepath.Join(dir, testCase.fileName))
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
	return dir, nil
}

func setupTest(t *testing.T) {
	t.Helper()
	infoLog.Printf("====== start test %s =====\n", t.Name())
}

func TestReadDir(t *testing.T) {
	resultExpected := make(Environment)
	for _, testCase := range testFiles {
		resultExpected[testCase.fileName] = testCase.expected
	}

	t.Run("Empty dir", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "prefix")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		env, err := ReadDir(dir)
		require.Equal(t, env, Environment{})
		require.NoError(t, err)
	})

	t.Run("Wrong dir path", func(t *testing.T) {
		_, err := ReadDir("blabla")
		require.Error(t, err)
	})

	t.Run("preparing test dir", func(t *testing.T) {
		setupTest(t)

		const testDirName = "testDir"
		defer os.RemoveAll(testDirName)

		_, err := prepareTestDir(t, testDirName)
		require.NoError(t, err)
	})

	t.Run("invalid filename", func(t *testing.T) {
		setupTest(t)

		testDirName := "testDir"
		dir, err := os.MkdirTemp("", testDirName)
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		newFile, _ := os.Create(filepath.Join(dir, "VAR="))
		newFile.Close()

		_, err = ReadDir(dir)
		require.ErrorIs(t, err, ErrUnsupportedFileName)
	})

	t.Run("test cases", func(t *testing.T) {
		setupTest(t)

		testDirName := "testDir"
		testDirPath, _ := prepareTestDir(t, testDirName)
		defer os.RemoveAll(testDirName)

		environment, err := ReadDir(testDirPath)
		require.NoError(t, err)
		require.Equal(t, len(resultExpected), len(environment))
		for key, val := range environment {
			infoLog.Println("check file ", key)
			require.Equal(t, resultExpected[key], val)
		}
	})
}
