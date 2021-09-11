package main

import (
	"os"

	"github.com/alecthomas/kingpin"
)

type argsParser interface {
	parse([]string) (config, error)
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

	return argsParser(kparser)
}

func (k *kingpinParser) parse(args []string) (config, error) {
	k.App.Name = args[0]
	_, err := k.App.Parse(args[1:])
	if err != nil {
		return emptyConf, err
	}

	healthCheckUrl := os.Getenv("HEALTHCHECK_REPORT_URL")
	if healthCheckUrl == "" {
		healthCheckUrl = DEFAULT_HEATHCHECK_URL
	}
	lineLoginAPIUrl := os.Getenv("LINE_LOGIN_API_URL")
	if lineLoginAPIUrl == "" {
		lineLoginAPIUrl = DEFAULT_LINE_LOGIN_API_URL
	}

	return config{
		File:                     k.File,
		LINE_LOGIN_CODE:          os.Getenv("LINE_LOGIN_CODE"),
		LINE_LOGIN_REDIRECT_URI:  os.Getenv("LINE_LOGIN_REDIRECT_URI"),
		LINE_LOGIN_CLIENT_ID:     os.Getenv("LINE_LOGIN_CLIENT_ID"),
		LINE_LOGIN_CLIENT_SECRET: os.Getenv("LINE_LOGIN_CLIENT_SECRET"),
		LINE_LOGIN_API_URL:       lineLoginAPIUrl,
		HEALTHCHECK_URL:          healthCheckUrl,
	}, nil
}
