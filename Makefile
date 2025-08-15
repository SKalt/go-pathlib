.PHONY: test coverage docs
test:
	go test -coverprofile=.coverage ./...

coverage: test
	go tool cover -html=.coverage

docs:
	 go tool golang.org/x/pkgsite/cmd/pkgsite .

lint:
	golangci-lint run