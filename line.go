package main

import (
	"encoding/json"
	"net/url"
	"strings"
)

type lineAPIER interface {
	getAccessToken(siteTotal int, successTotal int, failTotal int, executeTime int64) (lineToken, error)
}

type lineAPI struct {
	Config  config
	Request httpHCER
}

type lineToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func newLineAPI(config config, request httpHCER) lineAPIER {
	return &lineAPI{
		Config:  config,
		Request: request,
	}
}

func (lapi *lineAPI) getAccessToken(siteTotal int, successTotal int, failTotal int, executeTime int64) (lineToken, error) {
	data := lineToken{}
	ch := make(chan httpResult)
	body := url.Values{}
	body.Set("grant_type", lineLoginGrantType)
	body.Set("code", lapi.Config.LINE_LOGIN_CODE)
	body.Set("redirect_uri", lapi.Config.LINE_LOGIN_REDIRECT_URI)
	body.Set("client_id", lapi.Config.LINE_LOGIN_CLIENT_ID)
	body.Set("client_secret", lapi.Config.LINE_LOGIN_CLIENT_SECRET)
	headers := map[string]string{
		"Content-Type": contentTypeUrlEncoded,
	}

	go lapi.Request.post(lapi.Config.LINE_LOGIN_API_URL, strings.NewReader(body.Encode()), headers, ch)

	resp := <-ch
	close(ch)
	if resp.err != nil {
		return data, resp.err
	}

	if resp.statusCode != 200 {
		return data, errorf("Status: %d, Fail to request line api access token.", resp.statusCode)
	}

	json.Unmarshal(resp.body, &data)

	return data, nil
}
