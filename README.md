[![GitHub Actions](https://shields.io/github/actions/workflow/status/hanggrian/propertieslint/code-analysis.yaml)](https://github.com/hanggrian/propertieslint/actions/workflows/code-analysis.yaml)
[![Codecov](https://shields.io/codecov/c/gh/hanggrian/propertieslint)](https://app.codecov.io/gh/hanggrian/propertieslint/)
[![Renovate](https://shields.io/badge/renovate-enabled-brightgreen)](https://developer.mend.io/github/hanggrian/propertieslint/)
[![Go](https://img.shields.io/github/go-mod/go-version/hanggrian/propertieslint)](https://go.dev/doc/go1.25)

# propertieslint

Small Go CLI for linting Java-style `.properties` files.

## Download

```sh
go install github.com/hanggrian/propertieslint@latest
```

## Usage

The CLI walks directories recursively and lints every `.properties` file it
finds.

```sh
go run ./cmd/propertieslint ./path/to/dir
```

It accepts `-c` or `--config` and reads a JSON file with rule toggles. If no
config path is supplied, it looks for `.propertieslint.json` in the current
directory.

Example:

```json
{
  "duplicate-key": false
}
```

Supported keys:

- `duplicate-blank-line`: Multiple consecutive blank lines.
- `duplicate-key`: Entries with duplicate keys.
- `invalid-escape`: Entries with invalid escape sequences.
- `missing-separator`: Entries without `=` or `:` separators.
- `missing-value`: Entries with a key but no value.
- `no-leading-blank-line`: No blank line at the beginning of the file.
- `trailing-newline`: No newline at the end of the file.
- `unterminated-line-continuation`: Entries with unterminated line
  continuations.
- `untrimmed-entry`: Entries with leading or trailing whitespace.
