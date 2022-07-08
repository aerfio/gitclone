.PHONY: test
test:
	go test -race -v ./...

.PHONY: install
install:
	go install .