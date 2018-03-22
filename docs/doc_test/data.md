{% include navigation.html %}
{% raw %}
## Data manipulation

Using a data file with the following content in a format that doesn't follow a standard.
```!Data
IntegerValue = 1
FloatValue = 1.23
StringValue = "Foo bar"
EquationResult = @(2 + 2 * 3 ** 6)
ListValue = ["value1", "value2"]
DictValue = {"key1": "value1", "key2": "value2"}
```

### toYaml

| Razor | Gotemplate
| ---   | ---
| ```@toYaml(data("!Data"))``` | ```{{ toYaml (data "!Data") }}```

```
DictValue:
  key1: value1
  key2: value2
EquationResult: 46658
FloatValue: 1.23
IntegerValue: 1
ListValue:
- value1
- value2
StringValue: Foo bar
```

### toJson

| Razor | Gotemplate
| ---   | ---
| ```@toPrettyJson(data("!Data"))``` | ```{{ toPrettyJson (data "!Data") }}```

```
{
  "DictValue": {
    "key1": "value1",
    "key2": "value2"
  },
  "EquationResult": 46658,
  "FloatValue": 1.23,
  "IntegerValue": 1,
  "ListValue": [
    "value1",
    "value2"
  ],
  "StringValue": "Foo bar"
}
```

### toHcl

| Razor | Gotemplate
| ---   | ---
| ```@toPrettyHcl(data("!Data"))``` | ```{{ toPrettyHcl (data "!Data") }}```

```
EquationResult = 46658
FloatValue     = 1.23
IntegerValue   = 1
ListValue      = ["value1", "value2"]
StringValue    = "Foo bar"

DictValue {
  key1 = "value1"
  key2 = "value2"
}
```

### Nested conversions

This test shows how you can convert from and to other formats.

| Razor | Gotemplate
| ---   | ---
| ```@toPrettyTFVars(data(toTFVars(fromHcl(toHcl(fromJson(toJson(data("!Data"))))))))``` | ```{{ toPrettyTFVars (data (toTFVars (fromHcl (toHcl (fromJson (toJson (data "!Data"))))))) }}```

```
EquationResult = 46658
FloatValue     = 1.23
IntegerValue   = 1
ListValue      = ["value1", "value2"]
StringValue    = "Foo bar"

DictValue {
  key1 = "value1"
  key2 = "value2"
}
```
{% endraw %}