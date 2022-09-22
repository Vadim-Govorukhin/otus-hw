package storage

import "errors"

var (
	ErrorWrongTypeStorage = errors.New("wrong type of storage")

	// Db Errors.
	ErrorLoadDriver            = errors.New("failed to load driver")
	ErrorConnectDB             = errors.New("failed to connect to db")
	ErrorPreparedQueryNotFound = errors.New("prepared query not found")

	// Event Errors.
	ErrorDayBusy     = errors.New("данное время уже занято другим событием")
	ErrorEventIDBusy = errors.New("cобытие с таким id уже существует")
	ErrorWrongID     = errors.New("cобытие с таким id не существует")
)
