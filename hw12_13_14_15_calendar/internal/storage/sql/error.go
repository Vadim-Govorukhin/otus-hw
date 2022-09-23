package sqlstorage

import "errors"

var (
	// ErrorLoadDriver = errors.New("failed to load driver")
	ErrorConnectDB    = errors.New("failed to connect to db")
	ErrorPrepareQuery = errors.New("failed to prepare queries")
)
