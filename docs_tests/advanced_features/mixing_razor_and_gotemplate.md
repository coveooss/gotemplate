# Operators mixin

Note that you cannot combine razor extended expression (+, -,  /, *, etc.) with go template expression such as in:

## Razor syntax

```go
# In this statement, | is interpreted as bitwise or between 2 and 4
@(2 + (2 | mul(4)))

# While in this statement (no binary operator), | is interpreted as go template piping operator
@(sum 2 (2 | mul 4))
```

## gotemplate syntax

```go
{{ add 2 (bor 2 (mul 4)) }}
{{ sum 2 (2 | mul 4) }}
```

### Result

```go
8
10
```
