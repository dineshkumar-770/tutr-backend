package utils

import (
	"encoding/json" 
	"net/http"
)

func SendResponseWithOK(w http.ResponseWriter, arg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(arg)

}

func SendResponseWithServerError(w http.ResponseWriter, arg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(arg)
}

func SendResponseWithStatusBadRequest(w http.ResponseWriter, arg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(arg)
}
func SendResponseWithStatusNotFound(w http.ResponseWriter, arg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(arg)
}

func SendResponseWithUnauthorizedError(w http.ResponseWriter) {
	resp := ResponseStr{
		Status:     "failed",
		Message:    "unauthorized request",
		MyResponse: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(resp)
}

func SendResponseWithMissingValues(w http.ResponseWriter){
	resp := ResponseStr{
		Status:     "failed",
		Message:    "missing values not allowed",
		MyResponse: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)
	json.NewEncoder(w).Encode(resp)
}