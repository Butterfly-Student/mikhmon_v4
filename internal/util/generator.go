package util

import (
	"math/rand"
)

const (
	digits = "0123456789"
	lower  = "abcdefghijklmnopqrstuvwxyz"
	upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandN returns a random string of digits of length n.
func RandN(length int) string { return randFrom(digits, length) }

// RandLC returns a random lowercase string of length n.
func RandLC(length int) string { return randFrom(lower, length) }

// RandUC returns a random uppercase string of length n.
func RandUC(length int) string { return randFrom(upper, length) }

// RandULC returns a random string of uppercase and lowercase letters.
func RandULC(length int) string { return randFrom(upper+lower, length) }

// RandNLC returns a random string of digits and lowercase letters.
func RandNLC(length int) string { return randFrom(digits+lower, length) }

// RandNUC returns a random string of digits and uppercase letters.
func RandNUC(length int) string { return randFrom(digits+upper, length) }

// RandNULC returns a random string of digits, uppercase, and lowercase letters.
func RandNULC(length int) string { return randFrom(digits+upper+lower, length) }

// randFrom picks `length` characters randomly from charset using Fisher-Yates shuffle,
// matching the PHP str_shuffle() + substr() behavior.
func randFrom(charset string, length int) string {
	chars := []byte(charset)
	rand.Shuffle(len(chars), func(i, j int) { chars[i], chars[j] = chars[j], chars[i] })
	if length > len(chars) {
		length = len(chars)
	}
	return string(chars[:length])
}
