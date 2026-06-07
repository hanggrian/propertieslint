.PHONY: install lint test cov format doc


install:
	@if [ "$(CLEAN)" = "true" ]; then \
		go mod download; \
	else \
		go mod tidy; \
	fi

lint:
	go run . .

test:
	go test ./linter/...

cov:
	go test -coverprofile=coverage.out ./linter/...

format:
	go fmt ./...

doc:
	rm -rf build/doc2go/
	doc2go -out build/doc2go/ ./linter/...
