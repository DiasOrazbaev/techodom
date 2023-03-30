package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"techodom/internal/entity"
	"techodom/internal/utils"
	"techodom/pkg/mw"
)

type adminUsecase interface {
	FindByID(id string) (*entity.Redirect, error)
	All(page int, perPage int) ([]*entity.Redirect, error)
	Create(old, new string) error
	Update(old, new string, id int) error
	Delete(id int) error
}

type AdminRedirect struct {
	usecase adminUsecase
	log     *zap.Logger
	key     string
}

func NewAdminRedirect(usecase adminUsecase, log *zap.Logger, key string) *AdminRedirect {
	return &AdminRedirect{usecase: usecase, log: log, key: key}
}

func (a *AdminRedirect) Get(ctx echo.Context) error {
	// pagination
	page := ctx.QueryParam("page")
	if page == "" {
		page = "1"
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	limit := ctx.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	// get all redirects
	redirects, err := a.usecase.All(p, l)
	if err != nil {
		a.log.Error("error while getting all redirects", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, redirects)
}

func (a *AdminRedirect) Create(ctx echo.Context) error {
	var r entity.Redirect
	if err := ctx.Bind(&r); err != nil {
		a.log.Error("error while binding redirect", zap.Error(err))
		return ctx.NoContent(http.StatusBadRequest)
	}
	if err := a.usecase.Create(r.HistoryLink, r.ActiveLink); err != nil {
		a.log.Error("error while creating redirect", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusCreated)
}

func (a *AdminRedirect) Update(ctx echo.Context) error {
	var r entity.Redirect
	if err := ctx.Bind(&r); err != nil {
		a.log.Error("error while binding redirect", zap.Error(err))
		return ctx.NoContent(http.StatusBadRequest)
	}
	id := ctx.Param("id")

	ic, err := strconv.Atoi(id)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err = a.usecase.Update(r.HistoryLink, r.ActiveLink, ic); err != nil {
		a.log.Error("error while updating redirect", zap.Error(err))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}

func (a *AdminRedirect) Delete(ctx echo.Context) error {
	id := ctx.Param("id")

	ic, err := strconv.Atoi(id)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err := a.usecase.Delete(ic); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

func (a *AdminRedirect) Register(e *echo.Echo) {
	r := e.Group("/admin")
	r.Use(mw.NewAdminMW(a.key).IsAdmin)
	r.GET("/redirects", a.Get)
	r.POST("/redirects", a.Create)
	r.PUT("/redirects/:id", a.Update)
	r.DELETE("/redirects:id", a.Delete)
	r.GET("/redirects/:id", a.GetByID)
}

func (a *AdminRedirect) GetByID(c echo.Context) error {
	id := c.Param("id")
	redirect, err := a.usecase.FindByID(id)
	if errors.Is(err, utils.ErrNotFound) {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, redirect)
}
