package repository

import "errors"

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("record not found")
)

