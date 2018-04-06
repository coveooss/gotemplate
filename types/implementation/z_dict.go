package implementation

/*
// Only convert object to list or dictionary implementation (not deep conversion).
func (g genImpl) list(l baseIList) *baseList {
	result := baseList(*l.AsArray())
	return &result
}
func (g genImpl) dict(d baseIDict) baseIDict { return baseDict(d.AsMap()) }

// Deep converts an baseIDict object into current dictionary implementation.
func (g genImpl) Dict(d baseIDict) baseIDict {
	keys := d.KeysAsString()
	for i := range keys {
		if value, converted := g.value(d.Get(keys[i])); converted {
			d.Set(keys[i], value)
		}
	}
	return g.dict(d)
}

// Deep converts an baseIList object into current list implementation.
func (g genImpl) List(l baseIList) baseIList {
	for i := range *l.AsArray() {
		if value, converted := g.value(l.Get(i)); converted {
			l.Set(i, value)
		}
	}
	return g.list(l)
}

// Converts any convertible object into current list and dictionary implementation.
func (g genImpl) value(value interface{}) (interface{}, bool) {
}
*/
