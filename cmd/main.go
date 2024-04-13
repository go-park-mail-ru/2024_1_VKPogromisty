package main

import (
	"fmt"
	"net/http"
	"socio/app/routers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "socio/docs"
)

// swag init

// @title			Socio API
// @version		1.0
// @description	First version of Socio API
// @contact.name	Petr Mitin
// @contact.url	https://github.com/Petr09Mitin
// @contact.email	petr09mitin@mail.ru
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	http.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8001/swagger/doc.json"),
	))
	go http.ListenAndServe(":8001", nil)

	router := mux.NewRouter()
	err := routers.MountRootRouter(router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
