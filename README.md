# propertieslint

`propertieslint` is a small Go CLI for linting Java-style `.properties` files.

## What it checks

- Missing key/value separators
- Empty values after an explicit separator
- Duplicate keys
- Invalid escape sequences
- Unterminated line continuations

## Configuration

The CLI accepts `-c` or `--config` and reads a JSON file with rule toggles.
If no config path is supplied, it looks for `propertieslint.json` in the current directory.

Example:

```json
{
  "duplicate-key": false
}
```

Supported keys:

- `missing-key-value-separator`
- `missing-value`
- `duplicate-key`
- `invalid-escape`
- `unterminated-line-continuation`

## Usage

The CLI walks directories recursively and lints every `.properties` file it
finds.

```sh
go run ./cmd/propertieslint ./path/to/dir
```
