# Website-Copier in Golang

## Installation
Ensure you have Go 1.23.5 installed on your system.

Clone the repository:
```sh
git clone https://github.com/b0tShaman/website-copier.git
cd website-copier
```

Install dependencies:
```sh
go mod tidy
```

## Dependencies
This project uses the following external dependencies:
- [logrus](https://github.com/sirupsen/logrus): For structured logging.

## Usage
### 1. Run the Program
Build and run the Go program:
```sh
go run main.go urls.csv
```

### 2. Output
- Each downloaded webpage is saved as a separate `.txt` file with a random name in the `output` directory.
- The `output` directory will be created automatically if it does not exist.

## Running Unit Tests
To run the unit tests, use the following command:
```sh
go test ./...
```
This will execute all tests within the project.

## Assumptions and Design Decisions
- The number of concurrent download goroutines is capped at **50** with 1 additional goroutine to control all 50.
- The program dynamically adjusts the number of goroutines based on the number of URLs.
- Graceful shutdown is implemented to allow the program to exit cleanly within 5 seconds when interrupted.


