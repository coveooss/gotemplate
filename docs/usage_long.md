{% include navigation.html %}

{% raw %}
## Data structures
```go
// JSON data definition
{{ $json_data := `"number": 10, "text": "Hello World!, "data": { "a": 1 }, "list": ["value 1, "value 2"]` | json }}

// YAML data definition
{{ $yaml_data := yaml `
number: 10
text: Hello world!
data:
    a: 1
    b: 2
list:
    - value 1
    - value 2
` }}

// HCL data definition from template
{{ define "source" }}
number = 10
text = "Hello world!"

data {
    a = 1
    b = 2
}

list = [
    "value 1",
    "value 2",
]
{{ end }}

{{ $hcl_data := hcl "source" }}

// YAML data definition for file
{{ $yaml_file := yaml "source.yml" }}

// Assign arbitrary data definition to a variable
// data can handle any HCL, JSON, YAML or string coming from the string directly or a template name or a file name
{{ $foo := `a = 1 b = 2 c = "Hello world!"` | data }}
```

## Examples
```go
// Random text
{{ $text := lorem "paragraph" 5 10 }}

// Random email
{{ $email := lorem "email" }}

// Concatenation of various objects
{{ $string := concat "Hello" 1 " pi = " 3.14 " 3 x 5 = " (mul 3 5)}}

// Boolean conversion
{{ $bool := bool "True" }}

// Quote all elements of a list
{{ $quoted := list 1 2 3 | formatList "\"%v\"" }}

// Colored strings
{{ $colored := color "red" "I am red" }}
{{ $colored := color "red,bgwhite" "I am red with white background" }}
{{ $colored := color "red" "bgyellow" "I am %s with %s background" "red" "yellow" }}

// Define a script
{{ define "script" }}
#! /usr/bin/env python

for i in range({{ . | default 10 }}):
    print("- ".format(i))
{{ end }}

// Run the script (render the data as output by the python script)
{{ run "script" }}

// Run the script (returns a list of number from 0 to 24)
{{ $values := exec "script" 25 }}

// Get the content of a folder
{{ $values := run "ls" "-l" }}

// Get the content of the root folder (note that arguments could be either separated of embedded with the command)
{{ $values := run "ls -l /" }}

// Get the script in a variable and execute it later
{{ $script := include "script" 40 }}
{{ exec $script }}

// Add a function that could be called later (should be created in a .template file and used later as a function)
// Note that it is not possible to
{{ alias "generate_numbers" "exec" "script" }}
...
// Can be later used as a function in another template script
{{ $numbers := generate_numbers 200 }}
{{ $numbers := 1000 | generate_numbers }}
```
{% endraw %}
