# ICS Generator Service

This service generates ICS (iCalendar) files from JSON input data. It can run in both server and local modes, providing flexibility for different use cases.

## Purpose

The ICS Generator Service simplifies the process of creating iCalendar files for various events, holidays, or schedules. It accepts JSON-formatted event data and produces standard ICS files that can be imported into most calendar applications.

## Features

-   Supports both server (including containerized) and local operation modes
-   Handles single-day and multi-day events
-   Processes JSON input for easy integration with existing systems
-   Generates standard ICS files compatible with popular calendar applications

## Project Structure

```
ics-generator-service/
├── cmd/
│   └── server/
│       ├── main.go
│       └── main_test.go
├── input/
├── internal/
│   └── ics/
│       ├── generator.go
│       └── generator_test.go
├── output/
├── testdata/
│   ├── mixed_holidays_2024.json
│   └── public_holidays_2024.json
├── .dockerignore
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Key Components

### Main Application (cmd/server/)

-   `main.go`: The entry point of the application. It handles command-line flag parsing and determines whether to run in server or local mode.
-   `main_test.go`: Contains tests for the main application, including flag parsing.

### ICS Generator (internal/ics/)

-   `generator.go`: Contains the core logic for generating ICS files from JSON input. Key functions include:
    -   `HandleGenerateICS`: Handles HTTP requests in server mode
    -   `GenerateLocalICS`: Generates ICS files from local JSON files
    -   `GenerateICS`: Core function for converting event data to ICS format
-   `generator_test.go`: Contains tests for the ICS generation logic.

### Input (input/)

This directory is used for storing input JSON files when running in local mode.

### Output (output/)

Generated ICS files are saved in this directory by default when running in local mode.

### Test Data (testdata/)

Contains JSON files for testing:

-   `mixed_holidays_2024.json`: Sample data with a mix of single-day and multi-day events
-   `public_holidays_2024.json`: Sample data for public holidays (single-day events)

Example JSON format:

```json
[
    {
        "id": "event-1",
        "title": "Event Title",
        "start": "2024-01-01",
        "end": "2024-01-01",
        "description": "Event description"
    }
    // ... more events
]
```

### Configuration Files

-   `.dockerignore`: Specifies which files and directories should be excluded when building the Docker image
-   `.gitignore`: Specifies intentionally untracked files to ignore in Git
-   `Dockerfile`: Contains instructions for building the Docker image
-   `go.mod` and `go.sum`: Go module files for managing dependencies
-   `LICENSE`: Contains the license information for the project

## Installation

1. Ensure you have Go installed (version 1.22.2 or later recommended).
2. Clone this repository:
    ```sh
    git clone https://github.com/lelopez-io/ics-generator-service.git
    cd ics-generator-service
    ```
3. Install dependencies:
    ```sh
    go mod download
    ```

## Running Tests

To run all tests in the project:

```sh
go test ./...
```

## Usage

### Command-line Flags

-   `--server`: Run in server mode (default is local mode)
-   `--input`: Specify the input JSON file (required for local mode)
-   `--output`: Specify the output ICS file (default is "output/calendar.ics", for local mode)
-   `--port`: Specify the port to run the server on (default is 8080, for server mode)

### Local File Mode

To generate an ICS file from a local JSON file:

```sh
go run cmd/server/main.go --input input/events.json --output output/calendar.ics
```

### Local Server Mode

To run the service as a local web server:

```sh
go run cmd/server/main.go --server
```

The server will be accessible at `http://localhost:8080`. To generate an ICS file, send a POST request to `http://localhost:8080/generate-ics` with JSON event data in the request body.

### Contained Server Mode

To build and run the Docker container in one command, with automatic cleanup:

```sh
docker run --rm -p 8080:8080 $(docker build -q -t ics-generator-service .)
```

The service will be accessible at `http://localhost:8080`. Use it the same way as the local server mode.

## Dependencies

This project uses the following external library:

-   `github.com/arran4/golang-ical` for ICS file generation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms of the LICENSE file in the root directory of this project. Please see the LICENSE file for full details.
