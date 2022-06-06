package main

import (
	"github.com/labstack/echo"
)

func main() {

	router := echo.New()

	router.GET("/", func(context echo.Context) error {
		return context.JSON(200, "This endpoint Is allowed by roles: ReadOnly and ReadWrite")
	})

	router.GET("/readwrite", func(context echo.Context) error {
		return context.JSON(200, "This endpoint Is allowed by roles: ReadWrite")
	})

	router.Start("0.0.0.0:4000")
}
