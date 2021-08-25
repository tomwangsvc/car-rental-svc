package util

import (
	"math/rand"
)

//Lowercase and Uppercase and number
func Random(length int) string {
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP0123456789"
	var output string
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output = output + string(randomChar)
	}
	return output
}
