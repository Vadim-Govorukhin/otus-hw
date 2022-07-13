package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("check filesize", func(t *testing.T) {
		curDir, _ := os.Getwd()
		fromPath := filepath.Join(curDir, "testdata", "input.txt")

		fileSize, err := GetFileSize(fromPath)
		require.NoError(t, err, "Failed to check file size")
		require.Equal(t, int64(6617), fileSize)

		fromPath = ""
		fileSize, err = GetFileSize(fromPath)
		require.Equal(t, ErrUnsupportedFile, err)
		require.Equal(t, int64(0), fileSize)

	})

	t.Run("check agrs", func(t *testing.T) {

		err := CheckArgs(100, 10, 50)
		require.NoError(t, err, "Failed to check valid args")

		err = CheckArgs(10, 100, 50)
		require.Equal(t, ErrOffsetExceedsFileSize, err)

	})

	t.Run("check progressBar", func(t *testing.T) {

		progressStr := ProgressBar(0, 100, 100)
		require.Equal(t, "Выполнено 0%", progressStr)

		progressStr = ProgressBar(10, 100, 100)
		require.Equal(t, "Выполнено 10%", progressStr)

		progressStr = ProgressBar(100, 100, 100)
		require.Equal(t, "Выполнено 100%", progressStr)

		progressStr = ProgressBar(20, 80, 100)
		require.Equal(t, "Выполнено 25%", progressStr)

	})

	t.Run("check copier lumen", func(t *testing.T) {
		var limit int64 = 70
		var fileSize int64 = 288
		var offset int64 = 73

		var buffReader bytes.Buffer
		buffReader.WriteString(
			`Ты можешь помолчать, ты можешь петь,
			\t\tСтоять или бежать, но всё равно гореть.
		Огромный синий кит порвать не может сеть.
		Сдаваться или нет, но всё равно гореть.`)

		reader := bufio.NewReader(&buffReader) // creates a new reader
		reader.Discard(int(offset))            // discard the following offset bytes

		var buffWriter bytes.Buffer

		progressBar := func(i int64) string { return ProgressBar(i, limit, fileSize-offset) }
		err := makeCopy(reader, &buffWriter, limit, progressBar)
		require.NoError(t, err, "Failed to read from reader")
		s := buffWriter.String()
		require.Equal(t, "Стоять или бежать, но всё равно гореть.", s)

	})

	t.Run("check copier test input", func(t *testing.T) {
		var limit int64 = 1000
		var offset int64 = 100

		curDir, _ := os.Getwd()
		fromPath := filepath.Join(curDir, "testdata", "input.txt")

		os.Mkdir("tmp", 0750)
		defer os.Remove("tmp")
		toPath := filepath.Join(curDir, "tmp")

		err := Copy(fromPath, toPath, offset, limit)
		require.NoError(t, err, "Failed to check file size")

		b, _ := ioutil.ReadFile(filepath.Join(toPath, "out.txt"))
		outputText := string(b)

		b, _ = ioutil.ReadFile(filepath.Join(curDir, "testdata", "out_offset100_limit1000.txt"))
		goldenText := string(b)

		require.Equal(t, goldenText, outputText)
	})

}
