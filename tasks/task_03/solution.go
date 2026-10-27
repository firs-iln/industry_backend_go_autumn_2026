package main

import (
	"errors"
	"strconv"
)

var ErrZero = errors.New("zero is not allowed")

func fizzBuzz(n int) (string, error) {
	if n == 0 {
		return "", ErrZero
	}
	// Remainders of a negative n are negative or zero, so compare with 0
	// instead of taking the absolute value: it overflows for math.MinInt.
	switch {
	case n%15 == 0:
		return "FizzBuzz", nil
	case n%3 == 0:
		return "Fizz", nil
	case n%5 == 0:
		return "Buzz", nil
	}
	return strconv.Itoa(n), nil
}
