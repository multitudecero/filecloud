package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	id := os.Args[1]
	port := os.Args[2]
	folder := os.Args[3]
	outputFolder := os.Args[4]
	templateFolder := os.Args[5]

	data, err := os.ReadFile(templateFolder + "/app_template.cfg")
	if err != nil {
		log.Fatal(err)
	}

	content := string(data)
	content = strings.ReplaceAll(content, "$$ID", id)
	content = strings.ReplaceAll(content, "$$PORT", port)
	content = strings.ReplaceAll(content, "$$FOLDER", folder)

	err = os.WriteFile(outputFolder+"/app.cfg", []byte(content), 0777)
	if err != nil {
		log.Fatal(err)
	}
}
