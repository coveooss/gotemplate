
# Assignation

## Global variables

As in Go lang, you must initially declare your global variable using the `:=` assignment operator and subsequent overwrite use the `=` operator.

Sometime, global variables may exist without the code writer being aware of their existence. They could be defined in an imported file and we
need a mean to be able to assign theses variables without pre-assign validation. In that case, you could use the special `~=` assignment operator.

| Razor expression                            | Go Template                                                                     | Note
| ----------------                            | -----------                                                                     | ----
| `@string := "string value";`                | `<pre-assign check code>{{- set $ "string" "string value" }}`                   | Global declare and assign of string
| `@numeric1 := 10;`                          | `<pre-assign check code>{{- set $ "numeric1" 10 }}`                             | Global declare and assign of integer
| `@numeric2 := 1.23;`                        | `<pre-assign check code>{{- set $ "numeric2" 1.23 }}`                           | Global declare and assign of floating point
| `@numeric3 := 4E+4;`                        | `<pre-assign check code>{{- set $ "numeric3" 4E+4 }}`                           | Global declare and assign of large scientific notation number
| `@numeric4 := 5E-3;`                        | `<pre-assign check code>{{- set $ "numeric4" 5E-3 }}`                           | Global declare and assign of small scientific notation number
| `@hexa1 := 0x100;`                          | `<pre-assign check code>{{- set $ "hexa1" 0x100 }}`                             | Global declare and assign of hexadecimal number
| `@result1 := (2+3)*4;`                      | `<pre-assign check code>{{- set $ "result1" (mul (add 2 3) 4) }}`               | Global declare and assign of mathematic expression
| `@result2 := String("hello world!").Title;` | `<pre-assign check code>{{- set $ "result2" ((String "hello world!").Title) }}` | Global declare and assign of generic expression
| `@result2 = "Replaced";`                    | `<pre-assign check code>{{- set $ "result2" "Replaced" }}`                      | Global replacement of previously declared global variable
| `@result1 ~= "Replaced";`                   | `{{- set $ "result1" "Replaced" }}`                                             | Global declare and assign or replacement of previously declared global variable

## Local variables

First declaration of local variable must use the `:=` assignment operator and subsequent assignation must use only `=`.

There are many form used to declare local variable using razor syntax:

- `&#64;{variable} := <value or expression>`
- `&#64;{variable := <value or expression>}`
- `&#64;$variable := <value or expression>`

| Razor expression                             | Go Template                                        | Note
| ----------------                             | -----------                                        | ----
| `@{string := "string value"}`                | `{{- $string := "string value" }}`                 | Local declare and assign of string
| `@{numeric1} := 10`                          | `{{- $numeric1 := 10 }}`                           | Local declare and assign of integer
| `@$numeric2 := 1.23`                         | `{{- $numeric2 := 1.23 }}`                         | Local declare and assign of floating point
| `@{numeric3 := 4E+4}`                        | `{{- $numeric3 := 4E+4 }}`                         | Local declare and assign of large scientific number
| `@{numeric4 := 5E-3}`                        | `{{- $numeric4 := 5E-3 }}`                         | Local declare and assign of small scientific number
| `@{hexa1 := 0x100}`                          | `{{- $hexa1 := 0x100 }}`                           | Local declare and assign of hexadecimal number
| `@{result1 := (2+3)*4}`                      | `{{- $result1 := mul (add 2 3) 4 }}`               | Local declare and assign of mathematic expression
| `@{result2 := String("hello world!").Title}` | `{{- $result2 := (String "hello world!").Title }}` | Local declare and assign of generic expression
| `@{result2} = "Replaced";`                   | `{{- $result2 = "Replaced" }}`                     | Local replacement of previously declared local variable

### Exception

| Razor expression                                | Go Template                                        | Note
| ----------------                                | -----------                                        | ----
| `@{invalid} := print "hello" "world" | upper`   | `{{- $invalid := print }} "hello" "world" | upper` | Using a mixup of go template expression and razor expression could lead to undesired result
| `@{valid := print "hello" "world" | upper}`     | `{{- $valid := print "hello" "world" | upper }}`   | Enclosing the whole assignation statement within {} ensures that the whole expression is assigned
| `@($valid := print "hello" "world" | upper)`    | `{{- $valid := print "hello" "world" | upper }}`   | Using that syntax give the exact same result

### Assignation within expression

```go
@-foreach ($value := to(10))
    @{value}
@-end foreach
```

```go
@-foreach ($index, $value := to(10))
    @{index} = @($value * 2)
@-end foreach
```

```go
@-if ($result := 2+2 == 4)
    result = @{result}
@-end if
```

```go
@-with ($value := 2+2)
    value = @{value}
@-end with
```
