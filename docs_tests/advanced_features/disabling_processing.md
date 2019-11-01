# Disabling processing

You can completely disable gotemplate in a file with a comment. There's a comment that will completely disable gotemplate and one that will only disable razor, here's an example:

```go
# no-gotemplate!
{{ 4 + 5 }} # This will not be interpreted
@(4 + 5) # This will also not be interpreted

# no-razor!
@(4 + 5) # This will not be interpreted
```

Note that if the comment is found *anywhere* in the file, it will apply
