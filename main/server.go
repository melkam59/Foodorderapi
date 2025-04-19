package main

import (
	"foodorderapi/internals/config"
	"foodorderapi/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	config.Databaseinit()

	routes.Foodorderroutes(e)

	e.Logger.Fatal(e.Start(":14500"))

}
