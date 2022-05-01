GOCMD=go
GOTEST=$(GOCMD) test

test:
	@echo "\n\n==================== Start unit test and Integration Test ...... ====================\n\n"
	$(GOTEST) ./... -cover -race
	@echo "\n\n==================== Unit test and Integration Test Done ====================\n\n"

unittest:
	@echo "\n\n==================== Start unit test ...... ====================\n\n"
	@go test ./... --short -cover -race
	@echo "\n\n==================== Unit test done ====================\n\n"

lint:
	@golangci-lint run

run:
	@go run cmd/main/main.go
build:
	@go build -o kumparan-api cmd/main/main.go

docker-run:
	@docker stop kumparan_container && \
	docker rm -f kumparan_container && \
	docker image rm -f kumparan_api && \
	docker build -t kumparan_api . && \
	docker run --name=kumparan_container -d -it -p 9000:9000 kumparan_api

mock:
	@mockery --dir=./src/business --name=ScholarshipService --output=./mocks
	@mockery --dir=./src/business --name=ScholarshipRepository --output=./mocks

compose-rebuild:
	@docker-compose up --build -d

compose-up:
	@docker-compose up -d

compose-down:
	@docker-compose down && docker image rm kumparan-image