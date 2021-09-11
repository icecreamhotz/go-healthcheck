package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type httpHCER interface {
	get(url string, ch chan<- httpResult, wg *sync.WaitGroup)
	post(url, contentType string, data io.Reader, headers map[string]string, ch chan<- httpResult)
}

type httpResult struct {
	statusCode int
	body       []byte
	err        error
}

type httpHc struct{}

func newHttp() httpHCER {
	return &httpHc{}
}

func (h2p *httpHc) get(url string, ch chan<- httpResult, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		ch <- httpResult{
			statusCode: 0,
			body:       nil,
			err:        err,
		}
		return
	}
	defer resp.Body.Close()

	ch <- httpResult{
		statusCode: resp.StatusCode,
		body:       nil,
		err:        nil,
	}
}

func (h2p *httpHc) post(url string, contentType string, data io.Reader, headers map[string]string, ch chan<- httpResult) {
	resp, err := http.Post(url, contentType, data)
	if err != nil {
		ch <- httpResult{
			statusCode: 0,
			body:       nil,
			err:        err,
		}
		return
	}
	if len(headers) > 0 {
		for key, value := range headers {
			resp.Header.Set(key, value)
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- httpResult{
			statusCode: 0,
			body:       nil,
			err:        err,
		}
		return
	}

	ch <- httpResult{
		statusCode: resp.StatusCode,
		body:       body,
		err:        nil,
	}
}
