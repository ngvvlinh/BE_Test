package api

import (
	"encoding/json"
	"net/http"
)

func responseJsonData(jsonData string, w http.ResponseWriter) {
	data := map[string]interface{}{
		"success": true,
		"status":  http.StatusOK,
		"data":    jsonData,
	}
	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func responseJsonSuccessEx(message string, w http.ResponseWriter) {

	if len(message) < 1 {
		message = "Request succeeded"
	}

	data := map[string]interface{}{
		"success": true,
		"status":  http.StatusOK,
		"info":    message,
	}
	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func responseJsonSuccess(w http.ResponseWriter) {
	responseJsonSuccessEx("Request succeeded", w)
}

func responseJsonErrorEx(message string, status int, w http.ResponseWriter) {

	if status == 0 {
		status = http.StatusNotAcceptable
	}

	if len(message) < 1 {
		message = "Request not acceptable"
	}

	data := map[string]interface{}{
		"success": false,
		"status":  status,
		"error":   message,
	}
	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func responseJsonError(w http.ResponseWriter) {
	responseJsonErrorEx("Request not acceptable", http.StatusNotAcceptable, w)
}

func apiAuthenticate(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if len(auth) < 1 {
		return false
	}

	return false
}

//var jsonSuccess = map[string]string{"test": "test", "data": "test1"}

func JsonSuccess(message string, data ...interface{}) interface{} {
	return jsonSuccess(message, data)
}

func jsonSuccess(message string, data []interface{}) interface{} {

	obj := map[string]interface{}{
		"success": true,
		"code":    200,
		"message": message,
	}

	if len(data) > 0 && data[0] != nil {
		obj["data"] = data[0]
	} else {
		obj["data"] = []string{}
	}

	return obj
}

func JsonError(code int, message string) interface{} {
	return jsonError(code, []string{message})
}

func JsonErrors(code int, messages []string) interface{} {
	return jsonError(code, messages)
}

func jsonError(code int, messages []string) interface{} {
	obj := map[string]interface{}{
		"success": false,
		"code":    code,
		"errors":  messages,
	}

	return obj
}
