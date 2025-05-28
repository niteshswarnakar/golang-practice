package gin_server

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/niteshswarnakar/my-go-library/params"
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
	param := c.Request.URL.Query()
	fmt.Println("\nPREVIOUS PARAMS : ", param)
	myparams := params.GetQueryParams(param)
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
