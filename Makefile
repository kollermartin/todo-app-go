.PHONY: fmt

# Format all Go files in unix-like environment
fmt:
	@echo "Formatting Go files"
	@find . -name "*.go" -type f -exec go fmt {} \;

# Format all Go files in windows environment
fmt-win:
	@echo "Formatting Go files"
	@powershell -Command "Get-ChildItem -Recurse -Filter '*.go' | ForEach-Object { go fmt $$_.FullName }"

lint:
	golangci-lint run

build:
	go build -o bin/app ./cmd/http

run:
	./bin/app

test:
	go test -v ./... -count=1