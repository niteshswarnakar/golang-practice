package test

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
	return unicodeMap[char]
}

func GetNumber(num string) int {
	var number int = 0
	for i := range num {
		number += GetUnicodeCodePoint(string(num[len(num)-i]))
	}
	return number
}
