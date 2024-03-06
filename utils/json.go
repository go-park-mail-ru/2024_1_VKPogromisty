package utils

import (
	"encoding/json"
	"net/http"
	"socio/errors"
)

type JSONResponse struct {
	Body any `json:"body"`
}

func MarshalResponseBody(value any) (data []byte, err error) {
	data, err = json.Marshal(map[string]any{"body": value})
	if err != nil {
		err = errors.ErrJSONMarshalling
		data = MarshalResponseError(err.Error())
	}
	return
}

func MarshalResponseError(errMsg string) (data []byte) {
	data, _ = json.Marshal(map[string]string{"error": errMsg})
	return
}

func ServeJSONError(w http.ResponseWriter, err error) {
	msg, status := errors.ParseHTTPError(err)

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(status)
	w.Write(MarshalResponseError(msg))
}

func ServeJSONBody(w http.ResponseWriter, value any) {
	data, err := MarshalResponseBody(value)
	if err != nil {
		ServeJSONError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.Write(data)
}
