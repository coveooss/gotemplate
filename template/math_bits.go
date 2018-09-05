package template

import (
	"fmt"
	"strconv"
	"strings"
)

func leftShift(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return process(a, func(a interface{}) interface{} { return toInt64(a) << toUnsignedInteger(b) })
}

func rightShift(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return process(a, func(a interface{}) interface{} { return toInt64(a) >> toUnsignedInteger(b) })
}

func bitwiseOr(a, b interface{}, rest ...interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	result := toInt64(a) | toInt64(b)
	for i := range rest {
		result = result | toInt64(rest[i])
	}
	return result, nil
}

func bitwiseAnd(a, b interface{}, rest ...interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	result := toInt64(a) & toInt64(b)
	for i := range rest {
		result = result & toInt64(rest[i])
	}
	return result, nil
}

func bitwiseXor(a, b interface{}, rest ...interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	result := toInt64(a) ^ toInt64(b)
	for i := range rest {
		result = result ^ toInt64(rest[i])
	}
	return result, nil
}

func bitwiseClear(a, b interface{}, rest ...interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	result := toInt64(a) &^ toInt64(b)
	for i := range rest {
		result = result &^ toInt64(rest[i])
	}
	return result, nil
}

func hex(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return process(a, func(a interface{}) interface{} {
		return fmt.Sprintf("0x%X", toInt(a))
	})
}

func decimal(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return process(a, func(a interface{}) interface{} {
		return must(strconv.ParseInt(strings.TrimPrefix(fmt.Sprint(a), "0x"), 16, 64))
	})
}
