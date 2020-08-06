package api

import (
	"encoding/json"
	"log"
	"net/http"
)

//Status Constants
const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

type responseBody struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

func newSuccessResponseBody(data interface{}) responseBody {
	return responseBody{
		Status: StatusSuccess,
		Data:   data,
	}
}

func newFailResponseBody(data interface{}) responseBody {
	return responseBody{
		Status: StatusFail,
		Data:   data,
	}
}

func newErrorResponseBody(message string, code int, data interface{}) responseBody {
	return responseBody{
		Status:  StatusError,
		Message: message,
		Code:    code,
		Data:    data,
	}
}

func sendSuccessResponse(w http.ResponseWriter, data interface{}) {
	body := newSuccessResponseBody(data)
	b, marshalError := json.Marshal(body)
	//If there is a marshalling error, it is due to the interface passed in,
	//This should only occur if there is a bug with our struct that is passed in.
	if marshalError != nil {
		log.Println(marshalError)
		errorBytes, _ := json.Marshal(newErrorResponseBody("Error marshaling response", 0, nil))
		writeResponse(w, errorBytes, http.StatusInternalServerError)
	}
	writeResponse(w, b, http.StatusOK)
}

func unauthorizedResponse(w http.ResponseWriter, data interface{}) {
	body := newFailResponseBody(data)
	b, marshalError := json.Marshal(body)
	//If there is a marshalling error, it is due to the interface passed in,
	//This should only occur if there is a bug with our struct that is passed in.
	if marshalError != nil {
		log.Println(marshalError)
		errorBytes, _ := json.Marshal(newErrorResponseBody("Error marshaling response", 0, nil))
		writeResponse(w, errorBytes, http.StatusInternalServerError)
	}
	writeResponse(w, b, http.StatusUnauthorized)
}

func sendFailResponse(w http.ResponseWriter, data interface{}) {
	body := newFailResponseBody(data)
	b, marshalError := json.Marshal(body)
	//If there is a marshalling error, it is due to the interface passed in,
	//This should only occur if there is a bug with our struct that is passed in.
	if marshalError != nil {
		log.Println(marshalError)
		errorBytes, _ := json.Marshal(newErrorResponseBody("Error marshaling response", 0, nil))
		writeResponse(w, errorBytes, http.StatusInternalServerError)
	}
	writeResponse(w, b, http.StatusBadRequest)
}

//Add internal server error response....

func writeResponse(w http.ResponseWriter, payloadData []byte, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(payloadData)
}
