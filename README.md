[![GitHub Actions](https://shields.io/github/actions/workflow/status/hanggrian/propertieslint/code-analysis.yaml)](https://github.com/hanggrian/propertieslint/actions/workflows/code-analysis.yaml)
[![Codecov](https://shields.io/codecov/c/gh/hanggrian/propertieslint)](https://app.codecov.io/gh/hanggrian/propertieslint/)
[![Renovate](https://shields.io/badge/renovate-enabled-brightgreen)](https://developer.mend.io/github/hanggrian/propertieslint/)
[![GitHub Release](https://shields.io/github/release/hanggrian/propertieslint)](https://pkg.go.dev/github.com/hanggrian/propertieslint/)
[![Go](https://shields.io/github/go-mod/go-version/hanggrian/propertieslint)](https://go.dev/doc/go1.25)

# propertieslint

Small Go CLI for linting Java [Properties](https://docs.oracle.com/javase/tutorial/essential/environment/properties.html)
files. It catches common issues like duplicate keys, missing value and
dangling line continuations. Rules are configured with local
`.propertieslint.json` file or with `-c`/`--config` flag.

[View rules](https://hanggrian.github.io/propertieslint/rules/comment/comment-spacing/)

## Download

```sh
go install github.com/hanggrian/propertieslint@latest
```

## Usage

The CLI walks directories recursively and lints every `.properties` file it
finds.

```sh
propertieslint some_file.properties ./some_dir
```

> [!NOTE]
>
> All rules are enabled by default.

For example, to disable certain rules, refer to the rule ID within group name
and set it to `false`. To disable all rules in a group, set the group name to
`false`.

```json
{
  "comment": {
    "comment-spacing": false
  },
  "format": false
}
```
