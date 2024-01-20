package utils

import (
	"strings"
)

type Error struct {
	Msg string
}

func (err Error) Error() string {
	return err.Msg
}

func (err Error) Is(target error) bool {
	return strings.Contains(strings.ToLower(target.Error()), strings.ToLower(err.Msg))
}
