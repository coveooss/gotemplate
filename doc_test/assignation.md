# Assignation

## Global variables

| Razor expression                            | Go Template                                              | Note
| ----------------                            | -----------                                              | ----
| `@string := "string value";`                | `{{- set $ "string" ("string value") }}`                 | Global assignation of string
| `@numeric1 := 10;`                          | `{{- set $ "numeric1" (10) }}`                           | Global assignation of integer
| `@numeric2 := 1.23;`                        | `{{- set $ "numeric2" (1.23) }}`                         | Global assignation of floating point
| `@numeric3 := 4E+4;`                        | `{{- set $ "numeric3" (4E+4) }}`                         | Global assignation of large scientific notation number
| `@numeric4 := 5E-3;`                        | `{{- set $ "numeric4" (5E-3) }}`                         | Global assignation of small scientific notation number
| `@hexa1 := 0x100;`                          | `{{- set $ "hexa1" (0x100) }}`                           | Global assignation of hexadecimal number
| `@result1 := (2+3)*4;`                      | `{{- set $ "result1" (mul (add 2 3) 4) }}`               | Global assignation of mathematic expression
| `@result2 := string("hello world!").Title;` | `{{- set $ "result2" ((string "hello world!").Title) }}` | Global assignation of generic expression

## Local variables

| Razor expression                            | Go Template                                        | Note
| ----------------                            | -----------                                        | ----
| `$string := "string value";`                | `{{- $string := "string value" }}`                  | Local assignation of string
| `$numeric1 := 10;`                          | `{{- $numeric1 := 10 }}`                            | Local assignation of integer
| `$numeric2 := 1.23;`                        | `{{- $numeric2 := 1.23 }}`                          | Local assignation of floating point
| `$numeric3 := 4E+4;`                        | `{{- $numeric3 := 4E+4 }}`                          | Local assignation of large scientific number
| `$numeric4 := 5E-3;`                        | `{{- $numeric4 := 5E-3 }}`                          | Local assignation of small scientific number
| `$hexa1 := 0x100;`                          | `{{- $hexa1 := 0x100 }}`                            | Local assignation of hexadecimal number
| `$result1 := (2+3)*4;`                      | `{{- $result1 := mul (add 2 3) 4 }}`                | Local assignation of mathematic expression
| `$result2 := string("hello world!").Title;` | `{{- $result2 := (string "hello world!").Title }}`  | Local assignation of generic expression

### Exception

| Razor expression                              | Go Template                                        | Note
| ----------------                              | -----------                                        | ----
| `$invalid := print "hello" "world" | upper;`  | `{{- $invalid := print }} "hello" "world" | upper`  | Using a mixup of go template expression and razor expression could lead to undesired result
| `@$valid := print "hello" "world" | upper;`   | `{{- $valid := print "hello" "world" | upper }}`    | Adding a @ before assign expression solve the evaluation problem

### Assignation within expression

```go
@range ($value := to(10))
    @$value
@end range
```

```go
@range ($index, $value := to(10))
    @$index = @($value * 2)
@end range
```

```go
@if ($result := 2+2 == 4)
    result = @$result
@end if
```

```go
@with ($value := 2+2)
    value = @$value
@end with
```