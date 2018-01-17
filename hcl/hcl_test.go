package hcl

/*
func TestMarshalHCLVars(t *testing.T) {
	type test struct {
		Name  string `hcl:",omitempty"`
		Value int    `hcl:",omitempty"`
	}
	const (
		noIndent = ""
		indent   = "  "
	)
	var testNilPtr *test

	type args struct {
		value  interface{}
		indent string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Integer", args{2, noIndent}, "2"},
		{"Boolean", args{true, noIndent}, "true"},
		{"String", args{"Hello world", noIndent}, `"Hello world"`},
		{"String with newline", args{"Hello\nworld", noIndent}, `"Hello\nworld"`},
		{"String with newline (pretty)", args{"Hello\n\"world\"", indent}, "<<-EOF\n  Hello\n  \"world\"\n  EOF"},
		{"Null value", args{nil, noIndent}, "null"},
		{"Null struct", args{testNilPtr, noIndent}, "null"},
		{"List of integer", args{[]int{0, 1, 2, 3}, noIndent}, "[0,1,2,3]"},
		{"Map", args{map[string]interface{}{"a": 0, "bb": 1}, noIndent}, "a=0 bb=1"},
		{"Map (pretty)", args{map[string]interface{}{"a": 0, "bb": 1}, indent}, "a  = 0\nbb = 1\n"},
		{"Structure (pretty)", args{test{"name", 1}, indent}, "Name = \"name\"\nValue = 1\n"},
		{"Structure Ptr (pretty)", args{&test{"name", 1}, indent}, "Name = \"name\"\nValue = 1\n"},
		{"Array of 1 structure (pretty)", args{[]test{{"name", 1}}, indent}, "[{\n  Name = \"name\"\n  Value = 1\n}]"},
		{"Array of 2 structures (pretty)", args{[]test{{"val1", 1}, {"val2", 1}}, indent}, "[\n  {\n    Name = \"val1\"\n    Value = 1\n  },\n  {\n    Name = \"val2\"\n    Value = 1\n  },\n]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := string(marshalHCL(toBase(tt.args.value), true, "", tt.args.indent)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalHCLVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
