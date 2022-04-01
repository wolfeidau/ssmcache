ci: generate lint test ##=> Run all CI targets
.PHONY: ci

generate: ##=> generate all the things
	@echo "--- generate all the things"
	@go generate ./...
.PHONY: generate

lint: ##=> Lint all the things
	@echo "--- lint all the things"
	@docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.45.2 golangci-lint run -v
.PHONY: lint

clean: ##=> Clean all the things
	$(info [+] Clean all the things...")
.PHONY: clean

test: ##=> Run the tests
	$(info [+] Run tests...")
	@go test -v -cover ./...
.PHONY: test
