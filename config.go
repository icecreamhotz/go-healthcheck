package main

import (
	"os"
)

type configurater interface {
	checkArgs() error
	checkFileExist() error
	checkLineEnv() error
}

type config struct {
	File                     string
	LINE_LOGIN_CODE          string
	LINE_LOGIN_REDIRECT_URI  string
	LINE_LOGIN_CLIENT_ID     string
	LINE_LOGIN_CLIENT_SECRET string
	LINE_LOGIN_API_URL       string
	HEALTHCHECK_URL          string
}

func newConfig(file string) config {
	healthCheckUrl := os.Getenv("HEALTHCHECK_REPORT_URL")
	if healthCheckUrl == "" {
		healthCheckUrl = DEFAULT_HEATHCHECK_URL
	}
	lineLoginAPIUrl := os.Getenv("LINE_LOGIN_API_URL")
	if lineLoginAPIUrl == "" {
		lineLoginAPIUrl = DEFAULT_LINE_LOGIN_API_URL
	}

	return config{
		File: file,
		LINE_LOGIN_CODE: os.Getenv("LINE_LOGIN_CODE"),
		LINE_LOGIN_REDIRECT_URI: os.Getenv("LINE_LOGIN_REDIRECT_URI"),
		LINE_LOGIN_CLIENT_ID: os.Getenv("LINE_LOGIN_CLIENT_ID"),
		LINE_LOGIN_CLIENT_SECRET: os.Getenv("LINE_LOGIN_CLIENT_SECRET"),
		LINE_LOGIN_API_URL: lineLoginAPIUrl,
		HEALTHCHECK_URL: healthCheckUrl,
	}
}

func (c *config) checkArgs() error {
	checks := []func() error{
		c.checkFileExist,
		c.checkLineEnv,
	}

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	return nil
}

func (c *config) checkFileExist() error {
	if _, err := os.Stat(c.File); os.IsNotExist(err) {
		return errInvalidURL
	}

	return nil
}

func (c *config) checkLineEnv() error {
	if c.LINE_LOGIN_CODE == "" || c.LINE_LOGIN_REDIRECT_URI == "" || c.LINE_LOGIN_CLIENT_ID == "" || c.LINE_LOGIN_CLIENT_SECRET == "" {
		return errLineEnvBlank
	}

	return nil
}
