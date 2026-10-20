---
weight: 10
---

# No empty comment

Comment blocks must have content. A block containing only empty comments is
invalid. If a block contains at least one non-empty comment, leading or trailing
empty comments are handled by [comment-trimming](comment-trimming.md), and
interior empty comments are allowed as separators.

**Before &#10060;**

```
#
#
```

**After &#9989;**

```
```
