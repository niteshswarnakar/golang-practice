package test

import (
	"fmt"
	"testing"
)

func TestFunc(t *testing.T) {
	num := GetNumber("1234567890")
	fmt.Println(num)
	t.Log(num)
}
