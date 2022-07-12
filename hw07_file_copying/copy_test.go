package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	/*
		t.Run("check filesize", func(t *testing.T) {

		})
	*/
	t.Run("check agrs", func(t *testing.T) {

		err := CheckArgs(100, 10, 50)
		require.NoError(t, err, "Failed to check valid args")

		err = CheckArgs(10, 100, 50)
		require.Equal(t, ErrOffsetExceedsFileSize, err)

	})

	t.Run("check progressBar", func(t *testing.T) {

		progressStr := ProgressBar(0, 10, 100, 100)
		require.Equal(t, "Выполнено 0%", progressStr)

		progressStr = ProgressBar(1, 10, 100, 100)
		require.Equal(t, "Выполнено 10%", progressStr)

		progressStr = ProgressBar(10, 10, 100, 100)
		require.Equal(t, "Выполнено 100%", progressStr)

		progressStr = ProgressBar(1, 20, 80, 100)
		require.Equal(t, "Выполнено 25%", progressStr)

	})
	/*
		fs := fstest.MapFS{
			"hello.txt": {
				Data: []byte("hello, world"),
			},
		}
		data, err := fs.ReadFile("hello.txt")
		if err != nil {
			panic(err)
		}
		println(string(data) == "hello, world")

		//*
			var buffer bytes.Buffer
			buffer.WriteString("fake, csv, data")
			content, err := readFile(&buffer)
			if err != nil {
				t.Error("Failed to read csv data")
			}
			fmt.Print(content)
	*/
}
