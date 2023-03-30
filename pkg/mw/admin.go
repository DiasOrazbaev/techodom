package mw

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type AdminMW struct {
	Key string
}

func NewAdminMW(key string) *AdminMW {
	return &AdminMW{Key: key}
}

func (a *AdminMW) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("X-Admin") != a.Key {
			return c.NoContent(http.StatusForbidden)
		}
		return next(c)
	}
}
