OUTPUT = main 
SERVICE_NAME = location-history-api

.PHONY: test
test:
	go test ./...

build-local:
	go build -o $(OUTPUT) ./cmd/$(SERVICE_NAME)/main.go

test:
	go test ./...

run: build-local
	@echo ">> Running application ..."
	HISTORY_SERVER_LISTEN_ADDR= \
	LOCATION_HISTORY_TTL_SECONDS=5000 \
	./$(OUTPUT)