package main

import (
	"errors"
)

const (
	DEFAULT_HEATHCHECK_URL     = "https://backend-challenge.line-apps.com/healthcheck/report"
	DEFAULT_LINE_LOGIN_API_URL = "https://api.line.me/oauth2/v2.1/token"
)

var (
	// common variables
	emptyConf = config{}
	parser    = newKingpinParser()

	lineLoginGrantType    = "authorization_code"
	contentTypeJSON       = "application/json"
	contentTypeUrlEncoded = "application/x-www-form-urlencoded"

	// errors
	errInvalidURL = errors.New(
		"Csv file doesn't not exists.")
	errLineEnvBlank = errors.New(
		"Line env doest't not exists.")
)
