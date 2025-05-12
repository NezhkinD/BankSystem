package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateLuhnNumber генерирует номер карты заданной длины, валидный по алгоритму Луна
func GenerateLuhnNumber(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	number := make([]int, length-1)
	for i := range number {
		number[i] = r.Intn(10)
	}

	sum := 0
	double := true
	for i := len(number) - 1; i >= 0; i-- {
		digit := number[i]
		if double {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}
		sum += digit
		double = !double
	}

	lastDigit := (10 - (sum % 10)) % 10
	number = append(number, lastDigit)

	var s string
	for _, d := range number {
		s += strconv.Itoa(d)
	}
	return s
}
