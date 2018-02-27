package template

import (
	"math"

	"github.com/Masterminds/sprig"
)

const (
	mathBase         = "Mathematic (base)"
	mathStatistics   = "Mathematic (stats)"
	mathTrigonometry = "Mathematic (trigonometry)"
	mathBits         = "Mathematic (bit operations)"
	mathConstants    = "Mathematic (constants)"
	mathUtilities    = "Mathematic (utilities)"
)

var mathFuncs = funcTableMap{
	// Base
	"add":   {add, mathBase, []string{"sum"}, []string{}, ""},
	"sub":   {subtract, mathBase, []string{"subtract"}, []string{}, ""},
	"div":   {divide, mathBase, []string{"divide", "quotient"}, []string{}, ""},
	"mul":   {multiply, mathBase, []string{"multiply", "prod", "product"}, []string{}, ""},
	"mod":   {modulo, mathBase, []string{"modulo"}, []string{}, ""},
	"modf":  {modf, mathBase, []string{}, []string{}, ""},
	"rem":   {remainder, mathBase, []string{"remainder"}, []string{}, ""},
	"pow":   {power, mathBase, []string{"power"}, []string{}, ""},
	"pow10": {power10, mathBase, []string{"power10"}, []string{}, ""},
	"exp":   {exp, mathBase, []string{"exponent"}, []string{}, ""},
	"exp2":  {exp2, mathBase, []string{"exponent2"}, []string{}, ""},
	"expm1": {expm1, mathBase, []string{}, []string{}, ""},

	// Statistics
	"avg": {average, mathStatistics, []string{"average"}, []string{}, ""},
	"min": {min, mathStatistics, []string{"minimum", "smallest"}, []string{}, ""},
	"max": {max, mathStatistics, []string{"maximum", "biggest"}, []string{}, ""},

	// Trigonometry
	"rad":    {rad, mathTrigonometry, []string{"radian"}, []string{}, ""},
	"deg":    {deg, mathTrigonometry, []string{"degree"}, []string{}, ""},
	"acos":   {acos, mathTrigonometry, []string{"arcCosine", "arcCosinus"}, []string{}, ""},
	"acosh":  {acosh, mathTrigonometry, []string{"arcHyperbolicCosine", "arcHyperbolicCosinus"}, []string{}, ""},
	"asin":   {asin, mathTrigonometry, []string{"arcSine", "arcSinus"}, []string{}, ""},
	"asinh":  {asinh, mathTrigonometry, []string{"arcHyperbolicSine", "arcHyperbolicSinus"}, []string{}, ""},
	"atan":   {atan, mathTrigonometry, []string{"arcTangent"}, []string{}, ""},
	"atan2":  {atan2, mathTrigonometry, []string{"arcTangent2"}, []string{}, ""},
	"atanh":  {atanh, mathTrigonometry, []string{"arcHyperbolicTangent"}, []string{}, ""},
	"cos":    {cos, mathTrigonometry, []string{"cosine", "cosinus"}, []string{}, ""},
	"cosh":   {cosh, mathTrigonometry, []string{"hyperbolicCosine", "hyperbolicCosinus"}, []string{}, ""},
	"sin":    {sin, mathTrigonometry, []string{"sine", "sinus"}, []string{}, ""},
	"sinh":   {sinh, mathTrigonometry, []string{"hyperbolicSine", "hyperbolicSinus"}, []string{}, ""},
	"sincos": {sincos, mathTrigonometry, []string{"sineCosine", "sinusCosinus"}, []string{}, ""},
	"ilogb":  {ilogb, mathTrigonometry, []string{}, []string{}, ""},
	"log":    {logFunc, mathTrigonometry, []string{}, []string{}, ""},
	"log10":  {log10, mathTrigonometry, []string{}, []string{}, ""},
	"log1p":  {log1p, mathTrigonometry, []string{}, []string{}, ""},
	"log2":   {log2, mathTrigonometry, []string{}, []string{}, ""},
	"logb":   {logb, mathTrigonometry, []string{}, []string{}, ""},
	"tan":    {tan, mathTrigonometry, []string{"tangent"}, []string{}, ""},
	"tanh":   {tanh, mathTrigonometry, []string{"hyperbolicTangent"}, []string{}, ""},

	// Binary operators
	"lshift": {leftShift, mathBits, []string{"leftShift"}, []string{}, ""},
	"rshift": {rightShift, mathBits, []string{"rightShift"}, []string{}, ""},
	"bor":    {bitwiseOr, mathBits, []string{"bitwiseOR"}, []string{}, ""},
	"band":   {bitwiseAnd, mathBits, []string{"bitwiseAND"}, []string{}, ""},
	"bxor":   {bitwiseXor, mathBits, []string{"bitwiseXOR"}, []string{}, ""},
	"bclear": {bitwiseClear, mathBits, []string{"bitwiseClear"}, []string{}, ""},

	// Utilities
	"abs":       {abs, mathUtilities, []string{}, []string{}, ""},
	"sqrt":      {sqrt, mathUtilities, []string{}, []string{}, ""},
	"to":        {to, mathUtilities, []string{}, []string{}, ""},
	"until":     {until, mathUtilities, []string{}, []string{}, ""},
	"frexp":     {frexp, mathUtilities, []string{}, []string{}, ""},
	"ldexp":     {ldexp, mathUtilities, []string{}, []string{}, ""},
	"gamma":     {gamma, mathUtilities, []string{}, []string{}, ""},
	"lgamma":    {lgamma, mathUtilities, []string{}, []string{}, ""},
	"hypot":     {hypot, mathUtilities, []string{}, []string{}, ""},
	"isInf":     {isInfinity, mathUtilities, []string{"isInfinity"}, []string{}, ""},
	"isNaN":     {isNaN, mathUtilities, []string{}, []string{}, ""},
	"nextAfter": {nextAfter, mathUtilities, []string{}, []string{}, ""},
	"signBit":   {signBit, mathUtilities, []string{}, []string{}, ""},
	"hex":       {hex, mathUtilities, []string{"hexa", "hexaDecimal"}, []string{}, ""},
	"dec":       {decimal, mathUtilities, []string{"decimal"}, []string{}, ""},
}

