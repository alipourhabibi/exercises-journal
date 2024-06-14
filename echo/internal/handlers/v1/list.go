package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alipourhabibi/exercises-journal/echo/internal/core/list"
	"github.com/labstack/echo"
)

type server struct {
	list *list.ListService
}

func New(e *echo.Echo) (*server, error) {
	s := &server{}
	l, err := list.New(
		list.BootList(),
	)
	if err != nil {
		return nil, err
	}
	s.list = l
	v1 := e.Group("/api/v1")

	v1.PUT("/numbers", s.Insert)
	v1.DELETE("/numbers/:index", s.Remove)
	v1.GET("/numbers/value/:value", s.Find)
	v1.GET("/numbers/index/:index", s.Get)

	return s, nil
}

func (s *server) Insert(c echo.Context) error {
	data := list.ListEntity{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ok := s.list.Insert(data.Index, data.Value)
	if !ok {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}
	c.JSON(http.StatusCreated, data)
	return nil
}

func (s *server) Remove(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}

	ok := s.list.Remove(uint(index))
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
	}

	c.NoContent(http.StatusOK)
	return nil

}

func (s *server) Find(c echo.Context) error {
	valueStr := c.Param("value")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid value")
	}

	index, ok := s.list.Find(value)
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Value not found")
	}

	data := list.ListEntity{
		Index: index,
		Value: value,
	}
	c.JSON(http.StatusOK, data)
	return nil
}

func (s *server) Get(c echo.Context) error {

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}

	value, ok := s.list.Get(uint(index))
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
	}
	data := list.ListEntity{
		Index: uint(index),
		Value: value,
	}

	c.JSON(http.StatusOK, data)
	return nil
}
