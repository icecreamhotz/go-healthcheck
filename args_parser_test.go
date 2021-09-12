package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	programName = "go-healthcheck"
)

func TestInvalidArgsParsingWithFile(t *testing.T) {
	filename := "assets/example.csv"
	expectations := []struct {
		in  []string
		out error
	}{
		{
			[]string{programName, filename},
			nil,
		},
	}

	for _, v := range expectations {
		p := newKingpinParser()
		file, err := p.parse(v.in)
		assert.Equal(t, file, filename)
		assert.Equal(t, err, v.out)
	}
}

func TestInvalidArgsParsingWithOutFile(t *testing.T) {
	expectations := []struct {
		in  []string
		out string
	}{
		{
			[]string{programName},
			"required argument 'file' not provided",
		},
	}

	for _, v := range expectations {
		p := newKingpinParser()
		if _, err := p.parse(v.in); err == nil || err.Error() != v.out {
			t.Error(err, v.out)
		}
	}
}
