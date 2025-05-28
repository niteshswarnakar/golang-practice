package ziptest

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnZip() {
	filename := "ziptest/testfile.zip"
	fmt.Println("Unzipping file:", filename)

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return
	}
	defer zipReader.Close()

	// Create a directory to extract the archive into
	extractPath := "ziptest/extracted"
	os.MkdirAll(extractPath, os.ModePerm)

	// Read all the files from zip archive
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer rc.Close()
		path := filepath.Join(extractPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		outFile, err := os.Create(path)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}
		defer outFile.Close()

		// Copy contents to file
		_, err = io.Copy(outFile, rc)
		if err != nil {
			fmt.Println("Error extracting file:", err)
			continue
		}

		fmt.Printf("Extracted %s successfully\n", file.Name)
	}
}
