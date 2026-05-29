---
weight: 20
---

# Invalid escape

Entries with invalid escape sequences. This includes short unicode escapes (e.g.
`\u123`), invalid octal escapes (e.g. `\8`), and invalid single-character
escapes (e.g. `\x`).

**Before &#10060;**

```
unicode=\u123
```

**After &#9989;**

```
unicode=\u1234
```
