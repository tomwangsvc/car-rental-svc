package time

import (
	"strings"
	"time"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_math "github.com/tomwangsvc/lib-svc/math"
)

const (
	LayoutBoxer         = "02/01/2006"
	LayoutLocal         = "2006-01-02T15:04:05.999999999"
	DefaultTimezoneName = "Antarctica/McMurdo"
)

// NormalizeTimeAndTruncateToMilliseconds normalises time to UTC and truncates to milliseconds
// NOTE: This used to be NormalizeTime, however that did not describe the functionality
// 		If you don't want to truncate milliseconds, just use the UTC() function
func NormalizeTimeAndTruncateToMilliseconds(t time.Time) time.Time {
	return t.UTC().Truncate(time.Millisecond)
}

// NormalizeTimeAndTruncateToMillisecondsFormatted normalises time to UTC and truncates to milliseconds and formats as string
// NOTE: This used to be NormalizeTime, however that did not describe the functionality
// 		If you don't want to truncate milliseconds, just use the UTC() function
func NormalizeTimeAndTruncateToMillisecondsFormatted(t time.Time) string {
	return NormalizeTimeAndTruncateToMilliseconds(t).Format(time.RFC3339)
}

// NormalizeTimeToStartOfDay normalises time to the start of the day then to UTC
func NormalizeTimeToStartOfDay(t time.Time) time.Time {
	return t.UTC().Truncate(time.Hour * 24)
}

// NormalizeTimeToEndOfDay normalises time to the end of the day then to UTC
func NormalizeTimeToEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
}

// ParseFormattedTime parses and normalises an RFC3339 formatted time string
func ParseFormattedTime(ft string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(ft))
	if err != nil {
		return nil, lib_errors.Wrapf(err, "Failed parsing time %v as RFC3339", t)
	}
	return &t, nil
}

// ParseFormattedTimeWithFullPrecision parses a RFC3339Nano formatted time string
func ParseFormattedTimeWithFullPrecision(ft string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(ft))
	if err != nil {
		return nil, lib_errors.Wrapf(err, "Failed parsing time %v as RFC3339Nano", t)
	}
	return &t, nil
}

func ParseLocalTime(localTime, locationName string, ignoreDst bool) (*time.Time, error) {
	location, err := time.LoadLocation(locationName)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed loading location")
	}
	t, err := time.ParseInLocation(LayoutLocal, localTime, location)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed parsing string to time")
	}
	if ignoreDst {
		_, smallestOffset := t.Zone()
		for i := 1; i < 12; i++ {
			if _, offset := t.AddDate(0, i, 0).Zone(); offset < smallestOffset {
				smallestOffset = offset
			}
		}
		t, err = time.ParseInLocation(LayoutLocal, localTime, time.FixedZone("", smallestOffset))
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed parsing string to time")
		}
	}

	return &t, nil
}

// IsAfterToday returns true if the time passed is greater than or equal to start of tomorrow in UTC
// -> Tomorrow or some other day in the future
func IsAfterToday(t time.Time) bool {
	now := time.Now()
	startOfTomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return t.Equal(startOfTomorrow) || t.After(startOfTomorrow)
}

// TodayIsAfter returns true if today is greater than or equal to start of tomorrow from the time passed
// -> Time passed is yesterday or some other day in the past
func TodayIsAfter(t time.Time) bool {
	startOfTomorrow := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
	now := time.Now()
	return now.Equal(startOfTomorrow) || now.After(startOfTomorrow)
}

// IsBeforeSomeDayInPast returns true if the time given is more than some number of days in the past
func IsBeforeSomeDayInPast(t time.Time, days int) bool {
	now := time.Now()
	startOfSomeDayInPast := time.Date(now.Year(), now.Month(), now.Day()-days, 0, 0, 0, 0, now.Location())
	return t.Equal(startOfSomeDayInPast) || t.Before(startOfSomeDayInPast)
}

// IsFirstTimeAfterSecondTime returns true if the first time is after the second time using second precision
// -> firstTime > secondTime
func IsFirstTimeAfterSecondTime(firstTime time.Time, secondTime time.Time) bool {
	return IsFirstTimeAfterSecondTimePlusDays(firstTime, secondTime, 0)
}

// IsFirstTimeAfterSecondTimePlusDays returns true if the first time is after the second time plus some number of days using second precision
// -> firstTime > secondTime + days
func IsFirstTimeAfterSecondTimePlusDays(firstTime time.Time, secondTime time.Time, days int) bool {
	secondTime = secondTime.AddDate(0, 0, days)
	return firstTime.After(secondTime)
}

// IsFirstTimeBeforeSecondTime returns true if the first time is after the second time using second precision
// -> firstTime < secondTime
func IsFirstTimeBeforeSecondTime(firstTime time.Time, secondTime time.Time) bool {
	return IsFirstTimeBeforeSecondTimePlusDays(firstTime, secondTime, 0)
}

// IsFirstTimeBeforeSecondTimePlusDays returns true if the first time is after the second time plus some number of days using second precision
// -> firstTime < secondTime + days
func IsFirstTimeBeforeSecondTimePlusDays(firstTime time.Time, secondTime time.Time, days int) bool {
	secondTime = secondTime.AddDate(0, 0, days)
	return firstTime.Before(secondTime)
}

// Diff returns the difference between two dates
// -> https://stackoverflow.com/questions/36530251/golang-time-since-with-months-and-years/36531443#36531443
func Diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = y2 - y1
	month = int(M2 - M1)
	day = d2 - d1
	hour = h2 - h1
	min = m2 - m1
	sec = s2 - s1

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func DaysElapsedAbsolute(a, b time.Time) int {
	if a.Before(b) {
		a, b = b, a
	}
	daysFloat := a.Sub(b).Seconds()
	days := int(daysFloat / 86400)
	if daysFloat/1 >= 0 {
		days++
	}
	return days
}

func MonthsElapsedAbsolute(a, b time.Time) int {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	months := int(M2 - M1)
	days := d2 - d1

	// Days can be negative, meaning we got exactly the month value we needed
	// However, if days are non-negative, it means we're into a new month and need to round up
	if days >= 0 {
		months++
	}

	years := y2 - y1
	months += years * 12

	return months
}

func ParseSecondsIntoSecondsAndNanoseconds(value float64) (seconds int64, nanoseconds int32, negative bool, err error) {
	var milliseconds int32
	seconds, milliseconds, negative, err = lib_math.SplitFloatIntoWholeAndFraction(value)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed splitting value into whole and fraction")
		return
	}
	nanoseconds = milliseconds * 1000
	return
}
