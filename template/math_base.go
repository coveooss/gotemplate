package template

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/coveo/gotemplate/utils"
)

func add(a interface{}, args ...interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	args = convertArgs(a, args...)

	values, err := toArrayOfFloats(args...)
	if err != nil {
		if len(args) == 2 {
			// If the first argument is an array of float, we process it with the generic processor function
			if af, err := toArrayOfFloats(args[0]); err == nil {
				if _, err := strconv.ParseFloat(fmt.Sprintf("%v", args[1]), 64); err == nil {
					return processFloat2(af, args[1], func(a, b float64) float64 {
						return a + b
					})
				}
			}
		}

		switch reflect.TypeOf(args[0]).Kind() {
		case reflect.String:
			break
		case reflect.Array, reflect.Slice:
			switch reflect.TypeOf(args[1]).Kind() {
			case reflect.Array, reflect.Slice:
				return utils.MergeLists(convertArgs(args[0]), convertArgs(args[1])), nil
			default:
				return append(convertArgs(args[0]), args[1]), nil
			}
		}

		// If it is not possible to convert all arguments into numeric values
		// we simply return the concatenation of their string representation
		// This allow support of "Foo" + "Bar" or "Foo" + 1
		return fmt.Sprint(args...), nil
	}

	var result float64
	for i := range values {
		result += toFloat(values[i])
	}
	return simplify(result), nil
}

func multiply(a interface{}, args ...interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	args = convertArgs(a, args...)

	values, err := toArrayOfFloats(args...)
	if err != nil {
		if len(args) == 2 {
			// If the first argument is an array of float, we process it with the generic processor function
			if af, err := toArrayOfFloats(args[0]); err == nil {
				if _, err := strconv.ParseFloat(fmt.Sprintf("%v", args[1]), 64); err == nil {
					return processFloat2(af, args[1], func(a, b float64) float64 {
						return a * b
					})
				}
				if af2, err := toArrayOfFloats(args[1]); err == nil {
					// If the second argument is also an array of float, we then multiply the two arrays
					result := make([]interface{}, len(af2))
					for i := range af2 {
						result[i], err = multiply(af, af2[i])
					}
					return result, nil
				}
			}

			switch a := args[0].(type) {
			case string:
				return strings.Repeat(a, toInt(args[1])), nil
			default:
				result := make([]interface{}, toInt(args[1]))
				for i := range result {
					result[i] = args[0]
				}
				return result, nil
			}
		}
	}

	switch len(values) {
	case 0:
		return 0, nil
	case 2:
		return processFloat2(values[0], values[1], func(a, b float64) float64 {
			return a * b
		})
	}

	var result float64 = 1
	for i := range values {
		result *= values[i]
	}
	return simplify(result), nil
}

func subtract(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat2(a, b, func(a, b float64) float64 { return a - b })
}

func divide(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat2(a, b, func(a, b float64) float64 {
		if b == 0 {
			panic(fmt.Errorf("Division by 0"))
		}
		return a / b
	})
}

func modulo(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat2(a, b, math.Mod)
}

func modf(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return process(a, math.Modf)
}

func power(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat2(a, b, math.Pow)
}

func power10(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, func(a float64) float64 {
		return math.Pow10(int(a))
	})
}
