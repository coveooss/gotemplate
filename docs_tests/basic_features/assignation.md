
# Assignation

## Global variables

As in Go lang, you must initially declare your global variable using the `:=` assignment operator and subsequent overwrite use the `=` operator.

| Razor expression                            | Go Template                                                                     | Note
| ----------------                            | -----------                                                                     | ----
| `@string := "string value";`                | `{{- set $ "string" "string value" }}`                   | Global declare and assign of string
| `@numeric1 := 10;`                          | `{{- set $ "numeric1" 10 }}`                             | Global declare and assign of integer
| `@numeric2 := 1.23;`                        | `{{- set $ "numeric2" 1.23 }}`                           | Global declare and assign of floating point
| `@numeric3 := 4E+4;`                        | `{{- set $ "numeric3" 4E+4 }}`                           | Global declare and assign of large scientific notation number
| `@numeric4 := 5E-3;`                        | `{{- set $ "numeric4" 5E-3 }}`                           | Global declare and assign of small scientific notation number
| `@hexa1 := 0x100;`                          | `{{- set $ "hexa1" 0x100 }}`                             | Global declare and assign of hexadecimal number
| `@result1 := (2+3)*4;`                      | `{{- set $ "result1" (mul (add 2 3) 4) }}`               | Global declare and assign of mathematic expression
| `@result2 := String("hello world!").Title;` | `{{- set $ "result2" ((String "hello world!").Title) }}` | Global declare and assign of generic expression
| `@result2 = "Replaced";`                    | `{{- set $ "result2" "Replaced" }}`                      | Global replacement of previously declared global variable

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

## Assignment operators

Using the Razor syntax, it is possible to use assignment operators such as `+=`, `/=`... The operators that are supported are:

| Operator    | Assignment   | Note
| ----------- | ------------ | ----
| `+`         | `+=`         | Addition
| `-`         | `-=`         | Subtraction
| `*`         | `*=`         | Multiplication
| `/`, `÷`    | `/=`, `÷=`   | Division
| `%`         | `%=`         | Modulo
| `&`         | `&=`         | Bitwise AND
| `|`         | `|=`         | Bitwise OR
| `^`         | `^=`         | Bitwise XOR
| `&^`        | `&^=`        | Bit clear
| `<<`, `«`   | `<<=`, `«==` | Left shift
| `>>`, `»`   | `>>=`, `»==` | Right shift

| Razor expression  | Go Template                         | Note
| ----------------  | -----------                         | ----
| `@num := 10`      | `{{- set $ "num" 10 }}`             | Add assignment operator on global
| `@num += 5`       | `{{- set $ "num" (add $.num 10) }}` | Add assignment operator on global
| `@{local} := 5`   | `{{- $local := 5 }}`                | Local assignation
| `@$local += 10`   | `{{- $local = add $local 10 }}`     | Add assignment operator on local
| `@{local} *= 20`  | `{{- $local = mul $local 20 }}`     | Multiply assignment operator on local
| `@{local /= 2}`   | `{{- $local = div $local 2 }}`      | Divide assignment operator on local

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
