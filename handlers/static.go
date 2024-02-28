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
		errors.ServeHttpError(&w, errors.ErrInvalidFileName)
		return
	}

	filePath := path.Join("./static", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		errors.ServeHttpError(&w, errors.ErrInvalidFilePathGen)
		return
	}
	file.Close()

	http.ServeFile(w, r, filePath)
}
