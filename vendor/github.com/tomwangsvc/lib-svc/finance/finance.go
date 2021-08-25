package finance

import (
	"strconv"

	lib_math "github.com/tomwangsvc/lib-svc/math"

	humanize "github.com/dustin/go-humanize"
)

// Variance is the error margin for float comparison
const Variance float64 = 0.00999999

// ValuesEqualWithinVariance checks that the two monetary values are within the constant Variance
func ValuesEqualWithinVariance(a, b float64) bool {
	return (a-b) < Variance && (b-a) < Variance
}

func FormatMoney(money float64) string {
	return humanize.FormatFloat("#,###.##", money)
}

func DisplayMoney(money float64) string {
	return strconv.FormatFloat(Round(money), 'f', 2, 64)
}

// DisplayMoneyAsFloat64 presents the value returned by DisplayMoney as a float64
func DisplayMoneyAsFloat64(money float64) float64 {
	v, _ := strconv.ParseFloat(DisplayMoney(money), 64)
	return v
}

func Round(money float64) float64 {
	return lib_math.Round(money, 0.01)
}
