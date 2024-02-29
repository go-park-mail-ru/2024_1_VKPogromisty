package handlers

import (
	"net/http"
	"os"
	"path"
	"socio/errors"

	"github.com/gorilla/mux"
)

type StaticHandler struct {
}

func (s *StaticHandler) HandleServeStatic(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["fileName"]
	if len(fileName) == 0 {
		msg, status := errors.ServeHttpError(errors.ErrInvalidFileName)
		http.Error(w, msg, status)
		return
	}

	filePath := path.Join("./static", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		msg, status := errors.ServeHttpError(errors.ErrInvalidFilePathGen)
		http.Error(w, msg, status)
		return
	}
	defer file.Close()

	http.ServeFile(w, r, filePath)
}
