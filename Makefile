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

# Format the swagger annotations and regenerate the OpenAPI docs
swagger:
	@echo "Update OpenAPI..."
	@go tool swag fmt
	@go tool swag init

.PHONY: all build run test clean watch swagger
