package pdf_read

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/ledongthuc/pdf"
)

func readPdf(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	// f, r, err := pdf.Open(path)
	// defer f.Close()
	// if err != nil {
	// 	return "", err
	// }

	fileinfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	r, err := pdf.NewReader(file, fileinfo.Size())

	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func PdfRead(filePath string) (string, error) {
	return readPdf(filePath)
}

func TextRead() {
	text := `Name: John Doe
			Phone: +1 555-123-4567
			Address: 123 Main St, Anytown, CA 91234, USA
			Email: john.doe@example.com`

	// re := regexp.MustCompile(`Name:\s*(?P<name>[A-Za-z\s]+)\s*Phone:\s*(?P<phone>\+?\d[\d\s\-\(\)]+)\s*Address:\s*(?P<address>[\w\s\-,.]+),?\s*USA?\s*Email:\s*(?P<email>[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
	re := regexp.MustCompile(`(?s).*`)

	match := re.FindStringSubmatch(text)
	result := make(map[string]string)

	fmt.Println("\n\n")
	fmt.Println("Match found:", match)
	fmt.Println("\n\n")

	if len(match) > 0 {
		for i, name := range re.SubexpNames() {
			fmt.Println("key : ", name)
			fmt.Println("value : ", match[i])
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
	} else {
		fmt.Println("No match found.")
	}

	fmt.Println("\n\nResult map:")
	for key, _ := range result {
		fmt.Printf("%s: \n", key)
	}
}

func ReadTextFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()
	// Read the file content
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	fmt.Println("File Name: ", fileInfo.Name())
	tempBytes := &bytes.Buffer{}
	file.WriteTo(tempBytes)

	return tempBytes.String(), nil
}
