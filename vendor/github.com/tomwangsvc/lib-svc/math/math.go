package math

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

// Min returns minimum value of two int64
func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Max returns maximum value of two int64
func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// Round rounds numbers to arbitrary precision
func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// Epsilon is the error margin for float comparison
const Epsilon float64 = 0.00000001

// Float64Equals calls Float64EqualsWithinPrecision requiring high precision for the calculation
func Float64Equals(a, b float64) bool {
	return float64EqualsWithinPrecision(a, b, Epsilon)
}

// Float64EqualsWithinPrecision does not check whether the numbers are exactly the same since this is not possible, but instead whether their difference is very small
func float64EqualsWithinPrecision(a, b, eps float64) bool {
	return (a-b) < eps && (b-a) < eps
}

func SplitFloatIntoWholeAndFraction(value float64) (whole int64, fraction int32, negative bool, err error) {
	if value < 0 {
		negative = true
	}

	parts := strings.Split(fmt.Sprintf("%f", value), ".")
	if whole, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
		err = lib_errors.Wrap(err, "Failed parsing whole part of value into int64")
		return
	}

	var tmp int64
	if tmp, err = strconv.ParseInt(parts[1], 10, 32); err != nil {
		err = lib_errors.Wrap(err, "Failed parsing fraction part of value into int32")
		return
	}
	fraction = int32(tmp)

	return
}
