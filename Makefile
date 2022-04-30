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
	@go build -o scholarship-api cmd/app/main.go

docker-run:
	@docker stop bangun_container && \
	docker rm -f bangun_container && \
	docker image rm -f bangun_api && \
	docker build -t bangun_api . && \
	docker run --name=bangun_container -d -it -p 7070:7070 bangun_api

mock:
	@mockery --dir=./src/business --name=ScholarshipService --output=./mocks
	@mockery --dir=./src/business --name=ScholarshipRepository --output=./mocks