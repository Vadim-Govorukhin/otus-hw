package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Используйте log.New() для создания логгера для записи информационных сообщений. Для этого нужно
// три параметра: место назначения для записи логов (os.Stdout), строка
// с префиксом сообщения (INFO или ERROR) и флаги, указывающие, какая
// дополнительная информация будет добавлена. Обратите внимание, что флаги
// соединяются с помощью оператора OR |.
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Создаем логгер для записи сообщений об ошибках таким же образом, но используем stderr как
// место для записи и используем флаг log.Lshortfile для включения в лог
// названия файла и номера строки где обнаружилась ошибка.
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func GetFileSize(fromPath string) (int64, error) {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		errorLog.Fatal(err)
		return 0, ErrUnsupportedFile
	}
	return fileInfo.Size(), nil
}

func CheckArgs(fromPath string, offset, limit int64) error {
	fileSize, err := GetFileSize(fromPath)
	if err != nil {
		errorLog.Fatal(err)
		return err
	}

	if fileSize > offset {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {

	err := CheckArgs(fromPath, offset, limit)
	if err != nil {
		errorLog.Fatal(err)
		return err
	}
	infoLog.Printf("Input args checked")

	outputFile, err := os.Create(filepath.Join(toPath, "out.txt"))
	if err != nil {
		errorLog.Fatal(err)
		return err
	}
	infoLog.Printf("Output file created")
	defer outputFile.Close()

	inputFile, err := os.Open(fromPath)
	if err != nil {
		errorLog.Fatal(err)
		return err
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile) // creates a new reader
	_, err = reader.Discard(int(offset)) // discard the following offset bytes
	if err != nil {
		errorLog.Fatal(err)
		return err
	}

	bufLimit := limit
	if limit > 512 { /////
		bufLimit = 64 /////
	}
	data := make([]byte, bufLimit)

	for i := 0; ; i++ {
		_, err = reader.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			errorLog.Fatal(err)
			return err
		}
		outputFile.Write(data)
	}
	infoLog.Printf("Wrote data to new file")

	return nil
}
