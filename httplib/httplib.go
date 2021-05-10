package httplib

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func Get(url string) (int, []byte, error) {
	return parseResponse(makeRequest("GET", url, nil))
}
func Post(url string, requestBody []byte) (int, []byte, error) {
	return parseResponse(makeRequest("POST", url, bytes.NewBuffer(requestBody)))
}

func Put(url string, requestBody []byte) (int, []byte, error) {
	return parseResponse(makeRequest("PUT", url, bytes.NewBuffer(requestBody)))
}

func Delete(url string, requestBody []byte) (int, []byte, error) {
	return parseResponse(makeRequest("DELETE", url, bytes.NewBuffer(requestBody)))
}

func makeRequest(method string, url string, requestBody io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, requestBody)
	req.Header.Set("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(req)
	return response, err
}
func parseResponse(response *http.Response, err error) (int, []byte, error) {
	if err != nil {
		return -1, nil, err
	}
	body := response.Body
	defer body.Close()
	if response.StatusCode == http.StatusOK {
		responseBody, err := ioutil.ReadAll(body)
		if err != nil {
			return response.StatusCode, nil, err
		} else {
			return response.StatusCode, responseBody, nil
		}
	} else {
		return response.StatusCode, nil, nil
	}
}
