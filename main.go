package main

import (
	"fmt"
	"net/http"
	"os"
	"socio/routers"
)

const (
	PROTOCOL = "http://"
	HOST     = "localhost"
	PORT     = ":8080"
)

func main() {
	rootRouter := routers.NewRootRouter()
	os.Setenv("PROTOCOL", PROTOCOL)
	os.Setenv("HOST", HOST)
	os.Setenv("PORT", PORT)

	fmt.Printf("started on port %s\n", PORT)
	http.ListenAndServe(PORT, rootRouter)
}
