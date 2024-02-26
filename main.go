package main

import (
	"fmt"
	"net/http"
	"socio/routers"
)

const PORT = ":8080"

func main() {
	rootRouter := routers.NewRootRouter()

	fmt.Printf("started on port %s\n", PORT)
	http.ListenAndServe(PORT, rootRouter)
}
