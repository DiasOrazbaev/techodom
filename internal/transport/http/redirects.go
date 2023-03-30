package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"techodom/internal/utils"
)

type redirectUsecase interface {
	Find(code string) (string, error)
}

type UserRedirect struct {
	usecase redirectUsecase
	log     *zap.Logger
}

func NewUserRedirect(usecase redirectUsecase, log *zap.Logger) *UserRedirect {
	return &UserRedirect{usecase: usecase, log: log}
}

func (u *UserRedirect) Get(ctx echo.Context) error {
	link := ctx.QueryParam("link")
	newLink, err := u.usecase.Find(link)
	if errors.Is(err, utils.ErrNotFound) {
		return ctx.NoContent(http.StatusNotFound)
	} else if err != nil {
		u.log.Error("error while getting link", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if newLink == link {
		return ctx.JSON(http.StatusOK, "link is live")
		//return ctx.NoContent(http.StatusOK)
	} else {
		return ctx.JSON(http.StatusOK, "new link is "+newLink)
		//return ctx.Redirect(http.StatusMovedPermanently, newLink)
	}
}

func (u *UserRedirect) Register(e *echo.Echo) {
	e.GET("/redirects", u.Get)
}
