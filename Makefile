all: build

# Build the application
build:
	@echo "Building..."
	@go build main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Update swagger
swagger:
	@echo "Update OpenAPI..."
	@swag init

.PHONY: all build run test clean watch
