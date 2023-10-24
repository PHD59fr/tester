# Tester

This is a simple command-line tool for testing API calls in different environments, including development, staging, pre-production, and production. It is designed to help you verify the functionality and correctness of various API endpoints.

## Features

- Load API test scenario file from a YAML file.
- Send HTTP requests to specified endpoints.
- Validate response codes and expected JSON responses.

## Getting Started

### Quick Start
Run the tool with a test scenario file:

```shell
./tester -testFile sample.yaml
```

## Usage

Command-line options:

```shell
./tester -h
  -details
        Display request and response details
  -stopOnFailure
        Specify whether to display request and response details in case of failure
  -testFile string
        Specify the test scenario file
```

To use the tester, follow these steps:

1. Create a YAML configuration file specifying the API endpoints you want to test. See `sample.yaml` for a sample format.

2. Run the tool with the `-testFile` flag to specify the test scenario file:

```shell
./tester -testFile your-config-file.yaml -details
```

3. The tool will send HTTP requests to the specified endpoints and report the results.

## Test Scenario Configuration (YAML)

The test scenario file should contain an array of test cases, each specifying the following:

- `name`: A descriptive name for the test case.
- `url`: The API endpoint URL.
- `method`: HTTP request method (e.g., GET, POST, PUT, DELETE).
- `headers`: Key-value pairs of HTTP headers (optional).
- `body`: Request body, which can contain various types of data (optional).
- `multipartFields`: Multipart fields for requests with files (optional).
- `expectedStatus`: The expected HTTP response status code.
- `expectedResponse`: The expected response data, which can be in different formats (optional).
- `responseVariables`: A map of response variables to capture from the response and use in subsequent requests (optional).

**Example YAML format**:

```yaml
endpoints:
  - name: "Test Multipart"
    url: "https://dummyjson.com/product/add"
    method: POST
    headers:
      User-Agent: "Mozilla/5.0"
    multipartFields:
      file: "@test.png"
    expectedStatus: 200

  - name: "Test Post"
    url: "https://dummyjson.com/product/add"
    method: POST
    headers:
      User-Agent: "Mozilla/5.0"
    body:
      kikou: "test"
    expectedStatus: 200
    responseVariables:
      id: "id"

  - name: "Test responseVariables and return 404"
    url: "https://dummyjson.com/product/{{.id}}"
    method: GET
    headers:
      User-Agent: "Mozilla/5.0"
    expectedStatus: 404
    expectedResponse:
      message: "Product with id '101' not found"

  - name: "Test Get simple"
    url: "https://ipaddr.ovh"
    method: GET
    headers:
      User-Agent: "Mozilla/5.0"
    expectedStatus: 200
```

## Contributing

Contributions to this project are welcome. If you encounter any issues, have suggestions, or want to add features, please open an issue or create a pull request. Your contributions will help improve this automation tool for testing API scenarios.
