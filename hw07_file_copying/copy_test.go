package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const lumenText = `Ты можешь помолчать, ты можешь петь,
\t\tСтоять или бежать, но всё равно гореть.
Огромный синий кит порвать не может сеть.
Сдаваться или нет, но всё равно гореть.`

func TestSuppFunc(t *testing.T) {
	t.Run("check agrs", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		limit := int64(10)

		limit, err := CheckArgs(100, 10, limit)
		require.NoError(t, err, "Failed valid args")

		limit, err = CheckArgs(10, 100, limit)
		require.Equal(t, ErrOffsetExceedsFileSize, err, "Pass invalid args")
		require.Equal(t, int64(0), limit)

		limit, err = CheckArgs(100, 100, limit)
		require.NoError(t, err, "Failed offset = file size")
		require.Equal(t, int64(0), limit)

		limit, err = CheckArgs(100, 10, 0)
		require.NoError(t, err, "Failed zero limit")
		require.Equal(t, int64(100), limit)
	})

	t.Run("check progressBar", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		progressTemplate := "Completed %.2f%%"

		progressStr := ProgressBar(0, 100)
		expected := fmt.Sprintf(progressTemplate, 0.0)
		require.Equal(t, expected, progressStr)

		progressStr = ProgressBar(10, 100)
		expected = fmt.Sprintf(progressTemplate, 10.0)
		require.Equal(t, expected, progressStr)

		progressStr = ProgressBar(100, 100)
		expected = fmt.Sprintf(progressTemplate, 100.0)
		require.Equal(t, expected, progressStr)

		progressStr = ProgressBar(20, 80)
		expected = fmt.Sprintf(progressTemplate, 25.0)
		require.Equal(t, expected, progressStr)
	})
}

func TestCopyAsync(t *testing.T) {
	t.Run("check async copier lumen", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		var limit int64 = 70
		var fileSize int64 = 288
		var offset int64 = 70

		limit, _ = CheckArgs(fileSize, offset, limit)

		reader := strings.NewReader(lumenText)

		var buffWriter bytes.Buffer
		err := makeCopyAsync(reader, &buffWriter, limit, offset)
		require.NoError(t, err, "Failed to read from reader")

		s := buffWriter.String()
		require.Equal(t, "Стоять или бежать, но всё равно гореть.", s)
	})

	t.Run("check async copier test input", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		var limit int64 = 1000
		var offset int64 = 100
		isAsync := true

		curDir, _ := os.Getwd()
		fromPath := filepath.Join(curDir, "testdata", "input.txt")

		os.Mkdir("tmp", 0o750)
		defer os.RemoveAll("tmp")
		toPath := filepath.Join(curDir, "tmp", "out.txt")

		err := Copy(fromPath, toPath, offset, limit, isAsync)
		require.NoError(t, err, "Failed to check copy")

		b, _ := ioutil.ReadFile(toPath)
		outputText := string(b)

		b, _ = ioutil.ReadFile(filepath.Join(curDir, "testdata", "out_offset100_limit1000.txt"))
		goldenText := string(b)

		require.Equal(t, goldenText, outputText)
	})

	t.Run("check async copier test input", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		var limit int64
		var offset int64
		isAsync := true

		curDir, _ := os.Getwd()
		fromPath := filepath.Join(curDir, "testdata", "input.txt")

		os.Mkdir("tmp", 0o750)
		defer os.RemoveAll("tmp")
		toPath := filepath.Join(curDir, "tmp", "out.txt")

		err := Copy(fromPath, toPath, offset, limit, isAsync)
		require.NoError(t, err, "Failed to check copy")

		b, _ := ioutil.ReadFile(toPath)
		outputText := string(b)

		b, _ = ioutil.ReadFile(filepath.Join(curDir, "testdata", "out_offset0_limit0.txt"))
		goldenText := string(b)

		require.Equal(t, goldenText, outputText)
	})
}

func TestCopySync(t *testing.T) {
	t.Run("check copier lumen", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		var limit int64 = 70
		var fileSize int64 = 288
		var offset int64 = 70

		limit, _ = CheckArgs(fileSize, offset, limit)

		reader := strings.NewReader(lumenText)
		reader.Seek(offset, 0)

		var buffWriter bytes.Buffer
		err := makeCopySync(reader, &buffWriter, limit)
		require.NoError(t, err, "Failed to read from reader")

		s := buffWriter.String()
		require.Equal(t, "Стоять или бежать, но всё равно гореть.", s)
	})

	t.Run("check copier test input", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		var limit int64 = 1000
		var offset int64 = 100
		isAsync := false

		curDir, _ := os.Getwd()
		fromPath := filepath.Join(curDir, "testdata", "input.txt")

		os.Mkdir("tmp", 0o750)
		defer os.RemoveAll("tmp")
		toPath := filepath.Join(curDir, "tmp", "out.txt")

		err := Copy(fromPath, toPath, offset, limit, isAsync)
		require.NoError(t, err, "Failed to check copy")

		b, _ := ioutil.ReadFile(toPath)
		outputText := string(b)

		b, _ = ioutil.ReadFile(filepath.Join(curDir, "testdata", "out_offset100_limit1000.txt"))
		goldenText := string(b)

		require.Equal(t, goldenText, outputText)
	})

	t.Run("check async copier test input", func(t *testing.T) {
		infoLog.Printf("====== start test %s =====\n", t.Name())
		var limit int64
		var offset int64
		isAsync := false

		curDir, _ := os.Getwd()
		fromPath := filepath.Join(curDir, "testdata", "input.txt")

		os.Mkdir("tmp", 0o750)
		defer os.RemoveAll("tmp")
		toPath := filepath.Join(curDir, "tmp", "out.txt")

		err := Copy(fromPath, toPath, offset, limit, isAsync)
		require.NoError(t, err, "Failed to check copy")

		b, _ := ioutil.ReadFile(toPath)
		outputText := string(b)

		b, _ = ioutil.ReadFile(filepath.Join(curDir, "testdata", "out_offset0_limit0.txt"))
		goldenText := string(b)

		require.Equal(t, goldenText, outputText)
	})
}

/*
func BenchmarkPrimeNumbers(b *testing.B) {
    for i := 0; i < b.N; i++ {
        primeNumbers(num)
    }
}
*/
