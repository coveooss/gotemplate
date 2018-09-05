```go
Center(int) String
Compare(string) int
Contains(string) bool
ContainsAny(string) bool
ContainsRune(int32) bool
Count(string) int
EqualFold(string) bool
Fields() StringArray
FieldsFunc(func(int32) bool) StringArray
FieldsID() StringArray
GetContextAtPosition(int, string, string) String, int
GetWordAtPosition(int, ...string) String, int
HasPrefix(string) bool
HasSuffix(string) bool
Indent(string) String
IndentN(int) String
Index(string) int
IndexAny(string) int
IndexByte(uint8) int
IndexFunc(func(int32) bool) int
IndexRune(int32) int
Join(...string) String
LastIndex(string) int
LastIndexAny(string) int
LastIndexByte(uint8) int
LastIndexFunc(func(int32) bool) int
Len() int
Lines() StringArray
Map(func(int32) int32) String
Repeat(int) String
Replace(string, string) String
ReplaceN(string, string, int) String
SelectContext(int, string, string) String
SelectWord(int, ...string) String
Split(string) StringArray
SplitAfter(string) StringArray
SplitAfterN(string, int) StringArray
SplitN(string, int) StringArray
Str() string
String() string
Title() String
ToLower() String
ToLowerSpecial(unicode.SpecialCase) String
ToTitle() String
ToTitleSpecial(unicode.SpecialCase) String
ToUpper() String
ToUpperSpecial(unicode.SpecialCase) String
Trim(string) String
TrimFunc(func(int32) bool) String
TrimLeft(string) String
TrimLeftFunc(func(int32) bool) String
TrimPrefix(string) String
TrimRight(string) String
TrimRightFunc(func(int32) bool) String
TrimSpace() String
TrimSuffix(string) String
UnIndent() String
Wrap(int) String
```
