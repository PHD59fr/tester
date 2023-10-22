package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"syscall"
	"text/template"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type EndpointTest struct {
	Name              string                 `yaml:"name"`
	URL               string                 `yaml:"url"`
	Method            string                 `yaml:"method"`
	Headers           map[string]string      `yaml:"headers"`
	Body              map[string]interface{} `yaml:"body"`
	ExpectedStatus    int                    `yaml:"expectedStatus"`
	ExpectedResponse  map[string]interface{} `yaml:"expectedResponse"`
	ResponseVariables map[string]string      `yaml:"responseVariables"`
}

type TestScenario struct {
	Endpoints []EndpointTest `yaml:"endpoints"`
}

func main() {
	testScenarioFile := flag.String("testFile", "", "Specify the test scenario file")
	showDetails := flag.Bool("details", false, "Display request and response details")
	ignoreFail := flag.Bool("ignoreFail", true, "Specify whether to display request and response details in case of failure")
	flag.Parse()

	testScenario, err := loadTestScenario(*testScenarioFile)
	if err != nil {
		color.Red("Error: %v\n", err)
		return
	}

	passCount, failCount := 0, 0
	responseVariables := make(map[string]interface{})

	for _, endpoint := range testScenario.Endpoints {
		if err := processEndpoint(&endpoint, showDetails, responseVariables); err != nil {
			color.Red("[FAIL] [%s] %v\n", endpoint.Name, err)
			failCount++
			if !*ignoreFail {
				printTestSummary(passCount, failCount, len(testScenario.Endpoints))
				fmt.Print("The other tests are ignored because you have specified the ignoreFail flag as false\n")

				syscall.Exit(2)
			}
		} else {
			color.Green("[PASS] [%s]\n", endpoint.Name)
			passCount++
		}
	}

	printTestSummary(passCount, failCount, len(testScenario.Endpoints))
}

func loadTestScenario(filename string) (*TestScenario, error) {
	testScenarioData, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Annotate(err, "error reading the test scenario file")
	}

	var testScenario TestScenario
	if err := yaml.Unmarshal(testScenarioData, &testScenario); err != nil {
		return nil, errors.Annotate(err, "error decoding the YAML test scenario file")
	}

	return &testScenario, nil
}

func processEndpoint(endpoint *EndpointTest, showDetails *bool, responseVariables map[string]interface{}) error {
	replaceVariablesInEndpoint(endpoint, responseVariables)
	request, err := createRequest(endpoint)
	if err != nil {
		return errors.Annotate(err, fmt.Sprintf("[FAIL] [%s] %v", endpoint.Name, err))
	}

	if *showDetails {
		if err := dumpRequest(request); err != nil {
			return errors.Annotate(err, fmt.Sprintf("[FAIL] [%s] %v", endpoint.Name, err))
		}
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return errors.Annotate(err, "making the request to the endpoint")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if *showDetails {
		if err := dumpResponse(response); err != nil {
			return errors.Annotate(err, "dumping the response")
		}
	}

	if response.StatusCode != endpoint.ExpectedStatus {
		return errors.New(fmt.Sprintf("expected code %d, received %d", endpoint.ExpectedStatus, response.StatusCode))
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Annotate(err, "reading the endpoint's response")
	}

	if endpoint.ResponseVariables != nil && endpoint.ExpectedResponse != nil {
		var responseJSON map[string]interface{}
		if err := yaml.Unmarshal(responseBody, &responseJSON); err != nil {
			return errors.Annotate(err, "decoding the YAML response")
		}
		if err := checkExpectedResponse(responseJSON, endpoint.ExpectedResponse); err != nil {
			return err
		}
	}

	if endpoint.Body != nil {
		bodyStr := fmt.Sprintf("%v", endpoint.Body)
		bodyWithVariables := replaceVariables(bodyStr, responseVariables)
		request.Body = io.NopCloser(strings.NewReader(bodyWithVariables))
	}

	return nil
}

func replaceVariablesInEndpoint(endpoint *EndpointTest, responseVariables map[string]interface{}) {
	endpoint.URL = replaceVariables(endpoint.URL, responseVariables)
	for key, value := range endpoint.Headers {
		endpoint.Headers[key] = replaceVariables(value, responseVariables)
	}
}

func createRequest(endpoint *EndpointTest) (*http.Request, error) {
	endpointBody, err := json.Marshal(endpoint.Body)
	if err != nil {
		return nil, err
	}

	if endpoint.Body == nil {
		endpointBody = nil
	}

	request, err := http.NewRequest(endpoint.Method, endpoint.URL, bytes.NewReader(endpointBody))
	if err != nil {
		return nil, err
	}

	for key, value := range endpoint.Headers {
		request.Header.Set(key, value)
	}

	return request, nil
}

func dumpRequest(request *http.Request) error {
	requestDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return err
	}
	requestDumpLines := strings.Split(string(requestDump), "\n")
	for _, line := range requestDumpLines {
		color.Cyan("> %s", line)
	}
	return nil
}

func dumpResponse(response *http.Response) error {
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		return err
	}
	responseDumpLines := strings.Split(string(responseDump), "\n")
	for _, line := range responseDumpLines {
		color.Cyan("< %s", line)
	}
	return nil
}

func replaceVariables(input string, variables map[string]interface{}) string {
	tmpl, err := template.New("variables").Parse(input)
	if err != nil {
		return input
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, variables)
	if err != nil {
		return input
	}

	return output.String()
}

func checkExpectedResponse(actualResponse, expectedResponse map[string]interface{}) error {
	for key, expectedValue := range expectedResponse {
		actualValue, exists := actualResponse[key]
		if !exists {
			return errors.NotFoundf(`response key '%s' not found in the actual response`, key)
		}

		if fmt.Sprintf("%v", actualValue) != fmt.Sprintf("%v", expectedValue) {
			return errors.New(fmt.Sprintf(`response key '%s' does not match the expected value. Expected: %v, Actual: %v`, key, expectedValue, actualValue))
		}
	}
	return nil
}

func printTestSummary(passCount, failCount, totalTests int) {
	fmt.Printf("Tests Passed: %d / Tests Failed: %d / Coverage: %.2f%%\n", passCount, failCount, float64(passCount)/float64(totalTests)*100)
}
