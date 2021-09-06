package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Call invokes the Management api with the relative path: e.g. http://localhost:15672/api/[relPath]
func Call(pathAndQuery, method string) (*http.Response, error) {
	url := Config.url(pathAndQuery)
	return Config.call(url, method)
}

func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (cfg *cfg) url(pathAndQuery string) string {
	pathAndQuery = strings.TrimLeft(pathAndQuery, "/")
	return fmt.Sprintf("%s://%s:%d/api/%s", cfg.Scheme, cfg.Host, cfg.Port, pathAndQuery)
}

func (cfg *cfg) call(url, method string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}

	req.SetBasicAuth(cfg.User, cfg.Password)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	traceRequest(req)
	resp, err := client.Do(req)
	traceResponse(resp, err)
	return resp, err
}
