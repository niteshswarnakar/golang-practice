package main

import (
	"fmt"

	"github.com/niteshswarnakar/golang-practice/pdf_read"
)

func main() {
	filePath := "/home/iamnitesh/Desktop/work_space/golang-practice/files/campain.pdf"
	// content, err := pdf_read.ReadTextFile(filePath)
	content, err := pdf_read.PdfRead(filePath)
	if err != nil {
		panic(err)
	}

	// Print the content of the PDF
	fmt.Printf("Content of the PDF:\n%s\n", content)
}
