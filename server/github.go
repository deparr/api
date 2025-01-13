package server

import (
	"net/http"

	"github.com/deparr/api/cache"
	"github.com/labstack/echo/v4"
)

func getPinned(ctx echo.Context) error {
	pinned := cache.GetGithubPinned()
	res := map[string]any{
		"data": pinned,
	}

	return ctx.JSON(http.StatusOK, res)
}


func getRecent(ctx echo.Context) error {
	recent := cache.GetGithubRecent()
	res := map[string]any{
		"data": recent,
	}

	return ctx.JSON(http.StatusOK, res)
}
