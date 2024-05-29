package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/alipourhabibi/exercises-journal/echo/config"
	v1 "github.com/alipourhabibi/exercises-journal/echo/internal/handlers/v1"
	"github.com/alipourhabibi/exercises-journal/linkedlist"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type server struct {
	e *echo.Echo
}

func New() *server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}

	return &server{
		e: e,
	}
}

type list struct {
	sync.RWMutex
	l *linkedlist.LinkedList
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (s *server) Start(ctx context.Context) error {
	_, err := v1.New(s.e)
	if err != nil {
		return err
	}

	return s.e.Start(fmt.Sprintf(":%d", config.Confs.Server.Port))
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
