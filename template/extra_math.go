package template

import (
	"math"
)

const (
	mathBase         = "Mathematic Fundamental"
	mathBits         = "Mathematic Bit Operations"
	mathStatistics   = "Mathematic Stats"
	mathTrigonometry = "Mathematic Trigonometry"
	mathUtilities    = "Mathematic Utilities"
)

var mathBaseFuncs = dictionary{
	"add":   add,
	"ceil":  ceil,
	"cbrt":  cbrt,
	"dim":   dim,
	"div":   divide,
	"exp":   exp,
	"exp2":  exp2,
	"expm1": expm1,
	"floor": floor,
	"mod":   modulo,
	"modf":  modf,
	"mul":   multiply,
	"pow":   power,
	"pow10": power10,
	"rem":   remainder,
	"sub":   subtract,
	"trunc": trunc,
}

var mathStatFuncs = dictionary{
	"avg": average,
	"max": max,
	"min": min,
}

var mathTrigFuncs = dictionary{
	"acos":   acos,
	"acosh":  acosh,
	"asin":   asin,
	"asinh":  asinh,
	"atan":   atan,
	"atan2":  atan2,
	"atanh":  atanh,
	"cos":    cos,
	"cosh":   cosh,
	"deg":    deg,
	"ilogb":  ilogb,
	"j0":     j0,
	"j1":     j1,
	"jn":     jn,
	"log":    logFunc,
	"log10":  log10,
	"log1p":  log1p,
	"log2":   log2,
	"logb":   logb,
	"rad":    rad,
	"sin":    sin,
	"sincos": sincos,
	"sinh":   sinh,
	"tan":    tan,
	"tanh":   tanh,
	"y0":     y0,
	"y1":     y1,
	"yn":     yn,
}

var mathBitsFuncs = dictionary{
	"band":   bitwiseAnd,
	"bclear": bitwiseClear,
	"bor":    bitwiseOr,
	"bxor":   bitwiseXor,
	"lshift": leftShift,
	"rshift": rightShift,
}

var mathUtilFuncs = dictionary{
	"abs":       abs,
	"dec":       decimal,
	"frexp":     frexp,
	"gamma":     gamma,
	"hex":       hex,
	"hypot":     hypot,
	"isInf":     isInfinity,
	"isNaN":     isNaN,
	"ldexp":     ldexp,
	"lgamma":    lgamma,
	"nextAfter": nextAfter,
	"signBit":   signBit,
	"sqrt":      sqrt,
	"to":        to,
	"until":     until,
}

var mathFuncsAliases = aliases{
	"abs":    {"absolute"},
	"acos":   {"arcCosine", "arcCosinus"},
	"acosh":  {"arcHyperbolicCosine", "arcHyperbolicCosinus"},
	"add":    {"sum"},
	"asin":   {"arcSine", "arcSinus"},
	"asinh":  {"arcHyperbolicSine", "arcHyperbolicSinus"},
	"atan":   {"arcTangent"},
	"atan2":  {"arcTangent2"},
	"atanh":  {"arcHyperbolicTangent"},
	"avg":    {"average"},
	"band":   {"bitwiseAND"},
	"bclear": {"bitwiseClear"},
	"bor":    {"bitwiseOR"},
	"bxor":   {"bitwiseXOR"},
	"ceil":   {"roundUp", "roundup"},
	"cos":    {"cosine", "cosinus"},
	"cosh":   {"hyperbolicCosine", "hyperbolicCosinus"},
	"dec":    {"decimal"},
	"deg":    {"degree"},
	"div":    {"divide", "quotient"},
	"exp":    {"exponent"},
	"exp2":   {"exponent2"},
	"floor":  {"roundDown", "rounddown", "int", "integer"},
	"hex":    {"hexa", "hexaDecimal"},
	"hypot":  {"hypotenuse"},
	"isInf":  {"isInfinity"},
	"j0":     {"firstBessel0"},
	"j1":     {"firstBessel1"},
	"jn":     {"firstBesselN"},
	"lshift": {"leftShift"},
	"max":    {"maximum", "biggest"},
	"min":    {"minimum", "smallest"},
	"mod":    {"modulo"},
	"mul":    {"multiply", "prod", "product"},
	"pow":    {"power"},
	"pow10":  {"power10"},
	"rad":    {"radian"},
	"rem":    {"remainder"},
	"rshift": {"rightShift"},
	"sin":    {"sine", "sinus"},
	"sincos": {"sineCosine", "sinusCosinus"},
	"sinh":   {"hyperbolicSine", "hyperbolicSinus"},
	"sqrt":   {"squareRoot"},
	"sub":    {"subtract"},
	"tan":    {"tangent"},
	"tanh":   {"hyperbolicTangent"},
	"trunc":  {"truncate"},
	"y0":     {"secondBessel0"},
	"y1":     {"secondBessel1"},
	"yn":     {"secondBesselN"},
}

