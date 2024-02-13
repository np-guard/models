REPOSITORY := github.com/np-guard/common

mod: go.mod
	@echo -- $@ --
	go mod tidy
	go mod download

fmt:
	@echo -- $@ --
	goimports -local $(REPOSITORY) -w .

lint:
	@echo -- $@ --
	CGO_ENABLED=0 go vet ./...
	golangci-lint run --new

precommit: mod fmt lint

test:
	@echo -- $@ --
	go test ./... -v -cover -coverprofile models.coverprofile
