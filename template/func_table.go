package template

type funcTable struct {
	function    interface{}
	group       string
	aliases     []string
	description string
}

type funcTableMap map[string]funcTable

func (ftm funcTableMap) convert() map[string]interface{} {
	result := make(map[string]interface{}, len(ftm))
	for key, val := range ftm {
		result[key] = val.function
		for i := range val.aliases {
			result[val.aliases[i]] = val.function
		}
	}
	return result
}
