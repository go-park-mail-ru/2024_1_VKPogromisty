package handlers

import (
	"net/http"
	"os"
	"path"
	"socio/utils"

	"github.com/gorilla/mux"
)

type StaticHandler struct {
}

func (s *StaticHandler) HandleServeStatic(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["fileName"]
	if len(fileName) == 0 {
		http.Error(w, utils.ErrBadRequest.Error(), 400)
		return
	}

	filePath := path.Join("./static", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, utils.ErrNotFound.Error(), 404)
		return
	}
	file.Close()

	http.ServeFile(w, r, filePath)
}
