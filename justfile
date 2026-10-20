install CLEAN="false":
    make install CLEAN={{ CLEAN }}
    uv sync {{ if CLEAN == "true" { "--locked" } else { "" } }}
    pnpm install {{ if CLEAN == "true" { "--frozen-lockfile" } else { "" } }}

[private]
lint1:
    make lint

[private]
lint2:
    uv run poe lint

[private]
lint3:
    pnpm lint

[parallel]
lint: lint1 lint2 lint3

format:
    just --fmt
    make format

doc:
    make doc

# generate the website with embedded Go API docs
generate-website: doc
    rm -rf website/static/api/
    mkdir -p website/static/api/
    mv build/doc2go/ website/static/api/godoc/
    pnpm exec hugo --source website --minify
