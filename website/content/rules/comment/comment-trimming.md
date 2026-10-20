---
weight: 10
---

# Comment trimming

Empty comments at the beginning or end of a comment block with content must be
removed. Empty comments between non-empty comments are allowed as separators.
A block containing only empty comments is handled by
[no-empty-comment](no-empty-comment.md).

**Before &#10060;**

```
#
# one
#
```

**After &#9989;**

```
# one
```
