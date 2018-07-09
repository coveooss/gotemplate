{% include navigation.html %}
{% raw %}

# OS commands

It is possible to run OS commands using the following go template functions:

* `exec` returns the result of a shell command as structured data.
* `run` returns the result of a shell command as a string.

## exec

### Razor
```
@$hello := exec("printf 'SomeData: test2\nSomeData2: test3'")
First result: @$hello.SomeData
Second result: @$hello.SomeData2
@$hello
```

### Gotemplate
```
{{- $hello := exec "printf 'SomeData: test2\nSomeData2: test3'" }}
First result: {{ $hello.SomeData }}
Second result: {{ $hello.SomeData2 }}
{{ $hello }}
```

### Result
```
First result: test2
Second result: test3
SomeData: test2
SomeData2: test3
```

## run

### Razor
```
@$hello := run("printf 'SomeData: test2\nSomeData2: test3'")
Should be `string`: @typeOf($hello)
@$hello
```

### Gotemplate
```
{{- $hello := run "printf 'SomeData: test2\nSomeData2: test3'" }}
Should be `string`: {{ typeOf $hello }}
{{ $hello }}
```

### Result
```
Should be `string`: string
SomeData: test2
SomeData2: test3
```


{% endraw %}