var mathFuncsArgs = arguments{
	"abs":             {"x"},
	"acos":            {"x"},
	"acosh":           {"x"},
	"asin":            {"x"},
	"asinh":           {"x"},
	"atan":            {"x"},
	"atan2":           {"x", "y"},
	"atanh":           {"x"},
	"cbrt":            {"x"},
	"ceil":            {"x"},
	"copysign":        {"x", "y"},
	"cos":             {"x"},
	"cosh":            {"x"},
	"dim":             {"x", "y"},
	"erf":             {"x"},
	"erfc":            {"x"},
	"exp":             {"x"},
	"exp2":            {"x"},
	"expm1":           {"x"},
	"float32bits":     {"f"},
	"float32frombits": {"b"},
	"float64bits":     {"f"},
	"float64frombits": {"b"},
	"floor":           {"x"},
	"frexp":           {"f"},
	"gamma":           {"x"},
	"hypot":           {"p", "q"},
	"ilogb":           {"x"},
	"inf":             {"sign"},
	"isInf":           {"f"},
	"isNaN":           {"f"},
	"j0":              {"x"},
	"j1":              {"x"},
	"jn":              {"n", "x"},
	"ldexp":           {"frac", "exp"},
	"lgamma":          {"x"},
	"log":             {"x"},
	"log10":           {"x"},
	"log1p":           {"x"},
	"log2":            {"x"},
	"logb":            {"x"},
	"max":             {"x", "y"},
	"min":             {"x", "y"},
	"mod":             {"x", "y"},
	"modf":            {"f"},
	"nextafter":       {"x", "y"},
	"nextafter32":     {"x", "y"},
	"pow":             {"x", "y"},
	"pow10":           {"n"},
	"remainder":       {"x", "y"},
	"signbit":         {"x"},
	"sin":             {"x"},
	"sincos":          {"x"},
	"sinh":            {"x"},
	"sqrt":            {"x"},
	"tan":             {"x"},
	"tanh":            {"x"},
	"trunc":           {"x"},
	"y0":              {"x"},
	"y1":              {"x"},
	"yn":              {"n", "x"},
}

