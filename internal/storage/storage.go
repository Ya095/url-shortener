package storage

import "errors"


// Ошибки при работе с БД.
var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists = errors.New("url exists")
)