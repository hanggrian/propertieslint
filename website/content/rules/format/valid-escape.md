---
weight: 20
---

# Valid escape

Entries must have valid escape sequences. For example, unicode escapes must
include four hexadecimal digits (`\u1234`, not `\u123`).

**Before &#10060;**

```
unicode=\u123
```

**After &#9989;**

```
unicode=\u1234
```
