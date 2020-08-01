package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"status":  "Success",
		"messages": message,
		"data": data,
	}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if message == "" {
		message = "Something went wrong"
	}
	response := map[string]interface{}{
		"status": "Error",
		"messages": message,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

