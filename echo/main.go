package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/alipourhabibi/exercises-journal/linkedlist"
	"github.com/go-playground/validator"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type list struct {
	sync.RWMutex
	l *linkedlist.LinkedList
}

type listData struct {
	Index uint `json:"index"`
	Value int  `json:"value" validate:"required"`
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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}

	l := linkedlist.New()
	list := list{
		l: l,
	}

	e.PUT("/list", func(c echo.Context) error {
		list.Lock()
		defer list.Unlock()

		data := listData{}
		if err := c.Bind(&data); err != nil {
			return err
		}
		if err := c.Validate(&data); err != nil {
			return err
		}

		ok := l.Insert(data.Index, data.Value)
		if !ok {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
		}
		return nil
	})

	e.DELETE("/list/:index", func(c echo.Context) error {
		list.Lock()
		defer list.Unlock()

		indexStr := c.Param("index")
		index, err := strconv.ParseUint(indexStr, 10, 32)
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
		}

		ok := list.l.Remove(uint(index))
		if !ok {
			return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
		}

		c.NoContent(http.StatusOK)
		return nil
	})

	e.GET("/list/value/:value", func(c echo.Context) error {
		list.RLock()
		defer list.RUnlock()

		fmt.Println(c.ParamNames())
		valueStr := c.Param("value")
		fmt.Println(valueStr)
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid value")
		}

		index, ok := list.l.Find(value)
		if !ok {
			return echo.NewHTTPError(echo.ErrNotFound.Code, "Value not found")
		}

		data := listData{
			Index: index,
			Value: value,
		}
		c.JSON(http.StatusOK, data)
		return nil
	})

	e.GET("/list/index/:index", func(c echo.Context) error {
		list.RLock()
		defer list.RUnlock()

		indexStr := c.Param("index")
		index, err := strconv.ParseUint(indexStr, 10, 32)
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
		}

		value, ok := list.l.Get(uint(index))
		if !ok {
			return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
		}
		data := listData{
			Index: uint(index),
			Value: value,
		}

		c.JSON(http.StatusOK, data)
		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
