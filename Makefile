WIREMOCK_BASE_URL=http://localhost:8989

build:
	go build -o ./bin/webcrawler ./cmd/main.go

run_local: docker_start build
	./bin/webcrawler -url=$(WIREMOCK_BASE_URL) -depth=5 -delay=10

docker_start:
	docker compose up -d

integration_test: docker_start
	curl -X POST $(WIREMOCK_BASE_URL)/__admin/reset
	curl -X POST --data "@$(PWD)/integration/wiremock/mappings/full_page.json" $(WIREMOCK_BASE_URL)/__admin/mappings/import
	go test -p 1 -count 1 -coverpkg=./... -coverprofile=coverage.out ./...

coverage_html:
	go tool cover -html=coverage.out
