package template

import (
	"fmt"
	"math"
	"strconv"

	"github.com/coveo/gotemplate/collections"
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

func toListOfFloats(values iList) (result iList, err error) {
	if values == nil {
		return collections.CreateList(), nil
	}
	values = convertArgs(nil, values.AsArray()...)
	result = values.Clone()
	defer func() {
		if err = trapError(err, recover()); err != nil {
			result = nil
		}
	}()
	for i := range result.AsArray() {
		result.Set(i, toFloat(result.Get(i)))
	}
	return
}

func asFloats(values iList) ([]float64, error) {
	result, err := toListOfFloats(values)
	if err != nil {
		return nil, err
	}
	return mustAsFloats(result), nil
}

func mustAsFloats(values iList) (result []float64) {
	result = make([]float64, values.Len())
	for i := range result {
		result[i] = values.Get(i).(float64)
	}
	return
}

func process(arg, handler interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	arguments := convertArgs(arg)
	if arguments.Len() == 0 {
		return
	}
	argArray := arguments.AsArray()
	switch len(argArray) {
	case 0:
		r = 0
	case 1:
		r = execute(argArray[0], handler)
	default:
		result := arguments.Clone()
		for i := range result.AsArray() {
			result.Set(i, execute(result.Get(i), handler))
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

func convertArgs(arg1 interface{}, args ...interface{}) (result collections.IGenericList) {
	if arg1 == nil {
		// There is no first argument, so we isolate it from the other args
		if len(args) == 0 {
			return collections.CreateList()
		}
		arg1, args = args[0], args[1:]
	}
	if len(args) == 0 {
		// There is a single argument, we try to convert it into a list
		return collections.AsList(arg1)
	}

	if list, err := collections.TryAsList(arg1); err == nil {
		return list.Create(0, len(args)+1).Append(arg1).Append(args...)
	}
	return collections.NewList(arg1).Append(args...)
}

func simplify(value float64) interface{} {
	return iif(math.Floor(value) == value, int64(value), value)
}

func compareNumerics(values []interface{}, useMinFunc bool) interface{} {
	if len(values) == 0 {
		return nil
	}
	numerics, err := asFloats(collections.AsList(values))
	if err != nil {
		return compareStrings(values, useMinFunc)
	}
	result := numerics[0]
	comp := iif(useMinFunc, math.Min, math.Max).(func(a, b float64) float64)
	for _, value := range numerics[1:] {
		result = comp(result, value)
	}
	return simplify(result)
}

func compareStrings(values []interface{}, useMinFunc bool) (result string) {
	sa := collections.ToStrings(values)
	result = sa[0]
	for _, value := range sa[1:] {
		if (useMinFunc && value < result) || (!useMinFunc && value > result) {
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
