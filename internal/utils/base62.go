package utils

import "strings"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Encode(number int64) string {
	if number == 0 {
		return string(alphabet[0])
	}
	var result strings.Builder

	for number > 0 {
		result.WriteByte(alphabet[number%62])
		number = number / 62
	}
	return result.String()
}
