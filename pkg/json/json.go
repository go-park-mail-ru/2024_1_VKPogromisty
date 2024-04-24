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
		return
	}

	return
}

func MarshalResponseError(errMsg string) (data []byte, err error) {
	data, err = json.Marshal(map[string]string{"error": errMsg})
	if err != nil {
		err = errors.ErrJSONMarshalling
		return
	}

	return
}

func ServeJSONBody(ctx context.Context, w http.ResponseWriter, value any, statusCode int) {
	contextlogger.LogInfo(ctx)

	data, err := MarshalResponseBody(value)
	if err != nil {
		ServeJSONError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json;")

	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		ServeJSONError(ctx, w, err)
		return
	}
}

func ServeGRPCStatus(ctx context.Context, w http.ResponseWriter, err errors.CustomError) {
	contextlogger.LogErr(ctx, err)

	msg, statusCode := errors.ParseGRPCError(err)

	w.Header().Set("Content-Type", "application/json;")

	data, marshallErr := MarshalResponseError(msg)
	if marshallErr != nil {
		contextlogger.LogErr(ctx, marshallErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	_, marshallErr = w.Write(data)
	if marshallErr != nil {
		contextlogger.LogErr(ctx, marshallErr)
		return
	}
}

func ServeJSONError(ctx context.Context, w http.ResponseWriter, err error) {
	contextlogger.LogErr(ctx, err)

	msg, status := errors.ParseHTTPError(err)

	w.Header().Set("Content-Type", "application/json;")

	data, err := MarshalResponseError(msg)
	if err != nil {
		contextlogger.LogErr(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		contextlogger.LogErr(ctx, err)
		return
	}
}
