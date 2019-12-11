# Data collections (lists/slices and dicts/maps)

## Maps

| Razor                                           | Gotemplate                                          | Note
| ---                                             | ---                                                 | ---
| `@razorDict := dict("test", 1, "test2", 2);`    | `{{- set $ "goDict" (dict "test" 1 "test2" 2) }}`   | Creation
| `@razorDict2 := dict("test3", 3, "test5", 5);`  | `{{- set $ "goDict2" (dict "test3" 3 "test5" 5) }}` | Creation
| `@set(.razorDict, "test2", 3);`                 | `{{- set .goDict "test2" 3 }}`                      | Update
| `@set(.razorDict, "test3", 4);`                 | `{{- set .goDict "test3" 4 }}`                      | Update
| `@razorDict = merge(razorDict, razorDict2);`    | `{{- set $ "goDict" (merge .goDict .goDict2) }}`    | Merge (First dict has priority)
| `@razorDict.test3;`                             | `{{ $.goDict.test3 }}`                              | Should be `4`
| `@razorDict.test5;`                             | `{{ $.goDict.test5 }}`                              | Should be `5`
| `@razorDict["test", "test3", "undef"]`          | `{{ extract $.goDict "test" "test3" "undef" }}`     | Extract values, should be `[1,4,null]`
| `@razorDict["test":"test9"]`                    | `{{ slice $.goDict "test" "test9" }}`               | Slice values, should be `[1,3,4,5]`
| `@razorDict["test9":"test"]`                    | `{{ slice $.goDict "test9" "test" }}`               | Slice values, should be `[5,4,3,1]`
| `@keys(razorDict)`                              | `{{ keys $.goDict }}`                               | Get keys, should be `["test","test2","test3","test5"]`
| `@values(razorDict)`                            | `{{ values $.goDict }}`                             | Get values, should be `[1,3,4,5]`

### Looping (Maps)

#### Razor (Maps)

```go
@-foreach($key, $value := razorDict)
    @{key}, @{value}, @get($.razorDict, $key)
@-end foreach
```

#### Gotemplate (Maps)

```go
{{- range $key, $value := .goDict }}
    {{ $key }}, {{ $value }}, {{ get $.goDict $key }}
{{- end }}
```

#### Result (Maps)

```go
    test, 1, 1
    test2, 3, 3
    test3, 4, 4
    test5, 5, 5
```

## Slices

| Razor                                            | Gotemplate                                             | Note
| ---                                              | ---                                                    | ---
| `@razorList := list("test1", "test2", "test3");` | `{{- set $ "goList" (list "test1" "test2" "test3") }}` | Creation
| `@razorList = append(razorList, "test4");`       | `{{- set $ "goList" (append .goList "test4") }}`       | Append
| `@razorList = prepend(razorList, "test0");`      | `{{- set $ "goList" (prepend .goList "test0") }}`      | Prepend
| `@razorList;`                                    | `{{ $.goList }}`                                       | Should be `["test0","test1","test2","test3","test4"]`
| `@contains(razorList, "test1", "test2");`        | `{{- contains .goList "test1" "test2" }}`              | Check if element is in list
| `@has("test1", razorList);`                      | `{{- has "test1" .goList }}`                           | has is an alias to contains
| `@has(razorList, "test1");`                      | `{{- has .goList "test1" }}`                           | has Support inversion of argument if the first one is not a list
| `@has(razorList, "test1", "test2");`             | `{{- has .goList "test1" "test2" }}`                   | has can also test for many elements
| `@razorList.Contains("test1", "test2");`         | `{{ .goList.Contains "test1" "test2" }}`               | List also support using methods
| `@razorList.Reverse();`                          | `{{ .goList.Reverse }}`                                | Should be `["test4","test3","test2","test1","test0"]`

### Looping (Slice)

#### Razor (Slice)

```go
@-foreach($index, $value := razorList)
    @{index}, @{value}, @extract($.razorList, $index)
@-end foreach
```

#### Gotemplate (Slice)

```go
{{- range $index, $value := .goList }}
    {{ $index }}, {{ $value }}, {{ extract $.goList $index }}
{{- end }}
```

#### Result (Slice)

```go
    0, test0, test0
    1, test1, test1
    2, test2, test2
    3, test3, test3
    4, test4, test4
```