var mathFuncsHelp = descriptions{
	"abs":             "Returns the absolute value of x.\nSpecial cases are:\n    abs(±Inf) = +Inf\n    abs(NaN) = NaN",
	"acos":            "Returns the arccosine, in radians, of x.\nSpecial case is:\n    acos(x) = NaN if x < -1 or x > 1",
	"acosh":           "Returns the inverse hyperbolic cosine of x.\nSpecial cases are:\n    acosh(+Inf) = +Inf\n    acosh(x) = NaN if x < 1\n    acosh(NaN) = NaN",
	"asin":            "Returns the arcsine, in radians, of x.\nSpecial cases are:\n    asin(±0) = ±0\n    asin(x) = NaN if x < -1 or x > 1",
	"asinh":           "Returns the inverse hyperbolic sine of x.\nSpecial cases are:\n    asinh(±0) = ±0\n    asinh(±Inf) = ±Inf\n    asinh(NaN) = NaN",
	"atan":            "Returns the arctangent, in radians, of x.\nSpecial cases are:\n    atan(±0) = ±0\n    atan(±Inf) = ±Pi/2",
	"atan2":           "Returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.\nSpecial cases are (in order):\n    atan2(y, NaN) = NaN\n    atan2(NaN, x) = NaN\n    atan2(+0, x>=0) = +0\n    atan2(-0, x>=0) = -0\n    atan2(+0, x<=-0) = +Pi\n    atan2(-0, x<=-0) = -Pi\n    atan2(y>0, 0) = +Pi/2\n    atan2(y<0, 0) = -Pi/2\n    atan2(+Inf, +Inf) = +Pi/4\n    atan2(-Inf, +Inf) = -Pi/4\n    atan2(+Inf, -Inf) = 3Pi/4\n    atan2(-Inf, -Inf) = -3Pi/4\n    atan2(y, +Inf) = 0\n    atan2(y>0, -Inf) = +Pi\n    atan2(y<0, -Inf) = -Pi\n    atan2(+Inf, x) = +Pi/2\n    atan2(-Inf, x) = -Pi/2",
	"atanh":           "Returns the inverse hyperbolic tangent of x.\nSpecial cases are:\n    atanh(1) = +Inf\n    atanh(±0) = ±0\n    atanh(-1) = -Inf\n    atanh(x) = NaN if x < -1 or x > 1\n    atanh(NaN) = NaN",
	"cbrt":            "Returns the cube root of x.\nSpecial cases are:\n    cbrt(±0) = ±0\n    cbrt(±Inf) = ±Inf\n    cbrt(NaN) = NaN",
	"ceil":            "Returns the least integer value greater than or equal to x.\nSpecial cases are:\n    ceil(±0) = ±0\n    ceil(±Inf) = ±Inf\n    ceil(NaN) = NaN",
	"copysign":        "Returns a value with the magnitude of x and the sign of y",
	"cos":             "Returns the cosine of the radian argument x.\nSpecial cases are:\n    cos(±Inf) = NaN\n    cos(NaN) = NaN",
	"cosh":            "Returns the hyperbolic cosine of x.\nSpecial cases are:\n    cosh(±0) = 1\n    cosh(±Inf) = +Inf\n    cosh(NaN) = NaN",
	"dim":             "Returns the maximum of x-y or 0.\nSpecial cases are:\n    dim(+Inf, +Inf) = NaN\n    dim(-Inf, -Inf) = NaN\n    dim(x, NaN) = dim(NaN, x) = NaN",
	"erf":             "Returns the error function of x.\nSpecial cases are:\n    Erf(+Inf) = 1\nErf(-Inf) = -1\nErf(NaN) = NaN",
	"erfc":            "Returns the complementary error function of x.\nSpecial cases are:\n    Erfc(+Inf) = 0\nErfc(-Inf) = 2\nErfc(NaN) = NaN",
	"exp":             "Returns e**x, the base-e exponential of x.\nSpecial cases are:\n    exp(+Inf) = +Inf\n    exp(NaN) = NaN\nVery large values overflow to 0 or +Inf. Very small values underflow to 1.",
	"exp2":            "Returns 2**x, the base-2 exponential of x.\nSpecial cases are the same as exp.",
	"expm1":           "Returns e**x - 1, the base-e exponential of x minus 1. It is more\naccurate than exp(x) - 1 when x is near zero.\nSpecial cases are:\n    expm1(+Inf) = +Inf\n    expm1(-Inf) = -1\n    expm1(NaN) = NaN\nVery large values overflow to -1 or +Inf",
	"float32bits":     "Returns the IEEE 754 binary representation of f",
	"float32frombits": "Returns the floating point number corresponding to the\nIEEE 754 binary representation b",
	"float64bits":     "Returns the IEEE 754 binary representation of f",
	"float64frombits": "Returns the floating point number corresponding the IEEE\n754 binary representation b",
	"floor":           "Returns the greatest integer value less than or equal to x.\nSpecial cases are:\n    floor(±0) = ±0\n    floor(±Inf) = ±Inf\n    floor(NaN) = NaN",
	"frexp":           "Breaks f into a normalized fraction and an integral power of two. Returns frac and exp satisfying f == frac × 2**exp, with the absolute value of frac in the interval [½, 1).\nSpecial cases are:\n    frexp(±0) = ±0, 0\n    frexp(±Inf) = ±Inf, 0\n    frexp(NaN) = NaN, 0",
	"gamma":           "Returns the Gamma function of x.\nSpecial cases are:\n    gamma(+Inf) = +Inf\n    gamma(+0) = +Inf\n    gamma(-0) = -Inf\n    gamma(x) = NaN for integer x < 0\n    gamma(-Inf) = NaN\n    gamma(NaN) = NaN",
	"hypot":           "Returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow.\nSpecial cases are:\n    hypot(±Inf, q) = +Inf\n    hypot(p, ±Inf) = +Inf\n    hypot(NaN, q) = NaN\n    hypot(p, NaN) = NaN",
	"ilogb":           "Returns the binary exponent of x as an integer.\nSpecial cases are:\n    ilogb(±Inf) = MaxInt32\n    ilogb(0) = MinInt32\n    ilogb(NaN) = MaxInt32",
	"inf":             "Returns positive infinity if sign >= 0, negative infinity if sign <\n0",
	"isInf":           "Reports whether f is an infinity, according to sign. If sign > 0, isInf reports whether f is positive infinity. If sign < 0, IsInf reports whether f is negative infinity. If sign == 0, IsInf reports whether f is either infinity",
	"isNaN":           "Reports whether f is an IEEE 754 'not-a-number' value",
	"j0":              "Returns the order-zero Bessel function of the first kind.\nSpecial cases are:\n    j0(±Inf) = 0\n    j0(0) = 1\n    j0(NaN) = NaN",
	"j1":              "Returns the order-one Bessel function of the first kind.\nSpecial cases are:\n    j1(±Inf) = 0\n    j1(NaN) = NaN",
	"jn":              "Returns the order-n Bessel function of the first kind.\nSpecial cases are:\n    jn(n, ±Inf) = 0\n    jn(n, NaN) = NaN",
	"ldexp":           "Ldexp is the inverse of Frexp. Returns frac × 2**exp.\nSpecial cases are:\n    ldexp(±0, exp) = ±0\n    ldexp(±Inf, exp) = ±Inf\n    ldexp(NaN, exp) = NaN",
	"lgamma":          "Returns the natural logarithm and sign (-1 or +1) of Gamma(x).\nSpecial cases are:\n    lgamma(+Inf) = +Inf\n    lgamma(0) = +Inf\n    lgamma(-integer) = +Inf\n    lgamma(-Inf) = -Inf\n    lgamma(NaN) = NaN",
	"log":             "Returns the natural logarithm of x.\nSpecial cases are:\n    log(+Inf) = +Inf\n    log(0) = -Inf\n    log(x < 0) = NaN\n    log(NaN) = NaN",
	"log10":           "Returns the decimal logarithm of x. The special cases are the same as for log.",
	"log1p":           "Returns the natural logarithm of 1 plus its argument x. It is more accurate than log(1 + x) when x is near zero.\nSpecial cases are:\n    log1p(+Inf) = +Inf\n    log1p(±0) = ±0\n    log1p(-1) = -Inf\n    log1p(x < -1) = NaN\n    log1p(NaN) = NaN",
	"log2":            "Returns the binary logarithm of x. The special cases are the same as for log.",
	"logb":            "Returns the binary exponent of x.\nSpecial cases are:\n    logb(±Inf) = +Inf\n    logb(0) = -Inf\n    logb(NaN) = NaN",
	"max":             "Returns the larger of x or y.\nSpecial cases are:\n    max(x, +Inf) = max(+Inf, x) = +Inf\n    max(x, NaN) = max(NaN, x) = NaN\n    max(+0, ±0) = max(±0, +0) = +0\n    max(-0, -0) = -0",
	"min":             "Returns the smaller of x or y.\nSpecial cases are:\n    min(x, -Inf) = min(-Inf, x) = -Inf\n    min(x, NaN) = min(NaN, x) = NaN\n    min(-0, ±0) = min(±0, -0) = -0",
	"mod":             "Returns the floating-point remainder of x/y. The magnitude of the result is less than y and its sign agrees with that of x.\nSpecial cases are:\n    mod(±Inf, y) = NaN\n    mod(NaN, y) = NaN\n    mod(x, 0) = NaN\n    mod(x, ±Inf) = x\n    mod(x, NaN) = NaN",
	"modf":            "Returns integer and fractional floating-point numbers that sum to f. Both values have the same sign as f.\nSpecial cases are:\n    modf(±Inf) = ±Inf, NaN\n    modf(NaN) = NaN, NaN",
	"naN":             "Returns an IEEE 754 'not-a-number' value.",
	"nextAfter":       "Returns the next representable float64 value after x towards y.\nSpecial cases are:\n    Nextafter(x, x)   = x\nNextafter(NaN, y) = NaN\nNextafter(x, NaN) = NaN",
	"nextAfter32":     "Returns the next representable float32 value after x towards y.\nSpecial cases are:\n    Nextafter32(x, x)   = x\nNextafter32(NaN, y) = NaN\nNextafter32(x, NaN) = NaN",
	"pow":             "Returns x**y, the base-x exponential of y.\nSpecial cases are (in order):\n    pow(x, ±0) = 1 for any x\n    pow(1, y) = 1 for any y\n    pow(x, 1) = x for any x\n    pow(NaN, y) = NaN\n    pow(x, NaN) = NaN\n    pow(±0, y) = ±Inf for y an odd integer < 0\n    pow(±0, -Inf) = +Inf\n    pow(±0, +Inf) = +0\n    pow(±0, y) = +Inf for finite y < 0 and not an odd integer\n    pow(±0, y) = ±0 for y an odd integer > 0\n    pow(±0, y) = +0 for finite y > 0 and not an odd integer\n    pow(-1, ±Inf) = 1\n    pow(x, +Inf) = +Inf for |x| > 1\n    pow(x, -Inf) = +0 for |x| > 1\n    pow(x, +Inf) = +0 for |x| < 1\n    pow(x, -Inf) = +Inf for |x| < 1\n    pow(+Inf, y) = +Inf for y > 0\n    pow(+Inf, y) = +0 for y < 0\n    pow(-Inf, y) = Pow(-0, -y)\n    pow(x, y) = NaN for finite x < 0 and finite non-integer y",
	"pow10":           "Returns 10**n, the base-10 exponential of n.\nSpecial cases are:\n    pow10(n) =0 for n < -323\n    pow10(n) = +Inf for n > 308",
	"rem":             "Returns the IEEE 754 floating-point remainder of x/y.\nSpecial cases are:\n    rem(±Inf, y) = NaN\n    rem(NaN, y) = NaN\n    rem(x, 0) = NaN\n    rem(x, ±Inf) = x\n    rem(x, NaN) = NaN",
	"signbit":         "Returns true if x is negative or negative zero",
	"sin":             "Returns the sine of the radian argument x.\nSpecial cases are:\n    sin(±0) = ±0\n    sin(±Inf) = NaN\n    sin(NaN) = NaN",
	"sincos":          "Returns Sin(x), Cos(x).\nSpecial cases are:\n    sincos(±0) = ±0, 1\n    sincos(±Inf) = NaN, NaN\n    sincos(NaN) = NaN, NaN",
	"sinh":            "Returns the hyperbolic sine of x.\nSpecial cases are:\n    sinh(±0) = ±0\n    sinh(±Inf) = ±Inf\n    sinh(NaN) = NaN",
	"sqrt":            "Returns the square root of x.\nSpecial cases are:\n    sqrt(+Inf) = +Inf\n    sqrt(±0) = ±0\n    sqrt(x < 0) = NaN\n    sqrt(NaN) = NaN",
	"tan":             "Returns the tangent of the radian argument x.\nSpecial cases are:\n    tan(±0) = ±0\n    tan(±Inf) = NaN\n    tan(NaN) = NaN",
	"tanh":            "Returns the hyperbolic tangent of x.\nSpecial cases are:\n    tanh(±0) = ±0\n    tanh(±Inf) = ±1\n    tanh(NaN) = NaN",
	"trunc":           "Returns the integer value of x.\nSpecial cases are:\n    trunc(±0) = ±0\n    trunc(±Inf) = ±Inf\n    trunc(NaN) = NaN",
	"y0":              "Returns the order-zero Bessel function of the second kind.\nSpecial cases are:\n    y0(+Inf) = 0\n    y0(0) = -Inf\n    y0(x < 0) = NaN\n    y0(NaN) = NaN",
	"y1":              "Returns the order-one Bessel function of the second kind.\nSpecial cases are:\n    y1(+Inf) = 0\n    y1(0) = -Inf\n    y1(x < 0) = NaN\n    y1(NaN) = NaN",
	"yn":              "Returns the order-n Bessel function of the second kind.\nSpecial cases are:\n    yn(n, +Inf) = 0\n    yn(n ≥ 0, 0) = -Inf\n    yn(n < 0, 0) = +Inf if n is odd, -Inf if n is even\n    yn(n, x < 0) = NaN\n    yn(n, NaN) = NaN",
}

