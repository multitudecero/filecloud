package main

import (
	"bytes"
	"filecloud/common/dl/models"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// отправляем

	data, err := os.ReadFile("./example.txt")
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest(http.MethodPut, "http://localhost:8080", bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add(models.FileNameHeader, "example.txt")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode >= 300 {
		content, _ := io.ReadAll(response.Body)
		log.Fatalf("status: %+v, text: %+v", response.StatusCode, string(content))
	}

	// получаем

	request, err = http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add(models.FileNameHeader, "example.txt")

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	content, _ := io.ReadAll(response.Body)

	if response.StatusCode >= 300 {
		log.Fatalf("status: %+v, text: %+v", response.StatusCode, content)
	}

	log.Println(string(content))
}
