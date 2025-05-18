.PHONY: test test-integration test-unit test-handlers setup-test-db run

# Run the application
run:
	go run cmd/api/main.go --config config.yaml --file data/users_data.json

# Start test database
setup-test-db:
	docker-compose -f test/docker-compose.test.yaml up -d
	sleep 5  # Wait for database to be ready

# Run all tests
test: test-unit test-integration test-handlers

# Run unit tests
test-unit:
	go test -v -cover ./service/...

# Run integration tests
test-integration: setup-test-db
	go test -v -cover ./test/integration/...

# Run handler tests
test-handlers:
	go test -v -cover ./api/http/handlers/...

# Clean up test database
clean-test-db:
	docker-compose -f test/docker-compose.test.yaml down -v 