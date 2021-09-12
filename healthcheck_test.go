package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type lineAPIMock struct {
	DoFunc func(siteTotal int, successTotal int, failTotal int, executeTime int64) (lineToken, error)
}

func (c *lineAPIMock) getAccessToken(siteTotal int, successTotal int, failTotal int, executeTime int64) (lineToken, error) {
	if c.DoFunc != nil {
		return c.DoFunc(siteTotal, successTotal, failTotal, executeTime)
	}
	return lineToken{}, nil
}

func TestReadFile(t *testing.T) {
	t.Run("it should has file exist and not error", func(t *testing.T) {
		client := &http.Client{}
		config := newConfig("assets/example.csv")
		httpHc := newHttp(client)
		lineAPI := newLineAPI(config, httpHc)
		hc := newHealthCheck(config, httpHc, lineAPI)

		_, err := hc.readFile()
		assert.Equal(t, err, nil)
	})

	t.Run("it should has not file exist and error", func(t *testing.T) {
		client := &http.Client{}
		config := newConfig("does/not/exist.forreal")
		httpHc := newHttp(client)
		lineAPI := newLineAPI(config, httpHc)
		hc := newHealthCheck(config, httpHc, lineAPI)

		_, err := hc.readFile()
		assert.NotEqual(t, err, nil)
	})
}

func TestCheckHealthWebsite(t *testing.T) {
	t.Run("it should valid checking site total and success total", func(t *testing.T) {
		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
			}, nil
		}
		config := newConfig("")
		httpHc := newHttp(clientMock)
		lineAPI := newLineAPI(config, httpHc)
		hc := newHealthCheck(config, httpHc, lineAPI)

		contents := [][]string{
			{"https://www.facebook.com"},
			{"https://www.google.com"},
			{"https://www.youtube.com"},
		}

		checkingSiteTotal, successTotal, failTotal, _ := hc.checkHealthWebsite(contents)

		assert.Equal(t, checkingSiteTotal, 3)
		assert.Equal(t, successTotal, 3)
		assert.Equal(t, failTotal, 0)
	})

	t.Run("it should valid checking site total and fail total", func(t *testing.T) {
		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New(
				"Error from web server",
			)
		}
		config := newConfig("")
		httpHc := newHttp(clientMock)
		lineAPI := newLineAPI(config, httpHc)
		hc := newHealthCheck(config, httpHc, lineAPI)

		contents := [][]string{
			{"https://www.facebook.com"},
			{"https://www.google.com"},
			{"https://www.youtube.com"},
		}

		checkingSiteTotal, successTotal, failTotal, _ := hc.checkHealthWebsite(contents)

		assert.Equal(t, checkingSiteTotal, 3)
		assert.Equal(t, successTotal, 0)
		assert.Equal(t, failTotal, 3)
	})

	t.Run("it should valid only once url", func(t *testing.T) {
		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
			}, nil
		}
		config := newConfig("")
		httpHc := newHttp(clientMock)
		lineAPI := newLineAPI(config, httpHc)
		hc := newHealthCheck(config, httpHc, lineAPI)

		contents := [][]string{
			{"abc"},
			{"https://www.google.com"},
			{"def"},
		}

		checkingSiteTotal, successTotal, failTotal, _ := hc.checkHealthWebsite(contents)

		assert.Equal(t, checkingSiteTotal, 1)
		assert.Equal(t, successTotal, 1)
		assert.Equal(t, failTotal, 0)
	})
}

func TestReportStatistic(t *testing.T) {
	t.Run("it should error when call request", func(t *testing.T) {

		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New(
				"Error from web server",
			)
		}
		config := newConfig("")
		httpHc := newHttp(clientMock)
		lineAPI := newLineAPI(config, httpHc)
		hc := newHealthCheck(config, httpHc, lineAPI)

		err := hc.reportStatistic(1, 1, 1, 1)

		assert.NotEqual(t, err, nil)
	})

	t.Run("it should error when call line access token api", func(t *testing.T) {
		expectation := errors.New(
			"Error from get accesstoken",
		)
		clientMock := &ClientMock{}
		config := newConfig("")
		httpHc := newHttp(clientMock)
		lineMock := &lineAPIMock{}
		lineMock.DoFunc = func(siteTotal int, successTotal int, failTotal int, executeTime int64) (lineToken, error) {
			return lineToken{}, expectation
		}
		hc := newHealthCheck(config, httpHc, lineMock)

		err := hc.reportStatistic(1, 1, 1, 1)

		assert.Equal(t, err, expectation)
	})

	t.Run("it should error when call line access token api and return any status which not 200", func(t *testing.T) {

		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}
		config := newConfig("")
		httpHc := newHttp(clientMock)
		lineMock := &lineAPIMock{}
		hc := newHealthCheck(config, httpHc, lineMock)

		err := hc.reportStatistic(1, 1, 1, 1)

		assert.Equal(t, err, errorf("Status: %d, Fail to request health check.", 400))
	})

	t.Run("it should valid without error", func(t *testing.T) {
		clientMock := &ClientMock{}
		clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}
		config := newConfig("")
		httpHc := newHttp(clientMock)
		lineMock := &lineAPIMock{}
		hc := newHealthCheck(config, httpHc, lineMock)

		err := hc.reportStatistic(1, 1, 1, 1)

		assert.Equal(t, err, nil)
	})
}
