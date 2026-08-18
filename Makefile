.PHONY: swagger deps build run clean

# Generate Swagger UI documentation
swagger:
	swag init -g cmd/monitor-agent/main.go --parseDependency --parseInternal

# Update Go dependencies and clean up go.mod/go.sum
deps: swagger
	go mod tidy

# Build the executable binary (refreshes swagger docs beforehand)
build: swagger
	go build -o bin/monitor-agent.exe ./cmd/monitor-agent/main.go

# Run the agent in development mode
run: swagger
	go run ./cmd/monitor-agent/main.go

# Clean generated build artifacts and documentation
clean:
	rm -rf docs/
	rm -rf bin/