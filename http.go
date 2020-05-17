package freshdesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const httpClientTimeout = time.Second * 10

func (c *ApiClient) postJSON(path string, requestBody []byte, out interface{}, expectedStatus int) error {
	httpClient := &http.Client{
		Timeout: httpClientTimeout,
	}
	if c.logger != nil {
		c.logger.Println(string(requestBody))
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://%s.freshservice.com%s", c.domain, path), bytes.NewReader(requestBody))

	req.SetBasicAuth(c.apiKey, "X")
	req.Header.Add("Content-type", "application/json")
	//c.logReq(req)

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != expectedStatus {
		body, err := ioutil.ReadAll(res.Body)
		apiError := ""
		if err == nil {
			var jsonBuffer bytes.Buffer
			err := json.Indent(&jsonBuffer, body, "", "\t")
			if err == nil {
				apiError = string(jsonBuffer.Bytes())
			}
		}
		if res.StatusCode == http.StatusBadRequest && c.logger != nil {
			if c.logger != nil {
				c.logger.Println(apiError)
			}
		}
		return APIError{
			fmt.Errorf("received status code %d (%d expected)", res.StatusCode, expectedStatus),
			apiError,
		}
	}

	err = json.NewDecoder(res.Body).Decode(out)

	return err
}

func (c *ApiClient) put(path string, requestBody []byte, out interface{}, expectedStatus int) error {
	httpClient := &http.Client{
		Timeout: httpClientTimeout,
	}
	if c.logger != nil {
		c.logger.Println(string(requestBody))
	}
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%s.freshservice.com%s", c.domain, path), bytes.NewReader(requestBody))

	req.SetBasicAuth(c.apiKey, "X")
	req.Header.Add("Content-type", "application/json")
	//c.logReq(req)

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != expectedStatus {
		body, err := ioutil.ReadAll(res.Body)
		var apiError string
		if err == nil {
			var jsonBuffer bytes.Buffer
			err := json.Indent(&jsonBuffer, body, "", "\t")
			if err == nil {
				apiError = string(jsonBuffer.Bytes())
			}
		}
		if res.StatusCode == http.StatusBadRequest && c.logger != nil {
			if c.logger != nil {
				c.logger.Println(apiError)
			}
		}
		return APIError{
			fmt.Errorf("received status code %d (%d expected)", res.StatusCode, expectedStatus),
			apiError,
		}
	}

	err = json.NewDecoder(res.Body).Decode(out)

	return err
}

func (c *ApiClient) get(path string, out interface{}) (http.Header, error) {
	httpClient := &http.Client{
		Timeout: httpClientTimeout,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s.freshservice.com%s", c.domain, path), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiKey, "X")
	if c.logger != nil {
		c.logger.Println(path)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	c.logRes(res)

	// body, err := ioutil.ReadAll(res.Body)
	// var f interface{}
	// err = json.Unmarshal(body, &f)
	// if err != nil {
	// 	log.Printf("Unable to Unmarshal request")
	// }

	// data := f.(map[string]interface{})
	// log.Println(data)

	json.NewDecoder(res.Body).Decode(out)

	return res.Header, err
}

func (c *ApiClient) getNextLink(headers http.Header) string {
	link := headers.Get("link")
	if link != "" {
		return strings.TrimPrefix(strings.TrimSuffix(link, ">; rel=\"next\""), fmt.Sprintf("<https://%s.freshservice.com", c.domain))
	}
	return ""
}

func (c *ApiClient) delete(path string) error {
	httpClient := &http.Client{
		Timeout: httpClientTimeout,
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://%s.freshservice.com%s", c.domain, path), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.apiKey, "X")
	req.Header.Add("Content-type", "application/json")
	//c.logReq(req)

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("received status code %d (200 expected)", res.StatusCode)
	}

	return nil
}
