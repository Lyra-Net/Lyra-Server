API_GATEWAY_DIR=api-gateway
IDENTITY_SERVICE_DIR=identity-service
SONG_SERVICE_DIR=song-service
STREAM_SERVICE_DIR=stream-service

GO_RUN=go run main.go
NEST_START=pnpm start

.PHONY: all start start-api start-identity start-song start-stream clean build


start:
	@echo "🚀 Starting all backend services..."
	concurrently \
		"cd $(API_GATEWAY_DIR) && $(NEST_START)" \
		"cd $(IDENTITY_SERVICE_DIR) && $(GO_RUN)" \
		"cd $(SONG_SERVICE_DIR) && $(GO_RUN)" \
		"cd $(STREAM_SERVICE_DIR) && $(NEST_START)"

start-api:
	cd $(API_GATEWAY_DIR) && $(NEST_START)

start-identity:
	cd $(IDENTITY_SERVICE_DIR) && $(GO_RUN)

start-song:
	cd $(SONG_SERVICE_DIR) && $(GO_RUN)

start-stream:
	cd $(STREAM_SERVICE_DIR) && $(NEST_START)

build:
	pnpm --filter ./$(API_GATEWAY_DIR) build
	pnpm --filter ./$(STREAM_SERVICE_DIR) build

clean:
	rm -rf $(API_GATEWAY_DIR)/dist $(API_GATEWAY_DIR)/node_modules
	rm -rf $(STREAM_SERVICE_DIR)/dist $(STREAM_SERVICE_DIR)/node_modules
