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
	lineLoginCode := os.Getenv("LINE_LOGIN_CODE")
	lineLoginRedirectUri := os.Getenv("LINE_LOGIN_REDIRECT_URI")
	lineLoginClientID := os.Getenv("LINE_LOGIN_CLIENT_ID")
	lineLoginClientSecret := os.Getenv("LINE_LOGIN_CLIENT_SECRET")

	if lineLoginCode == "" || lineLoginRedirectUri == "" || lineLoginClientID == "" || lineLoginClientSecret == "" {
		return errLineEnvBlank
	}

	return nil
}
