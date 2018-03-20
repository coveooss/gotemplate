# Data collections (lists/slices and dicts/maps)

## Maps
| Razor | Gotemplate | Note |
| --- | --- | --- |
| @.razorDict := dict("test", 1, "test2", 2); | {{- set . "goDict" (dict "test" 1 "test2" 2) }} | Creation |
| @.razorDict2 := dict("test3", 3); | {{- set . "goDict2" (dict "test3" 3) }} | Creation |
| @.razorDict.test2 := 3; | {{- set . "goDict.test2" 3 }} | Update |
| @.razorDict.test3 := 4; | {{- set . "goDict.test3" 4 }} | Update |
| @.razorDict := merge(.razorDict, .razorDict2); | {{- set . "goDict" (merge .goDict .goDict2) }} | Merge (First dict has priority) |
| @.razorDict.test3; | {{ get . "goDict.test3" }} | Should be 4
