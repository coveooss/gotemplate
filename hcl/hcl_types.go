package hcl

// List allows generic lists to be rendered as HCL list
type List []interface{}

func (hl List) String() string {
	result, _ := MarshalInternal([]interface{}(hl))
	return string(result)
}

// Map allows generic maps to be rendered as HCL map
type Map map[string]interface{}

func (hm Map) String() string {
	result, _ := MarshalInternal(map[string]interface{}(hm))
	return string(result)
}
