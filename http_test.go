package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	if c.DoFunc != nil {
		return c.DoFunc(req)
	}
	return &http.Response{}, nil
}

func TestHttpHcGetDoSuccess(t *testing.T) {
	ch := make(chan httpResult)
	wg := &sync.WaitGroup{}
	clientMock := &ClientMock{}
	clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
		}, nil
	}
	h2p := newHttp(clientMock)

	wg.Add(1)

	go h2p.get("some-url", ch, wg)

	expectation := <-ch

	assert.Equal(t, expectation.statusCode, 200)
	assert.Equal(t, expectation.err, nil)

	wg.Wait()
}

func TestHttpHcGetDoFail(t *testing.T) {
	ch := make(chan httpResult)
	wg := &sync.WaitGroup{}
	clientMock := &ClientMock{}
	clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New(
			"Error from web server",
		)
	}
	h2p := newHttp(clientMock)

	wg.Add(1)

	go h2p.get("some-url", ch, wg)

	expectation := <-ch

	assert.NotEqual(t, expectation.err, nil)

	wg.Wait()
}

func TestHttpHcPostDoSuccess(t *testing.T) {
	ch := make(chan httpResult)
	clientMock := &ClientMock{}
	content := `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}}`
	body := ioutil.NopCloser(bytes.NewReader([]byte(content)))
	clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       body,
		}, nil
	}
	h2p := newHttp(clientMock)

	go h2p.post("", nil, nil, ch)

	expectation := <-ch

	assert.Equal(t, string(expectation.body), content)
	assert.Equal(t, expectation.statusCode, 200)
	assert.Equal(t, expectation.err, nil)
}

func TestHttpHcPostDoFailReadBody(t *testing.T) {
	ch := make(chan httpResult)
	clientMock := &ClientMock{}

	clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		}))
		defer ts.Close()

		return http.Get(ts.URL)
	}
	h2p := newHttp(clientMock)

	go h2p.post("", nil, nil, ch)

	expectation := <-ch

	assert.NotEqual(t, expectation.err, nil)
}

func TestHttpHcPostDoFail(t *testing.T) {
	ch := make(chan httpResult)
	clientMock := &ClientMock{}
	clientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New(
			"Error from web server",
		)
	}
	h2p := newHttp(clientMock)

	go h2p.post("", nil, nil, ch)

	expectation := <-ch

	assert.NotEqual(t, expectation.err, nil)
}
