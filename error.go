package main

import (
	"errors"
	"fmt"
)

func errorf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}
