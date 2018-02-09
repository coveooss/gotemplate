package template

func average(a interface{}, args ...interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	args = convertArgs(a, args...)
	if len(args) == 0 {
		return 0, nil
	}
	var sum interface{}
	sum, err = add(args[0], args[1:]...)
	return simplify(toFloat(sum) / float64(len(args))), nil
}

func min(values ...interface{}) interface{} { return compareNumerics(values, true) }
func max(values ...interface{}) interface{} { return compareNumerics(values, false) }
