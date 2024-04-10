package json

import (
	"context"
	"encoding/json"
	"net/http"
	"socio/errors"
	"socio/pkg/contextlogger"
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

func ServeJSONBody(ctx context.Context, w http.ResponseWriter, value any) {
	contextlogger.LogInfo(ctx)

	data, err := MarshalResponseBody(value)
	if err != nil {
		ServeJSONError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.Write(data)
}

func ServeJSONError(ctx context.Context, w http.ResponseWriter, err error) {
	contextlogger.LogErr(ctx, err)

	msg, status := errors.ParseHTTPError(err)

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(status)
	w.Write(MarshalResponseError(msg))
}
