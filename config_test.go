package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigShouldReturnOsEnv(t *testing.T) {
	t.Run("it should return os env", func(t *testing.T) {
		t.Setenv("LINE_LOGIN_CODE", "LINE_LOGIN_CODE")
		t.Setenv("LINE_LOGIN_REDIRECT_URI", "LINE_LOGIN_REDIRECT_URI")
		t.Setenv("LINE_LOGIN_CLIENT_ID", "LINE_LOGIN_CLIENT_ID")
		t.Setenv("LINE_LOGIN_CLIENT_SECRET", "LINE_LOGIN_CLIENT_SECRET")
		t.Setenv("LINE_LOGIN_API_URL", "LINE_LOGIN_API_URL")
		t.Setenv("HEALTHCHECK_REPORT_URL", "HEALTHCHECK_REPORT_URL")

		configAssert := newConfig("assets/example.csv")
		configExpect := config{
			File:                     "assets/example.csv",
			LINE_LOGIN_CODE:          "LINE_LOGIN_CODE",
			LINE_LOGIN_REDIRECT_URI:  "LINE_LOGIN_REDIRECT_URI",
			LINE_LOGIN_CLIENT_ID:     "LINE_LOGIN_CLIENT_ID",
			LINE_LOGIN_CLIENT_SECRET: "LINE_LOGIN_CLIENT_SECRET",
			LINE_LOGIN_API_URL:       "LINE_LOGIN_API_URL",
			HEALTHCHECK_URL:          "HEALTHCHECK_REPORT_URL",
		}

		assert.EqualValues(t, configAssert, configExpect)
	})

	t.Run("it should return default env instead if not assign", func(t *testing.T) {
		t.Setenv("LINE_LOGIN_CODE", "LINE_LOGIN_CODE")
		t.Setenv("LINE_LOGIN_REDIRECT_URI", "LINE_LOGIN_REDIRECT_URI")
		t.Setenv("LINE_LOGIN_CLIENT_ID", "LINE_LOGIN_CLIENT_ID")
		t.Setenv("LINE_LOGIN_CLIENT_SECRET", "LINE_LOGIN_CLIENT_SECRET")
		t.Setenv("LINE_LOGIN_API_URL", "")
		t.Setenv("HEALTHCHECK_REPORT_URL", "")

		configAssert := newConfig("assets/example.csv")
		configExpect := config{
			File:                     "assets/example.csv",
			LINE_LOGIN_CODE:          "LINE_LOGIN_CODE",
			LINE_LOGIN_REDIRECT_URI:  "LINE_LOGIN_REDIRECT_URI",
			LINE_LOGIN_CLIENT_ID:     "LINE_LOGIN_CLIENT_ID",
			LINE_LOGIN_CLIENT_SECRET: "LINE_LOGIN_CLIENT_SECRET",
			LINE_LOGIN_API_URL:       DEFAULT_LINE_LOGIN_API_URL,
			HEALTHCHECK_URL:          DEFAULT_HEATHCHECK_URL,
		}

		assert.EqualValues(t, configAssert, configExpect)
	})
}

func TestCheckArgs(t *testing.T) {
	t.Run("it should valid", func(t *testing.T) {
		t.Setenv("LINE_LOGIN_CODE", "LINE_LOGIN_CODE")
		t.Setenv("LINE_LOGIN_REDIRECT_URI", "LINE_LOGIN_REDIRECT_URI")
		t.Setenv("LINE_LOGIN_CLIENT_ID", "LINE_LOGIN_CLIENT_ID")
		t.Setenv("LINE_LOGIN_CLIENT_SECRET", "LINE_LOGIN_CLIENT_SECRET")

		expectations := []struct {
			in  config
			out error
		}{
			{
				newConfig("assets/example.csv"),
				nil,
			},
		}
		for _, e := range expectations {
			assert.Equal(t, e.in.checkArgs(), e.out)
		}
	})

	t.Run("it should invalid", func(t *testing.T) {
		t.Setenv("LINE_LOGIN_CODE", "")
		t.Setenv("LINE_LOGIN_REDIRECT_URI", "")
		t.Setenv("LINE_LOGIN_CLIENT_ID", "")
		t.Setenv("LINE_LOGIN_CLIENT_SECRET", "")

		expectations := []struct {
			in  config
			out error
		}{
			{
				newConfig("does/not/exist.forreal"),
				errInvalidURL,
			},
			{
				newConfig("assets/example.csv"),
				errLineEnvBlank,
			},
		}
		for _, e := range expectations {
			assert.Equal(t, e.in.checkArgs(), e.out)
		}
	})
}