func (t *Template) addMathFuncs() {
	// Enhance mathematic functions
	t.AddFunctions(mathFuncs)

	constants := map[string]interface{}{
		"E":                      math.E,
		"Pi":                     math.Pi,
		"Phi":                    math.Phi,
		"Sqrt2":                  math.Sqrt2,
		"SqrtE":                  math.SqrtE,
		"SqrtPi":                 math.SqrtPi,
		"SqrtPhi":                math.SqrtPhi,
		"Ln2":                    math.Ln2,
		"Log2E":                  math.Log2E,
		"Ln10":                   math.Ln10,
		"Log10E":                 math.Log10E,
		"MaxFloat32":             math.MaxFloat32,
		"MaxFloat64":             math.MaxFloat64,
		"SmallestNonzeroFloat64": math.SmallestNonzeroFloat64,
		"MaxInt8":                math.MaxInt8,
		"MaxInt16":               math.MaxInt16,
		"MaxInt32":               math.MaxInt32,
		"MaxInt64":               math.MaxInt64,
		"MaxUint8":               math.MaxUint8,
		"MaxUint16":              math.MaxUint16,
		"MaxUint32":              math.MaxUint32,
		"MaxUint64":              uint(math.MaxUint64),
		"Nan":                    math.NaN(),
		"Infinity":               math.Inf(1),
		"Inf":                    math.Inf(1),
		"NegativeInfinity":       math.Inf(-1),
		"NegInf":                 math.Inf(-1),
	}

	if !constantInjected {
		// We do not want to inject the math constant twice
		t.setConstant(true, constants, "Math", "MATH")
		constantInjected = true
	}
}

var constantInjected bool
var round = sprig.GenericFuncMap()["round"].(func(a interface{}, p int, r_opt ...float64) float64)

// To classify
func to(params ...interface{}) (interface{}, error)    { return generateNumericArray(true, params...) }
func until(params ...interface{}) (interface{}, error) { return generateNumericArray(false, params...) }

func abs(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Abs(toFloat(a))), nil
}

// math.Cbrt
// math.Ceil
// math.Copysign

// math.Dim
// math.Erf
// math.Erfc

func exp(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Exp(toFloat(a))), nil
}

func exp2(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Exp2(toFloat(a))), nil
}

func expm1(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Expm1(toFloat(a))), nil
}

func frexp(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	f, e := math.Frexp(toFloat(a))
	return []interface{}{simplify(f), e}, nil
}

func gamma(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Gamma(toFloat(a))), nil
}

func infinity(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Inf(toInt(a))), nil
}

func isInfinity(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return math.IsInf(toFloat(a), toInt(b)), nil
}

func isNaN(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return math.IsNaN(toFloat(a)), nil
}

// math.J0
// math.J1
// math.Jn

func ldexp(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Ldexp(toFloat(a), toInt(b))), nil
}

func lgamma(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	f, e := math.Lgamma(toFloat(a))
	return []interface{}{simplify(f), e}, nil
}

func nextAfter(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Nextafter(toFloat(a), toFloat(b))), nil
}

func remainder(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return simplify(math.Remainder(toFloat(a), toFloat(b))), nil
}

func signBit(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return math.Signbit(toFloat(a)), nil
}

func sqrt(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Sqrt)
}

// math.Trunc
// math.Y0
// math.Y2
// math.Yn
