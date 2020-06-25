# Functions

## Defining functions

The following file is evaluated: [!extensions.gte](!extensions.gte)

## Variable set in extension are globally available

The value of `@GlobalExtensionVariable` should be `I am set and I am global`.

## Calling functions

### Without arguments

| Razor | Gotemplate
| ---   | ---
| ```@ChristmasTree()``` | ```{{ ChristmasTree }}```

#### Result

```text
                       ✾
                      ✾✾✾
                     ✾✾✾✾✾
                    ✾✾✾✾✾✾✾
                   ✾✾✾✾✾✾✾✾✾
                  ✾✾✾✾✾✾✾✾✾✾✾
                 ✾✾✾✾✾✾✾✾✾✾✾✾✾
                ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
               ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
              ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
             ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
            ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
           ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
          ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
         ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
        ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
       ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
      ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
     ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
    ✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾✾
```

### With arguments (colors not shown)

| Razor | Gotemplate
| ---   | ---
| ```@ChristmasTree(5, "Red", "Green", "a")``` | ```{{ ChristmasTree 5 "Red" "Green" "a" }}```

#### Result (small)

```text
        a
       aaa
      aaaaa
     aaaaaaa  
    aaaaaaaaa
```
