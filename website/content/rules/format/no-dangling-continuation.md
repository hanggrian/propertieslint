---
weight: 20
---

# No dangling continuation

Files cannot end with a dangling line continuation. A backslash (`\`) at the
end of a line continues the entry onto the next line, so it must be followed by
one.

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
