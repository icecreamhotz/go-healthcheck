package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type httpClienter interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpHCER interface {
	get(url string, ch chan<- httpResult, wg *sync.WaitGroup)
	post(url string, data io.Reader, headers map[string]string, ch chan<- httpResult)
}

type httpResult struct {
	statusCode int
	body       []byte
	err        error
}

type httpHc struct {
	Client httpClienter
}

func newHttp(client httpClienter) httpHCER {
	return &httpHc{
		Client: client,
	}
}

func (h2p *httpHc) get(url string, ch chan<- httpResult, wg *sync.WaitGroup) {
	defer wg.Done()

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ch <- httpResult{
			statusCode: 0,
			body:       nil,
			err:        err,
		}
		return
	}

	resp, err := h2p.Client.Do(r)
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
		body:       nil,
		err:        nil,
	}
}

func (h2p *httpHc) post(url string, data io.Reader, headers map[string]string, ch chan<- httpResult) {
	r, err := http.NewRequest(http.MethodPost, url, data)
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
			r.Header.Add(key, value)
		}
	}

	resp, err := h2p.Client.Do(r)
	if err != nil {
		ch <- httpResult{
			statusCode: 0,
			body:       nil,
			err:        err,
		}
		return
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
