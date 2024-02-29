package main

import (
	"fmt"
	"net/http"
	"os"
	"socio/routers"
	"socio/utils"
)

func main() {
	rootRouter := routers.NewRootRouter()
	os.Setenv("PROTOCOL", utils.PROTOCOL)
	os.Setenv("HOST", utils.HOST)
	os.Setenv("PORT", utils.PORT)

	fmt.Printf("started on port %s\n", utils.PORT)
	http.ListenAndServe(utils.PORT, rootRouter)
}
