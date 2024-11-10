BINARY_NAME=your-app-name
MAIN_PACKAGE=./cmd/main.go
GO=go

.PHONY: all build clean run test coverage vet lint docker-build docker-push help
# Default target
all: clean build

# Run the application
run:
	@echo "Running..."
	@${GO} run ${MAIN_PACKAGE}