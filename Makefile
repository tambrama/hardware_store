MAIN=cmd/app/main.go
BUILD_DIR =build
O_FILE=server

.PHONY: all build run clear

all: docker-up

build: 
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(O_FILE) $(MAIN)

run: build
	./$(BUILD_DIR)/$(O_FILE)

clean:
	rm -rf $(BUILD_DIR)

fmt:
	go fmt ./...

update_mod:
	go mod tidy

docker-up:
	docker-compose -f docker-compose.yml up --build -d

docker-down:
	docker-compose -f docker-compose.yml down

docker-rebuild:
	docker-compose -f docker-compose.yml down -v && \
	docker-compose -f docker-compose.yml up --build -d

docker-logs:
	docker-compose -f docker-compose.yml logs -f app

docs:
	swag init --parseDependency --parseInternal --generalInfo cmd/app/main.go --output ./docs