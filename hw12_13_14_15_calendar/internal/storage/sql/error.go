package sqlstorage

import "errors"

var (
	ErrorConnectDB    = errors.New("failed to connect to db")
	ErrorPrepareQuery = errors.New("failed to prepare queries")
)
