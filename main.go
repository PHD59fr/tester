package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/fatih/color"
)

type EndpointConfig struct {
	Name             string                 `json:"name"`
	URL              string                 `json:"url"`
	Method           string                 `json:"method"`
	Headers          map[string]string      `json:"headers"`
	Body             map[string]interface{} `json:"body"`
	ExpectedStatus   int                    `json:"expectedStatus"`
	ExpectedResponse map[string]interface{} `json:"expectedResponse"`
}

type Config struct {
	Endpoints []EndpointConfig `json:"endpoints"`
}

func main() {
	configFile := flag.String("testfile", "", "Specify the test configuration file")
	flag.Parse()

	configData, err := os.ReadFile(*configFile)
	if err != nil {
		color.Red("Error reading the test configuration file: %v\n", err)
		return
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		color.Red("Error decoding the JSON test configuration file: %v\n", err)
		return
	}

	passCount := 0
	failCount := 0

	for _, endpoint := range config.Endpoints {
		request, err := http.NewRequest(endpoint.Method, endpoint.URL, nil)
		if err != nil {
			color.Red("[FAIL] [%s] creating the request: %v\n", endpoint.Name, err)
			failCount++
			continue
		}

		for key, value := range endpoint.Headers {
			request.Header.Set(key, value)
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			color.Red("[FAIL] [%s] making the request to the endpoint: %v\n", endpoint.Name, err)
			failCount++
			continue
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)

		if response.StatusCode != endpoint.ExpectedStatus {
			color.Red("[FAIL] [%s] the expected code does not match the received code. Expected: %d, Received: %d\n", endpoint.Name, endpoint.ExpectedStatus, response.StatusCode)
			failCount++
			continue
		}

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {

			color.Red("[FAIL] [%s] reading the endpoint's response: %v\n", endpoint.Name, err)
			failCount++
			continue
		}
		var responseJSON map[string]interface{}
		if endpoint.ExpectedResponse != nil {
			if err := json.Unmarshal(responseBody, &responseJSON); err != nil {
				color.Red("[FAIL] [%s] decoding the JSON response: %v\n", endpoint.Name, err)
				failCount++
				continue
			}
		}

		if !reflect.DeepEqual(responseJSON, endpoint.ExpectedResponse) {
			color.Red("[FAIL] [%s] unexpected response. Expected: %v, Received: %v\n", endpoint.Name, endpoint.ExpectedResponse, responseJSON)
			failCount++
			continue
		}

		color.Green("[PASS] [%s]\n", endpoint.Name)
		passCount++
	}
	fmt.Printf("Tests Passed: %d / Tests Failed: %d / Coverage: %.2f%%\n", passCount, failCount, float64(passCount)/float64(len(config.Endpoints))*100)

}
