package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	t.Run("it should work without argument", func(t *testing.T) {
		expectation := errors.New("This is a expectation value")
		value := "This is a expectation value"
		assert.Equal(t, expectation, errorf(value))
	})

	t.Run("it should work with one argument", func(t *testing.T) {
		expectation := errors.New("This is a expectation value is 1")
		value := "This is a expectation value is %v"
		assert.Equal(t, expectation, errorf(value, 1))
	})

	t.Run("it should work with many arguments", func(t *testing.T) {
		expectation := errors.New("This is a expectation value is 1 2 3 4 5")
		value := "This is a expectation value is %v %v %v %v %v"
		assert.Equal(t, expectation, errorf(value, 1, 2, 3, 4, 5))
	})
}
