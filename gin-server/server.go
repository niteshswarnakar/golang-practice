package gin_server

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Server() {
	server := gin.Default()
	server.GET("/", handler)
	println("Server started at port 5000")
	if err := server.Run(":5000"); err != nil {
		panic(err)
	}
}

func handler(c *gin.Context) {
	params := c.Request.URL.Query()
	fmt.Println("\nPREVIOUS PARAMS : ", params)
	myparams := GetQueryParams(params)
	fmt.Printf("\nAFTER PARAMS : %s\n\n", myparams)
	for key, item := range myparams {
		value := item.Value
		fmt.Printf("%s : %v\n", key, value)
		fmt.Println("Type : ", fmt.Sprintf("%T", value))
		fmt.Println("Reflect Type : ", reflect.TypeOf(value).Kind())
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			fmt.Println("Slice Value : ", value)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	return
}

type MyAny struct {
	Value any
}
type MyQueryParams map[string]MyAny

func GetQueryParams(params url.Values) MyQueryParams {
	result := make(MyQueryParams)
	for key, value := range params {
		if len(value) > 1 {
			result[key] = MyAny{Value: value}
		} else if len(value) == 1 {
			val := value[0]
			if val == "true" || val == "1" {
				result[key] = MyAny{Value: true}
			} else if val == "false" || val == "0" {
				result[key] = MyAny{Value: false}
			} else {
				number, err := strconv.Atoi(val)
				if err != nil {
					fmt.Println("Error converting to int:", err)
					result[key] = MyAny{Value: val}
				} else {
					result[key] = MyAny{Value: number}
				}
			}
		} else {
			result[key] = MyAny{Value: nil}
		}
	}
	return result
}
