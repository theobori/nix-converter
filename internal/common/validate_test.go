package common

import "testing"

var numbersOk = []string{
	"0.00001",
	"-1",
	"-1.123",
	"123",
}

var numbersKo = []string{
	"0.00.001",
	"--1",
	"-.1123",
	"123a",
	"123.",
	"",
	"ajd",
}

func TestIsNumber(t *testing.T) {
	for _, numberOk := range numbersOk {
		if !IsNumber(numberOk) {
			t.Fatalf("the string '%s' is a number", numberOk)
		}
	}
}

func TestIsNotNumber(t *testing.T) {
	for _, numberKo := range numbersKo {
		if IsNumber(numberKo) {
			t.Fatalf("the string '%s' is not a number", numberKo)
		}
	}
}
