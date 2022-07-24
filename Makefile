GOLANGCI_VERSION = v1.47.0

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_VERSION)
	@mv bin/golangci-lint $@

ci: generate lint test ##=> Run all CI targets
.PHONY: ci

generate: ##=> generate all the things
	@echo "--- generate all the things"
	@go generate ./...
.PHONY: generate

.PHONY: lint
lint: bin/golangci-lint ##=> Lint all the things
	@echo "--- lint all the things"
	@bin/golangci-lint run

.PHONY: clean
clean: ##=> Clean all the things
	$(info [+] Clean all the things...")

.PHONY: test
test: ##=> Run the tests
	$(info [+] Run tests...")
	@go test -v -cover ./...
