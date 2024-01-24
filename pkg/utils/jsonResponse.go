package utils

import (
	"net/http"

	"github.com/goccy/go-json"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseDTO struct {
	StatusCode int
	Success    bool
	Message    string
	Data       interface{}
}

func JsonResponse(w http.ResponseWriter, response ResponseDTO) {
	res := Response{
		Message: response.Message,
		Data:    response.Data,
		Success: response.Success,
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(jsonResponse)
}
