package template

import (
	"fmt"
	"math"
	"strconv"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/utils"
)

func toInt(value interface{}) int {
	// We convert to the string representation to ensure that any type is converted to int
	return must(strconv.Atoi(fmt.Sprintf("%v", value))).(int)
}

func toInt64(value interface{}) int64 {
	// We convert to the string representation to ensure that any type is converted to int64
	return must(strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)).(int64)
}

func toUnsignedInteger(value interface{}) uint64 {
	// We convert to the string representation to ensure that any type is converted to uint64
	return must(strconv.ParseUint(fmt.Sprintf("%v", value), 10, 64)).(uint64)
}

func toFloat(value interface{}) float64 {
	// We convert to the string representation to ensure that any type is converted to float64
	return must(strconv.ParseFloat(fmt.Sprintf("%v", value), 64)).(float64)
}

func toArrayOfFloats(values ...interface{}) (result []float64, err error) {
	values = convertArgs(nil, values...)
	result = make([]float64, len(values))
	defer func() { err = trapError(err, recover()) }()
	for i := range values {
		result[i] = toFloat(values[i])
	}
	return
}

func process(arg, handler interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	args := convertArgs(arg)
	if err != nil {
		return
	}
	switch len(args) {
	case 0:
		r = 0
	case 1:
		r = execute(args[0], handler)
	default:
		result := make([]interface{}, len(args))
		for i := range args {
			result[i] = execute(args[i], handler)
		}
		r = result
	}
	return
}

func processFloat(arg interface{}, handler func(float64) float64) (r interface{}, err error) {
	return process(arg, handler)
}

func processFloat2(a, b interface{}, handler func(float64, float64) float64) (r interface{}, err error) {
	return processFloat(a, func(a float64) float64 {
		return handler(a, toFloat(b))
	})
}

func execute(arg, handler interface{}) interface{} {
	switch handler := handler.(type) {
	case func(float64) float64:
		return simplify(handler(toFloat(arg)))
	case func(interface{}) interface{}:
		return handler(arg)
	case func(float64) (float64, float64):
		r1, r2 := handler(toFloat(arg))
		return []interface{}{r1, r2}
	default:
		panic(fmt.Errorf("Unknown handler function %v", handler))
	}
}

func convertArgs(a interface{}, args ...interface{}) []interface{} {
	if a == nil {
		// There is no first argument, so we isolate it from the other args
		if len(args) == 0 {
			return args
		}
		a = args[0]
		args = args[1:]
	}
	if len(args) == 0 {
		// There is a single argument, we try to convert it into a list
		return collections.ToInterfaces(collections.ToStrings(a)...)
	}
	return append([]interface{}{a}, args...)
}

func simplify(value float64) interface{} {
	return utils.IIf(math.Floor(value) == value, int64(value), value)
}

func trapError(err error, rec interface{}) error {
	if rec != nil {
		switch err := rec.(type) {
		case error:
			return err
		default:
			return fmt.Errorf("%[1]T %[1]v", err)
		}
	}
	return err
}

func compareInternal(min bool, values []interface{}) interface{} {
	if len(values) == 0 {
		return nil
	}
	if values, err := toArrayOfFloats(values...); err == nil {
		result := values[0]
		for _, value := range values[1:] {
			if (min && value < result) || (!min && value > result) {
				result = value
			}
		}
		return simplify(result)
	}

	sa := collections.ToStrings(values)
	result := sa[0]
	for _, value := range sa[1:] {
		if (min && value < result) || (!min && value > result) {
			result = value
		}
	}
	return result
}

func compareNumerics(values []interface{}, min bool) interface{} {
	if len(values) == 0 {
		return nil
	}
	numerics, err := toArrayOfFloats(values...)
	if err != nil {
		return compareStrings(values, min)
	}
	result := numerics[0]
	comp := iif(min, math.Min, math.Max).(func(a, b float64) float64)
	for _, value := range numerics[1:] {
		result = comp(result, value)
	}
	return simplify(result)
}

func compareStrings(values []interface{}, min bool) (result string) {
	sa := collections.ToStrings(values)
	result = sa[0]
	for _, value := range sa[1:] {
		if (min && value < result) || (!min && value > result) {
			result = value
		}
	}
	return result
}

func generateNumericArray(limit bool, params ...interface{}) (result collections.IGenericList, err error) {
	defer func() { err = trapError(err, recover()) }()

	var start, stop float64
	var step float64 = 1
	var precision int
	switch len(params) {
	case 1:
		start = float64(iif(limit, 1, 0).(int))
		stop = toFloat(params[0])
	case 3:
		step = math.Abs(toFloat(params[2]))
		_, frac := collections.Split2(fmt.Sprintf("%g", step), ".")
		precision = len(frac)
		fallthrough
	case 2:
		start = toFloat(params[0])
		stop = toFloat(params[1])
	default:
		return nil, fmt.Errorf("Invalid arguments, must be start [stop] [step]")
	}
	if step == 0 {
		return nil, fmt.Errorf("Step cannot be zero")
	}
	array := make([]interface{}, 0, int64(math.Abs(stop-start)))
	forward := stop > start
	if !forward {
		step = -step
	}
	for current := start; (forward && current <= stop || !forward && current >= stop) && (limit || current != stop); {
		current = sprigRound(current, precision)
		array = append(array, simplify(current))
		current += step
	}
	result = collections.AsList(array)
	return
}
