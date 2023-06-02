fmt:
	gofmt -s -w .

test:
	go test ./... -cover

mock:
	mockery --all

coverage_cli:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage_html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

