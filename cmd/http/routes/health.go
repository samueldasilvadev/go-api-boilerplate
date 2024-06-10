package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Health struct {
}

func NewHealthRoute() *Health {
	return &Health{}
}

func (hs *Health) DeclarePrivateRoutes(_ *echo.Group, _ string) {}

func (hs *Health) DeclarePublicRoutes(server *echo.Group, _ string) {
	server.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	server.POST("/auth/generate", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
