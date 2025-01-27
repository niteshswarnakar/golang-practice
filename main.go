package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/niteshswarnakar/golang-practice/email_server"
)

type Letter struct {
	number string
	count  int
}

func GetLetter() (Letter, error) {
	var letter = Letter{
		number: "123",
		count:  3,
	}

	return letter, errors.New("this is error")
}

func main() {
	// params := []interface{}{"nitesh", 25, 3}
	// TestFunction(params...)
	// testmail.SendMail()
	// test.GetNumber("45621")
	// newNumber := learn.NewNumber(word)
	// fmt.Println(newNumber.Count())

	// letter := goo.Must(GetLetter())
	// letter, _ := GetLetter()

	// workerpool.WorkerPool()

	// myschedular.GoCronRunner()

	// password.Test()

	// database.Main()

	// storage.Storage()

	email_server.EmailTest()
}

func TestFunction(params ...interface{}) {
	name := params[0].(string)
	age := params[1].(int)
	grade := params[2].(float64)
	fmt.Println("Name: ", name)
	fmt.Println("Age: ", age)
	fmt.Println("Grade: ", grade)
}

func Main() {
	wg := &sync.WaitGroup{}
	_map := map[int][]string{
		2: {"a", "b", "c"},
		3: {"d", "e", "f"},
		4: {"g", "h", "i"},
		5: {"j", "k", "l"},
		6: {"m", "n", "o"},
		7: {"p", "q", "r", "s"},
		8: {"t", "u", "v"},
		9: {"w", "x", "y", "z"},
	}
	fmt.Println(_map)

	// var sentence string

	reset := make(chan bool)
	letter := make(chan Letter)

	wg.Add(1)
	go Writer(letter, reset, wg)
	wg.Wait()

}

func Writer(letter chan Letter, reset chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("WRITER FUNCTION")
}
