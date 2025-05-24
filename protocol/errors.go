package protocol

import "errors"

var (
	ErrQuit = errors.New("client requested quit")
)
