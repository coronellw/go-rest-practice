package server

import (
	"fmt"
	"net/http"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/labstack/echo"
)

func (s *EchoServer) GetAllCustomers(ctx echo.Context) error {
	emailAddress := ctx.QueryParam("emailAddress")

	customers, err := s.DB.GetAllCustomers(ctx.Request().Context(), emailAddress)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, customers)
}

func (s *EchoServer) AddCustomer(ctx echo.Context) error {
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	customer, err := s.DB.AddCustomer(ctx.Request().Context(), customer)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}

	}
	return ctx.JSON(http.StatusCreated, customer)
}

func (s *EchoServer) GetCustomerById(ctx echo.Context) error {
	customerId := ctx.Param("id")
	customer, err := s.DB.GetCustomerById(ctx.Request().Context(), customerId)

	switch err.(type) {
	case *dberrors.NotFoundError:
		return ctx.JSON(http.StatusNotFound, err)
	case nil:
		return ctx.JSON(http.StatusOK, customer)
	default:
		return ctx.JSON(http.StatusInternalServerError, err)
	}
}

func (s *EchoServer) UpdateCustomer(ctx echo.Context) error {
	customerId := ctx.Param("id")
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if customerId != customer.CustomerID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match id on body")
	}

	fmt.Printf("Updating user with id %s\n\nUpdate values %+v\n", customerId, customer)
	customer, err := s.DB.UpdateCustomer(ctx.Request().Context(), customer)

	if err != nil {

		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, customer)
}

func (s *EchoServer) DeleteCustomer(ctx echo.Context) error {
	customerId := ctx.Param("id")
	err := s.DB.DeleteCustomer(ctx.Request().Context(), customerId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusResetContent, "deleted")
}
