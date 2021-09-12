package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorFWithOutArgument(t *testing.T) {
	expectation := errors.New("This is a expectation value")
	value := "This is a expectation value"
	assert.Equal(t, expectation, errorf(value))
}

func TestErrorFWithOneArgument(t *testing.T) {
	expectation := errors.New("This is a expectation value is 1")
	value := "This is a expectation value is %v"
	assert.Equal(t, expectation, errorf(value, 1))
}

func TestErrorFWithAnyArgument(t *testing.T) {
	expectation := errors.New("This is a expectation value is 1 2 3 4 5")
	value := "This is a expectation value is %v %v %v %v %v"
	assert.Equal(t, expectation, errorf(value, 1, 2, 3, 4, 5))
}