func (t *Template) addMathFuncs() {
	// Enhance mathematic functions
	options := funcOptions{
		funcHelp:    mathFuncsHelp,
		funcArgs:    mathFuncsArgs,
		funcAliases: mathFuncsAliases,
	}

	t.AddFunctions(mathBaseFuncs, mathBase, options)
	t.AddFunctions(mathStatFuncs, mathStatistics, options)
	t.AddFunctions(mathTrigFuncs, mathTrigonometry, options)
	t.AddFunctions(mathBitsFuncs, mathBits, options)
	t.AddFunctions(mathUtilFuncs, mathUtilities, options)

	constants := dictionary{
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
		"MaxUint8":               math.MaxUint8,
		"MaxUint16":              math.MaxUint16,
		"MaxUint32":              math.MaxUint32,
		// Those values are commented because they causes problem with object serialization.
		// "MaxInt64":            math.MaxInt64,
		// "MaxUint64":           uint(math.MaxUint64),
		// "Nan":                 math.NaN(),
		// "Infinity":            math.Inf(1),
		// "Inf":                 math.Inf(1),
		// "NegativeInfinity":    math.Inf(-1),
		// "NegInf":              math.Inf(-1),
	}

	// We do not want to inject the math constant twice
	if !t.optionsEnabled[Math] {
		t.setConstant(true, constants, "Math", "MATH")
		t.optionsEnabled[Math] = true
	}
}

func to(params ...interface{}) (interface{}, error)    { return generateNumericArray(true, params...) }
func until(params ...interface{}) (interface{}, error) { return generateNumericArray(false, params...) }

func abs(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Abs)
}

func cbrt(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Cbrt)
}

func ceil(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Ceil)
}

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

func j0(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.J0)
}

func j1(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.J1)
}

func jn(n, x interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return math.Jn(toInt(n), toFloat(x)), nil
}

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

func y0(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Y0)
}

func y1(a interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return processFloat(a, math.Y1)
}

func yn(n, x interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return math.Yn(toInt(n), toFloat(x)), nil
}
