package comm

import (
	"encoding/json"
	"log"
	"net/http"
)

func SetResponse(response http.ResponseWriter, statusCode int, body interface{}, err error) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	if err != nil {
		log.Println(err.Error())
		_, err := response.Write([]byte(err.Error()))
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("ошибка возврата ошибки"))
		}
	} else if body != nil {
		data, err := json.MarshalIndent(body, "", "	")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("ошибка сериализации ответа: " + err.Error()))
		}
		response.Write(data)
	}
}
