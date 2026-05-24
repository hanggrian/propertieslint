install CLEAN="false":
    go mod {{ if CLEAN == "true" { "download" } else { "tidy" } }}
    uv sync {{ if CLEAN == "true" { "--locked" } else { "" } }}
    pnpm install {{ if CLEAN == "true" { "--frozen-lockfile" } else { "" } }}

[group: 'check']
lint-go:
    go run . .

[group: 'check']
lint-python:
    uv run poe lint

[group: 'check']
lint-node:
    pnpm lint

[group: 'check']
[parallel]
lint: lint-go lint-python lint-node

[group: 'check']
test:
    go test ./linter/...

[group: 'check']
cov:
    go test -coverprofile=coverage.out ./linter/...

format:
    just --fmt
    go fmt ./...

doc:
    rm -rf build/doc2go/
    doc2go -out build/doc2go/ ./linter/...
