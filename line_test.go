package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccessToken(t *testing.T) {
	t.Run("it should show error when request", func(t *testing.T) {
		config := newConfig("")
		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New(
				"Error from web server",
			)
		}
		h2p := newHttp(clientMock)
		lineAPI := newLineAPI(config, h2p)

		_, err := lineAPI.getAccessToken(1, 1, 1, 1)
		assert.NotEqual(t, err, nil)
	})

	t.Run("it should valid request and not return 200 status", func(t *testing.T) {
		config := newConfig("")
		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}
		h2p := newHttp(clientMock)
		lineAPI := newLineAPI(config, h2p)

		_, err := lineAPI.getAccessToken(1, 1, 1, 1)
		assert.Equal(t, err, errorf("Status: %d, Fail to request line api access token.", 400))
	})

	t.Run("it should valid request and return correct response", func(t *testing.T) {
		config := newConfig("")
		clientMock := &ClientMock{}
		data := `{"access_token": "token", "token_type":"bearer"}`
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(data))),
			}, nil
		}
		h2p := newHttp(clientMock)
		lineAPI := newLineAPI(config, h2p)

		actual, _ := lineAPI.getAccessToken(1, 1, 1, 1)
		expectation := lineToken{
			AccessToken: "token",
			TokenType:   "bearer",
		}
		assert.EqualValues(t, expectation, actual)
	})

}
