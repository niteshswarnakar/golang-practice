package test

import "fmt"

var unicodeMap = map[string]int{
	"1": 49,
	"2": 50,
	"3": 51,
	"4": 52,
	"5": 53,
	"6": 54,
	"7": 55,
	"8": 56,
	"9": 57,
	"0": 48,
}

func GetUnicodeCodePoint(char string) int {
	return unicodeMap[char] - 48
}

func GetNumber(num string) int {
	var number int = 0
	for _, item := range num {
		number = number * 10
		fmt.Println("number : ", number)
		number += GetUnicodeCodePoint(string(item))
	}
	return number
}
