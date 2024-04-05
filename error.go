package main

import "errors"

var (
	ErrNoContent = errors.New("content cannot be empty")
	ErrNoTitle   = errors.New("title has to be specified")
	ErrNoEvent   = errors.New("event not found")
)
