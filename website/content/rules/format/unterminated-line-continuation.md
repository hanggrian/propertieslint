---
weight: 20
---

# Unterminated line continuation

Entries with unterminated line continuations. Line continuations must end with a
backslash (`\`) followed by a newline.

**Before &#10060;**

```
multiline=line1\

escape=this is backslash \
```

**After &#9989;**

```
multiline=line1\
line2

escape=this is backslash \\
```
