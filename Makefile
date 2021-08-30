.DEFAULT_GOAL := test

.PHONY: test
test: test-unit test-integration test-docker

.PHONY: test-unit
test-unit:
	@go test -count=1 -v -race -cover ./...

.PHONY: docker-test
test-docker:
	@./tests/docker/test.sh

.PHONY: test-integration
test-integration:
	$(MAKE) -C ./instrumentation/opentelemetry/github.com/jackc/hyperpgx/integrationtest test

.PHONY: bench
bench:
	go test -v -run - -bench . -benchmem ./...

.PHONY: lint
lint:
	@echo "Running linters..."
	@golangci-lint run ./... && echo "Done."

.PHONY: deps
deps:
	@go get -v -t -d ./...

.PHONY: ci-deps
deps-ci:
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

check-examples:
	find ./instrumentation -type d -print | \
	grep examples/ | \
	xargs -I {} bash -c 'if [ -f "{}/main.go" ] ; then cd {}; go build -o ./build_example main.go ; fi'
	find . -name "build_example" -delete

.PHONY: fmt
fmt:
	gofmt -w -s ./

.PHONY: tidy
tidy:
	find . -path ./config -prune -o -name "go.mod" \
	| grep go.mod \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'cd {}; go mod tidy'

.PHONY: install-tools
install-tools: ## Install all the dependencies under the tools module
	$(MAKE) -C ./tools install

.PHONY: check-vanity-import
check-vanity-import:
	@porto -l .
	@if [[ "$(porto --skip-files ".*\\.pb\\.go$" -l . | wc -c | xargs)" -ne "0" ]]; then echo "Vanity imports are not up to date" ; exit 1 ; fi