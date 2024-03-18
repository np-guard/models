REPOSITORY := github.com/np-guard/models
JSON_PACKAGE_NAME := spec

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
	golangci-lint run

precommit: mod fmt lint

test:
	@echo -- $@ --
	go test ./... -v -cover -coverprofile models.coverprofile

pkg/${JSON_PACKAGE_NAME}/data_model.go: spec_schema.json
	@echo -- generate --
	# Install https://github.com/atombender/go-jsonschema
	go-jsonschema spec_schema.json --package ${JSON_PACKAGE_NAME} --struct-name-from-title --tags json --output $@
	goimports -local $(REPOSITORY) -w $@

generate: pkg/${JSON_PACKAGE_NAME}/data_model.go
