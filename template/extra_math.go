package template

import (
	"math"

	"github.com/Masterminds/sprig"
)

const (
	mathBase         = "Mathematic Fundamental"
	mathStatistics   = "Mathematic Stats"
	mathTrigonometry = "Mathematic Trigonometry"
	mathBits         = "Mathematic Bit Operations"
	mathUtilities    = "Mathematic Utilities"
)

var mathFuncs = funcTableMap{
	// Base
	"add":   {f: add, group: mathBase, aliases: []string{"sum"}, args: []string{}, desc: ""},
	"ceil":  {f: ceil, group: mathBase, aliases: []string{"roundUp", "roundup"}, args: []string{}, desc: ""},
	"dim":   {f: dim, group: mathBase, aliases: []string{}, args: []string{}, desc: ""},
	"div":   {f: divide, group: mathBase, aliases: []string{"divide", "quotient"}, args: []string{}, desc: ""},
	"exp":   {f: exp, group: mathBase, aliases: []string{"exponent"}, args: []string{}, desc: ""},
	"exp2":  {f: exp2, group: mathBase, aliases: []string{"exponent2"}, args: []string{}, desc: ""},
	"expm1": {f: expm1, group: mathBase, aliases: []string{}, args: []string{}, desc: ""},
	"floor": {f: floor, group: mathBase, aliases: []string{"roundDown", "rounddown", "int", "integer"}, args: []string{}, desc: ""},
	"mod":   {f: modulo, group: mathBase, aliases: []string{"modulo"}, args: []string{}, desc: ""},
	"modf":  {f: modf, group: mathBase, aliases: []string{}, args: []string{}, desc: ""},
	"mul":   {f: multiply, group: mathBase, aliases: []string{"multiply", "prod", "product"}, args: []string{}, desc: ""},
	"pow":   {f: power, group: mathBase, aliases: []string{"power"}, args: []string{}, desc: ""},
	"pow10": {f: power10, group: mathBase, aliases: []string{"power10"}, args: []string{}, desc: ""},
	"rem":   {f: remainder, group: mathBase, aliases: []string{"remainder"}, args: []string{}, desc: ""},
	"sub":   {f: subtract, group: mathBase, aliases: []string{"subtract"}, args: []string{}, desc: ""},
	"trunc": {f: trunc, group: mathBase, aliases: []string{"truncate"}, args: []string{}, desc: ""},

	// Statistics
	"avg": {f: average, group: mathStatistics, aliases: []string{"average"}, args: []string{}, desc: ""},
	"min": {f: min, group: mathStatistics, aliases: []string{"minimum", "smallest"}, args: []string{}, desc: ""},
	"max": {f: max, group: mathStatistics, aliases: []string{"maximum", "biggest"}, args: []string{}, desc: ""},

	// Trigonometry
	"rad":    {f: rad, group: mathTrigonometry, aliases: []string{"radian"}, args: []string{}, desc: ""},
	"deg":    {f: deg, group: mathTrigonometry, aliases: []string{"degree"}, args: []string{}, desc: ""},
	"acos":   {f: acos, group: mathTrigonometry, aliases: []string{"arcCosine", "arcCosinus"}, args: []string{}, desc: ""},
	"acosh":  {f: acosh, group: mathTrigonometry, aliases: []string{"arcHyperbolicCosine", "arcHyperbolicCosinus"}, args: []string{}, desc: ""},
	"asin":   {f: asin, group: mathTrigonometry, aliases: []string{"arcSine", "arcSinus"}, args: []string{}, desc: ""},
	"asinh":  {f: asinh, group: mathTrigonometry, aliases: []string{"arcHyperbolicSine", "arcHyperbolicSinus"}, args: []string{}, desc: ""},
	"atan":   {f: atan, group: mathTrigonometry, aliases: []string{"arcTangent"}, args: []string{}, desc: ""},
	"atan2":  {f: atan2, group: mathTrigonometry, aliases: []string{"arcTangent2"}, args: []string{}, desc: ""},
	"atanh":  {f: atanh, group: mathTrigonometry, aliases: []string{"arcHyperbolicTangent"}, args: []string{}, desc: ""},
	"cos":    {f: cos, group: mathTrigonometry, aliases: []string{"cosine", "cosinus"}, args: []string{}, desc: ""},
	"cosh":   {f: cosh, group: mathTrigonometry, aliases: []string{"hyperbolicCosine", "hyperbolicCosinus"}, args: []string{}, desc: ""},
	"ilogb":  {f: ilogb, group: mathTrigonometry, aliases: []string{}, args: []string{}, desc: ""},
	"log":    {f: logFunc, group: mathTrigonometry, aliases: []string{}, args: []string{}, desc: ""},
	"log10":  {f: log10, group: mathTrigonometry, aliases: []string{}, args: []string{}, desc: ""},
	"log1p":  {f: log1p, group: mathTrigonometry, aliases: []string{}, args: []string{}, desc: ""},
	"log2":   {f: log2, group: mathTrigonometry, aliases: []string{}, args: []string{}, desc: ""},
	"logb":   {f: logb, group: mathTrigonometry, aliases: []string{}, args: []string{}, desc: ""},
	"sin":    {f: sin, group: mathTrigonometry, aliases: []string{"sine", "sinus"}, args: []string{}, desc: ""},
	"sincos": {f: sincos, group: mathTrigonometry, aliases: []string{"sineCosine", "sinusCosinus"}, args: []string{}, desc: ""},
	"sinh":   {f: sinh, group: mathTrigonometry, aliases: []string{"hyperbolicSine", "hyperbolicSinus"}, args: []string{}, desc: ""},
	"tan":    {f: tan, group: mathTrigonometry, aliases: []string{"tangent"}, args: []string{}, desc: ""},
	"tanh":   {f: tanh, group: mathTrigonometry, aliases: []string{"hyperbolicTangent"}, args: []string{}, desc: ""},

	// Binary operators
	"lshift": {f: leftShift, group: mathBits, aliases: []string{"leftShift"}, args: []string{}, desc: ""},
	"rshift": {f: rightShift, group: mathBits, aliases: []string{"rightShift"}, args: []string{}, desc: ""},
	"bor":    {f: bitwiseOr, group: mathBits, aliases: []string{"bitwiseOR"}, args: []string{}, desc: ""},
	"band":   {f: bitwiseAnd, group: mathBits, aliases: []string{"bitwiseAND"}, args: []string{}, desc: ""},
	"bxor":   {f: bitwiseXor, group: mathBits, aliases: []string{"bitwiseXOR"}, args: []string{}, desc: ""},
	"bclear": {f: bitwiseClear, group: mathBits, aliases: []string{"bitwiseClear"}, args: []string{}, desc: ""},

	// Utilities
	"abs":       {f: abs, group: mathUtilities, aliases: []string{"absolute"}, args: []string{}, desc: ""},
	"sqrt":      {f: sqrt, group: mathUtilities, aliases: []string{"squareRoot"}, args: []string{}, desc: ""},
	"to":        {f: to, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"until":     {f: until, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"frexp":     {f: frexp, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"ldexp":     {f: ldexp, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"gamma":     {f: gamma, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"lgamma":    {f: lgamma, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"hypot":     {f: hypot, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"isInf":     {f: isInfinity, group: mathUtilities, aliases: []string{"isInfinity"}, args: []string{}, desc: ""},
	"isNaN":     {f: isNaN, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"nextAfter": {f: nextAfter, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"signBit":   {f: signBit, group: mathUtilities, aliases: []string{}, args: []string{}, desc: ""},
	"hex":       {f: hex, group: mathUtilities, aliases: []string{"hexa", "hexaDecimal"}, args: []string{}, desc: ""},
	"dec":       {f: decimal, group: mathUtilities, aliases: []string{"decimal"}, args: []string{}, desc: ""},
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

	if !t.mathConstantInjected {
		// We do not want to inject the math constant twice
		t.setConstant(true, constants, "Math", "MATH")
		t.mathConstantInjected = true
	}
}

var round = sprig.GenericFuncMap()["round"].(func(a interface{}, p int, r_opt ...float64) float64)

// To classify
func to(params ...interface{}) (interface{}, error)    { return generateNumericArray(true, params...) }
func until(params ...interface{}) (interface{}, error) { return generateNumericArray(false, params...) }

func abs(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Abs)
}

func ceil(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Ceil)
}

// math.Cbrt
// math.Copysign

func dim(a, b interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat2(a, b, math.Dim)
}

// math.Erf
// math.Erfc

func exp(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Exp)
}

func exp2(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Exp2)
}

func expm1(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Expm1)
}

func floor(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Floor)
}

func frexp(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	f, e := math.Frexp(toFloat(a))
	return []interface{}{simplify(f), e}, nil
}

func gamma(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Gamma)
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

func trunc(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Trunc)
}

// math.Y0
// math.Y2
// math.Yn
