package server

import (
	"errors"
	"net/http"

	"github.com/deparr/api/pkg/cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ListenAndServe(host, port string) error {
	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/gh/pinned", getPinned)

	addr := host + ":" + port

	cache.InitCache()

	if err := server.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
