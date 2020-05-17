package freshdesk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ApiClient struct {
	domain      string
	apiKey      string
	logger      *log.Logger
	Departments DepartmentManager
	Requesters  RequesterManager
	Tickets     TicketManager
}

type ClientOptions struct {
	Logger *log.Logger
}

func EmptyOptions() *ClientOptions {
	return nil
}

// Init initializes the package
func Init(domain, apiKey string, options *ClientOptions) ApiClient {
	client := ApiClient{
		domain: domain,
		apiKey: apiKey,
	}
	if options != nil {
		client.logger = options.Logger
	}
	if client.logger != nil {
		client.logger.Println("Freshservice Client initializing... Domain =", domain, "authorization =", apiKey)
	}
	client.Departments = newDepartmentManager(&client)
	client.Tickets = newTicketManager(&client)
	client.Requesters = newrequesterManager(&client)
	return client
}

func (client ApiClient) logErr(err error) {
	if err != nil && client.logger != nil {
		client.logger.Println(err.Error())
	}
}

func (client ApiClient) logReq(req *http.Request) {
	if client.logger != nil {
		client.logger.Println("Headers")
		for k, v := range req.Header {
			client.logger.Printf("%s: %s\n", k, v)
		}
		client.logger.Println("URL:", req.URL)
		if req.Body != nil {
			body, _ := ioutil.ReadAll(req.Body)
			var jsonBuffer bytes.Buffer
			json.Indent(&jsonBuffer, body, "", "\t")
			client.logger.Println(string(jsonBuffer.Bytes()))
		}
	}
}

func (client ApiClient) logRes(res *http.Response) {
	if client.logger != nil {
		client.logger.Println("Status:", res.StatusCode)
		if res.StatusCode != 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				client.logger.Println(err.Error())
			}
			var jsonBuffer bytes.Buffer
			json.Indent(&jsonBuffer, body, "", "\t")
			client.logger.Println(string(jsonBuffer.Bytes()))
		}
	}
}
