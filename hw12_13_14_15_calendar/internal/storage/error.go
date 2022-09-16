package storage

import "errors"

var (
	// Db Errors
	ErrorLoadDriver = errors.New("failed to load driver")
	ErrorConnectDB  = errors.New("failed to connect to db")

	// Event Errors
	ErrorDayBusy       = errors.New("данное время уже занято другим событием")
	ErrorEventIDBusy   = errors.New("cобытие с таким id уже существует")
	ErrorWrongUpdateID = errors.New("обновление события с изменением ID")
	// TODO
)
