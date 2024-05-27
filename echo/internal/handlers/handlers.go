package handlers

import (
	"net/http"
	"sync"

	v1 "github.com/alipourhabibi/exercises-journal/echo/internal/handlers/v1"
	"github.com/alipourhabibi/exercises-journal/linkedlist"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

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

func Launch() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}

	_, err := v1.New(e)
	if err != nil {
		return err
	}

	return e.Start(":8080")
}
