# OS commands

It is possible to run OS commands using the following go template functions:

* `exec` returns the result of a shell command as structured data.
* `run` returns the result of a shell command as a string.

## exec

### Razor (exec)

```go
@{example} := exec("printf 'SomeData: test2\nSomeData2: test3'")
First result: @{example.SomeData}
Second result: @{example.SomeData2}
@{example}

@{example2} := exec("printf 'Test'")
Should be `string`: @typeOf($example2)
@{example2}
```

### Gotemplate (exec)

```go
{{- $example := exec "printf 'SomeData: test2\nSomeData2: test3'" }}
First result: {{ $example.SomeData }}
Second result: {{ $example.SomeData2 }}
{{ $example }}

{{- $example2 := exec "printf 'Test'" }}
Should be `string`: {{ typeOf $example2 }}
{{ $example2 }}
```

### Result (exec)

```text
First result: test2
Second result: test3
SomeData: test2
SomeData2: test3

Should be `string`: string
Test
```

## run

### Razor (run)

```go
@{example} := run("printf 'SomeData: test2\nSomeData2: test3'")
Should be `string`: @typeOf($example)
@{example}
```

### Gotemplate (run)

```go
{{- $example := run "printf 'SomeData: test2\nSomeData2: test3'" }}
Should be `string`: {{ typeOf $example }}
{{ $example }}
```

### Result (run)

```text
Should be `string`: string
SomeData: test2
SomeData2: test3
```
