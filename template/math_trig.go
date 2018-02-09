package template

import (
	"math"
)

func rad(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, func(a float64) float64 {
		return a * math.Pi / 180
	})
}

func deg(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, func(a float64) float64 {
		return a * 180 / math.Pi
	})
}

func acos(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Acos)
}

func acosh(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Acosh)
}

func asin(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Asin)
}

func asinh(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Asinh)
}

func atan(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Atan)
}

func atan2(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat2(a, b, math.Atan2)
}

func atanh(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Atanh)
}

func cos(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Cos)
}

func cosh(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Cosh)
}

func hypot(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat2(a, b, math.Hypot)
}

func ilogb(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return math.Ilogb(toFloat(a)), nil
}

func logFunc(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Log)
}

func log10(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Log10)
}

func log1p(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Log1p)
}

func log2(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Log2)
}

func logb(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Logb)
}

func sin(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Sin)
}

func sincos(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	s, c := math.Sincos(toFloat(a))
	return []interface{}{simplify(s), simplify(c)}, nil
}

func sinh(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Sinh)
}

func tan(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Tan)
}

func tanh(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Tanh)
}
