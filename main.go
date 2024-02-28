package main

import (
	"fmt"
	"net/http"
	"socio/routers"
	"socio/utils"
)

func main() {
	rootRouter := routers.NewRootRouter()

	fmt.Printf("started on port %s\n", utils.PORT)
	http.ListenAndServe(utils.PORT, rootRouter)
}
