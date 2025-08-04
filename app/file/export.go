package file

import (
	"fmt"
	"os"
)

func InitDB(dbPath string, dbFile string) {
	filePath := fmt.Sprintf("./%s/%s", dbPath, dbFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("No file exists ata the provided path: %s\nInit with empty database\n", filePath)
		return
	}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading database gile from disk: %s\n", err)
		return
	}

	err = parseFile(fileContent)
	if err != nil {
		fmt.Printf("Error parsing database file: %s\n", err)
	}
}
