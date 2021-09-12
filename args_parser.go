package main

import (
	"github.com/alecthomas/kingpin"
)

type argsParser interface {
	parse([]string) (string, error)
}

type kingpinParser struct {
	App  *kingpin.Application
	File string
}

func newKingpinParser() argsParser {
	kparser := &kingpinParser{
		File: "",
	}

	app := kingpin.New("", "go-healthcheck challenge").
		Version("1.0.0")
	app.Arg("file", "Csv file is required.").Required().
		StringVar(&kparser.File)

	kparser.App = app

	return kparser
}

func (k *kingpinParser) parse(args []string) (string, error) {
	k.App.Name = args[0]
	_, err := k.App.Parse(args[1:])
	if err != nil {
		return "", err
	}

	return k.File, nil
}
