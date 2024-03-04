package handlers

import (
	"net/http"
	"os"
	"path"
	"socio/errors"
	"socio/utils"

	"github.com/gorilla/mux"
)

type StaticHandler struct {
}

func (s *StaticHandler) HandleServeStatic(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["fileName"]
	if len(fileName) == 0 {
		utils.ServeJSONError(w, errors.ErrInvalidFileName)
		return
	}

	filePath := path.Join(utils.StaticFilePath, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		utils.ServeJSONError(w, errors.ErrNotFound)
		return
	}
	defer file.Close()

	http.ServeFile(w, r, filePath)
}
