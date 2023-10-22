# Tester

This is a simple command-line tool for testing API calls in different environments, including development, staging, pre-production, and production. It is designed to help you verify the functionality and correctness of various API endpoints.

## Features

- Load API test scenario file from a YAML file.
- Send HTTP requests to specified endpoints.
- Validate response codes and expected JSON responses.

## Getting Started

### Quick start
Run the tool with a test scenario file:

```shell
./tester -testFile sample.yaml
```

## Usage

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
- `headers`: Key-value pairs of HTTP headers (Optional field).
- `body`: Request body as a map (usually for JSON requests, Optional field).
- `expectedStatus`: The expected HTTP response status code.
- `expectedResponse`: The expected JSON response data (Optional field).
- `responseVariables`: A map of response variables to capture from the response and use in subsequent requests (Optional field).

Example YAML format:
```yaml
endpoints:
  - name: "Test Post product"
    url: "https://dummyjson.com/product/add"
    method: POST
    headers:
      User-Agent: "Mozilla/5.0"
    body:
      kikou: "test"
    expectedStatus: 200
    responseVariables:
      id: "id"

  - name: "Test json responseVariables"
    url: "https://dummyjson.com/product/{{.id}}"
    method: GET
    headers:
      User-Agent: "Mozilla/5.0"
    expectedStatus: 404
    expectedResponse:
      message: "Product with id '101' not found"

  - name: "Test ip"
    url: "https://ipaddr.ovh"
    method: GET
    headers:
      User-Agent: "Mozilla/5.0"
    expectedStatus: 200
```

## Contributing

Contributions to this project are welcome. If you encounter any issues, have suggestions, or want to add features, please open an issue or create a pull request. Your contributions will help improve this automation tool for testing API scenarios.