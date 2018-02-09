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
	"add":   {add, mathBase, []string{"sum"}, ""},
	"sub":   {subtract, mathBase, []string{"subtract"}, ""},
	"div":   {divide, mathBase, []string{"divide", "quotient"}, ""},
	"mul":   {multiply, mathBase, []string{"multiply", "prod", "product"}, ""},
	"mod":   {modulo, mathBase, []string{"modulo"}, ""},
	"modf":  {modf, mathBase, []string{}, ""},
	"rem":   {remainder, mathBase, []string{"remainder"}, ""},
	"pow":   {power, mathBase, []string{"power"}, ""},
	"pow10": {power10, mathBase, []string{"power10"}, ""},
	"exp":   {exp, mathBase, []string{"exponent"}, ""},
	"exp2":  {exp2, mathBase, []string{"exponent2"}, ""},
	"expm1": {expm1, mathBase, []string{}, ""},

	// Statistics
	"avg": {average, mathStatistics, []string{"average"}, ""},
	"min": {min, mathStatistics, []string{"minimum", "smallest"}, ""},
	"max": {max, mathStatistics, []string{"maximum", "biggest"}, ""},

	// Trigonometry
	"rad":    {rad, mathTrigonometry, []string{"radian"}, ""},
	"deg":    {deg, mathTrigonometry, []string{"degree"}, ""},
	"acos":   {acos, mathTrigonometry, []string{"arcCosine", "arcCosinus"}, ""},
	"acosh":  {acosh, mathTrigonometry, []string{"arcHyperbolicCosine", "arcHyperbolicCosinus"}, ""},
	"asin":   {asin, mathTrigonometry, []string{"arcSine", "arcSinus"}, ""},
	"asinh":  {asinh, mathTrigonometry, []string{"arcHyperbolicSine", "arcHyperbolicSinus"}, ""},
	"atan":   {atan, mathTrigonometry, []string{"arcTangent"}, ""},
	"atan2":  {atan2, mathTrigonometry, []string{"arcTangent2"}, ""},
	"atanh":  {atanh, mathTrigonometry, []string{"arcHyperbolicTangent"}, ""},
	"cos":    {cos, mathTrigonometry, []string{"cosine", "cosinus"}, ""},
	"cosh":   {cosh, mathTrigonometry, []string{"hyperbolicCosine", "hyperbolicCosinus"}, ""},
	"sin":    {sin, mathTrigonometry, []string{"sine", "sinus"}, ""},
	"sinh":   {sinh, mathTrigonometry, []string{"hyperbolicSine", "hyperbolicSinus"}, ""},
	"sincos": {sincos, mathTrigonometry, []string{"sineCosine", "sinusCosinus"}, ""},
	"ilogb":  {ilogb, mathTrigonometry, []string{}, ""},
	"log":    {logFunc, mathTrigonometry, []string{}, ""},
	"log10":  {log10, mathTrigonometry, []string{}, ""},
	"log1p":  {log1p, mathTrigonometry, []string{}, ""},
	"log2":   {log2, mathTrigonometry, []string{}, ""},
	"logb":   {logb, mathTrigonometry, []string{}, ""},
	"tan":    {tan, mathTrigonometry, []string{"tangent"}, ""},
	"tanh":   {tanh, mathTrigonometry, []string{"hyperbolicTangent"}, ""},

	// Binary operators
	"lshift": {leftShift, mathBits, []string{"leftShift"}, ""},
	"rshift": {rightShift, mathBits, []string{"rightShift"}, ""},
	"bor":    {bitwiseOr, mathBits, []string{"bitwiseOR"}, ""},
	"band":   {bitwiseAnd, mathBits, []string{"bitwiseAND"}, ""},
	"bxor":   {bitwiseXor, mathBits, []string{"bitwiseXOR"}, ""},
	"bclear": {bitwiseClear, mathBits, []string{"bitwiseClear"}, ""},

	// Constants
	"E":               {func() float64 { return math.E }, mathConstants, []string{}, ""},
	"Pi":              {func() float64 { return math.Pi }, mathConstants, []string{}, ""},
	"Phi":             {func() float64 { return math.Phi }, mathConstants, []string{}, ""},
	"Sqrt2":           {func() float64 { return math.Sqrt2 }, mathConstants, []string{}, ""},
	"SqrtE":           {func() float64 { return math.SqrtE }, mathConstants, []string{}, ""},
	"SqrtPi":          {func() float64 { return math.SqrtPi }, mathConstants, []string{}, ""},
	"SqrtPhi":         {func() float64 { return math.SqrtPhi }, mathConstants, []string{}, ""},
	"Ln2":             {func() float64 { return math.Ln2 }, mathConstants, []string{}, ""},
	"Log2E":           {func() float64 { return math.Log2E }, mathConstants, []string{}, ""},
	"Ln10":            {func() float64 { return math.Ln10 }, mathConstants, []string{}, ""},
	"Log10E":          {func() float64 { return math.Log10E }, mathConstants, []string{}, ""},
	"MaxFloat":        {func() float64 { return math.MaxFloat64 }, mathConstants, []string{}, ""},
	"MinNonZeroFloat": {func() float64 { return math.SmallestNonzeroFloat64 }, mathConstants, []string{}, ""},
	"MaxInt":          {func() float64 { return math.MaxInt64 }, mathConstants, []string{}, ""},
	"MaxUInt":         {func() float64 { return math.MaxUint64 }, mathConstants, []string{}, ""},
	"Nan":             {func() float64 { return math.NaN() }, mathConstants, []string{}, ""},
	"inf":             {infinity, mathConstants, []string{"infinity"}, ""},

	// Utilities
	"abs":       {abs, mathUtilities, []string{}, ""},
	"sqrt":      {sqrt, mathUtilities, []string{}, ""},
	"to":        {to, mathUtilities, []string{}, ""},
	"until":     {until, mathUtilities, []string{}, ""},
	"frexp":     {frexp, mathUtilities, []string{}, ""},
	"ldexp":     {ldexp, mathUtilities, []string{}, ""},
	"gamma":     {gamma, mathUtilities, []string{}, ""},
	"lgamma":    {lgamma, mathUtilities, []string{}, ""},
	"hypot":     {hypot, mathUtilities, []string{}, ""},
	"isInf":     {isInfinity, mathUtilities, []string{"isInfinity"}, ""},
	"isNaN":     {isNaN, mathUtilities, []string{}, ""},
	"nextAfter": {nextAfter, mathUtilities, []string{}, ""},
	"signBit":   {signBit, mathUtilities, []string{}, ""},
	"hex":       {hex, mathUtilities, []string{"hexa", "hexaDecimal"}, ""},
	"dec":       {decimal, mathUtilities, []string{"decimal"}, ""},
}

var mathFuncsInject map[string]interface{}

func (t *Template) addMathFuncs() {
	if mathFuncsInject == nil {
		mathFuncsInject = mathFuncs.convert()
	}
	// Enhance mathematic functions
	t.Funcs(mathFuncsInject)
}

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
