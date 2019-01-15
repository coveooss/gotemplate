```go
@Add(interface{}, interface{}) IDictionary
AsMap() map[string]interface{}
Clone(...interface{}) IDictionary
Count() int
Create(...int) IDictionary
CreateList(...int) IGenericList
Default(interface{}, interface{}) interface{}
Delete(interface{}, ...interface{}) IDictionary, error
Flush(...interface{}) IDictionary
Get(...interface{}) interface{}
GetHelpers() IDictionaryHelper, IListHelper
GetKeys() IGenericList
GetValues() IGenericList
Has(...interface{}) bool
KeysAsString() StringArray
Len() int
Merge(IDictionary, ...IDictionary) IDictionary
Native() interface{}
Omit(interface{}, ...interface{}) IDictionary
Pop(...interface{}) interface{}
PrettyPrint() string
Set(interface{}, interface{}) IDictionary
String() string
Transpose() IDictionary
TypeName() String
```
