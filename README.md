# Tester

This is a simple command-line tool for testing API calls in different environments, including development, staging, pre-production, and production. It is designed to help you verify the functionality and correctness of various API endpoints.

## Features

- Load API test configurations from a JSON file.
- Send HTTP requests to specified endpoints.
- Validate response codes and expected JSON responses.

## Getting Started

### Quick start
Run the tool with a test configuration file:

```shell
./tester -testfile sample.json
```

## Usage

To use the tester, follow these steps:

1. Create a JSON configuration file specifying the API endpoints you want to test. See `example-config.json` for a sample configuration format.

2. Run the tool with the `-testfile` flag to specify the configuration file:

```shell
./tester -testfile your-config-file.json
```

3. The tool will send HTTP requests to the specified endpoints and report the results.

## Configuration File

The configuration file should contain an array of test cases, each specifying the following:

- `name`: A descriptive name for the test case.
- `url`: The API endpoint URL.
- `method`: HTTP request method (e.g., GET, POST).
- `headers`: Key-value pairs of HTTP headers.
- `body`: Request body as a JSON object.
- `expectedStatus`: The expected HTTP response status code.
- `expectedResponse`: The expected JSON response (can be `null` if not required).

```json
{
  "endpoints": [
    {
      "name": "Test json",
      "url": "https://jsonplaceholder.typicode.com/todos/1",
      "method": "GET",
      "headers": {"User-Agent":  "Mozilla/5.0"},
      "body": {},
      "expectedStatus": 200,
      "expectedResponse": {"userId": 1, "id": 1,"title":"delectus aut autem","completed": false}
    },
    {
      "name": "Test ip",
      "url": "https://ipaddr.ovh",
      "method": "GET",
      "headers": {"Accept-Encoding": "gzip, deflate, br"},
      "body": {},
      "expectedStatus": 200
    }
    /// Your test here !
  ]
}

```

## Contributing

Contributions to this project are welcome. If you encounter any issues, have suggestions, or want to add features, please open an issue or create a pull request.