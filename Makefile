.PHONY: test coverage clean build

# Build the project
build:
	go build -o bin/task-tracker main.go

# Run all tests
test:
	go test ./... -v

# Generate coverage report
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean generated files
clean:
	rm -f coverage.out coverage.html
	rm -f bin/task-tracker

# Run tests and coverage
test-coverage: test coverage 
