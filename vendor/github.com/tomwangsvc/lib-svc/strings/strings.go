package strings

import "fmt"

func Contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func Remove(src, inputs []string) []string {
	var d []string
	for _, v := range src {
		var found bool
		for _, input := range inputs {
			if input == v {
				found = true
			}
		}

		if !found {
			d = append(d, v)
		}
	}

	return d
}

func Overlaps(slice1, slice2 []string) bool {
	for _, m := range slice1 {
		for _, n := range slice2 {
			if m == n {
				return true
			}
		}
	}

	return false
}

func PrependString(src []string, input string) []string {
	var d []string
	for _, v := range src {
		d = append(d, fmt.Sprintf("%s%s", input, v))
	}

	return d
}

func AppendString(src []string, input string) []string {
	var d []string
	for _, v := range src {
		d = append(d, fmt.Sprintf("%s%s", v, input))
	}

	return d
}
