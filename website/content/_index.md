---
layout: landing
---

<br/>

# propertieslint {anchor=false}

Small Go CLI for linting Java [Properties](https://docs.oracle.com/javase/tutorial/essential/environment/properties.html)
files. It catches common issues like duplicate keys, missing value and
unterminated line continuations. Rules are configured with local
`.propertieslint.json` file or with `-c`/`--config` flag.

{{<button href="/rules/comment/comment-spaces/">}}Explore{{</button>}}

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
    "comment-spaces": false
  },
  "format": false
}
```
