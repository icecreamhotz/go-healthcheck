package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type healthChecker interface {
	readFile() ([][]string, error)
	checkHealthWebsite(contents [][]string) (int, int, int, int64)
	reportStatistic(siteTotal int, successTotal int, failTotal int, executeTime int64) error
}

type healthCheck struct {
	Config  config
	Request httpHCER
	LineAPI lineAPIER
}

type statistic struct {
	TotalWebsites int   `json:"total_websites"`
	Success       int   `json:"success"`
	Failure       int   `json:"failure"`
	TotalTime     int64 `json:"total_time"`
}

func newHealthCheck(cfg config, request httpHCER, lineAPI lineAPIER) healthChecker {
	return &healthCheck{
		Config:  cfg,
		Request: request,
		LineAPI: lineAPI,
	}
}

func (hc *healthCheck) readFile() ([][]string, error) {
	csvFile, err := os.Open(hc.Config.File)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}

	return csvLines, nil
}

func (hc *healthCheck) checkHealthWebsite(contents [][]string) (int, int, int, int64) {
	timeStart := time.Now()
	ch := make(chan httpResult)
	wg := &sync.WaitGroup{}
	checkingSiteTotal := 0
	successTotal := 0
	failTotal := 0

	for _, content := range contents {
		if content[0] != "" {
			u, err := url.Parse(content[0])
			if err == nil && u.Scheme != "" && u.Host != "" {
				checkingSiteTotal += 1
				wg.Add(1)
				go hc.Request.get(content[0], ch, wg)
			}
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		if res.err == nil {
			successTotal += 1
		} else {
			failTotal += 1
		}
	}

	return checkingSiteTotal, successTotal, failTotal, time.Since(timeStart).Nanoseconds()
}

func (hc *healthCheck) reportStatistic(siteTotal int, successTotal int, failTotal int, executeTime int64) error {
	lineData, err := hc.LineAPI.getAccessToken(siteTotal, successTotal, failTotal, executeTime)
	if err != nil {
		return err
	}

	ch := make(chan httpResult)
	headers := map[string]string{
		"Authorization": lineData.AccessToken,
		"Content-Type":  contentTypeJSON,
	}
	data := &statistic{
		TotalWebsites: siteTotal,
		Success:       successTotal,
		Failure:       failTotal,
		TotalTime:     executeTime,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(data)

	go hc.Request.post(hc.Config.HEALTHCHECK_URL, payloadBuf, headers, ch)

	resp := <-ch
	close(ch)

	if resp.err != nil {
		return resp.err
	}

	if resp.statusCode != 200 {
		return errorf("Status: %d, Fail to request health check.", resp.statusCode)
	}

	return nil
}

func main() {
	parser := newKingpinParser()
	file, err := parser.parse(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg := newConfig(file)
	if err := cfg.checkArgs(); err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	httpHc := newHttp(client)
	lineAPI := newLineAPI(cfg, httpHc)
	hc := newHealthCheck(cfg, httpHc, lineAPI)

	contents, err := hc.readFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nPerform website checking...")

	siteTotal, successTotal, failTotal, executeTime := hc.checkHealthWebsite(contents)

	fmt.Println("Done!")
	fmt.Println("\nChecked webistes: ", siteTotal)
	fmt.Println("Successful websites: ", successTotal)
	fmt.Println("Failure websites: ", failTotal)
	fmt.Println("Total times to finished checking website: ", executeTime)
	fmt.Println("")

	err = hc.reportStatistic(siteTotal, successTotal, failTotal, executeTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Sended the statistic of each website succeed.")

}
