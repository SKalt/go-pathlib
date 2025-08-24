.PHONY: test coverage docs
test:
	mise run test

coverage: test
	go tool cover -html=.coverage

bin/pkgsite:
	go install golang.org/x/pkgsite/cmd/pkgsite@latest

docs: bin/pkgsite
	pkgsite .

lint:
	mise run lint