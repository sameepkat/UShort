package encoding

import (
	"errors"
	"strings"
)

const (
	base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base        = 62
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrOverflow     = errors.New("integer overflow")
)

func Encode(num uint64) string {
	if num == 0 {
		return string(base62chars[0])
	}
	var builder strings.Builder
	builder.Grow(11)

	for num > 0 {
		remainder := num % base
		builder.WriteByte(base62chars[remainder])
		num /= base
	}
	result := builder.String()
	return reverseString(result)
}

func Decode(str string) (uint64, error) {
	if str == "" {
		return 0, ErrInvalidInput
	}
	var result uint64
	for _, char := range str {
		index := -1
		for j, c := range base62chars {
			if c == char {
				index = j
				break
			}
		}
		if index == -1 {
			return 0, ErrInvalidInput
		}

		//Check for overflow
		if result > (^uint64(0)-uint64(index))/base {
			return 0, ErrOverflow
		}

		result = result*base + uint64(index)
	}
	return result, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
