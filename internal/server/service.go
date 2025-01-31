package server

import (
	"fmt"
	"net/http"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/labstack/echo"
)

func (s *EchoServer) GetAllServices(ctx echo.Context) error {
	services, err := s.DB.GetAllServices(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, services)
}

func (s *EchoServer) AddService(ctx echo.Context) error {
	service := new(models.Service)
	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	service, err := s.DB.AddService(ctx.Request().Context(), service)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, service)
}

func (s *EchoServer) GetServiceById(ctx echo.Context) error {
	serviceId := ctx.Param("id")
	service, err := s.DB.GetServiceById(ctx.Request().Context(), serviceId)

	switch err.(type) {
	case *dberrors.NotFoundError:
		return ctx.JSON(http.StatusNotFound, err)
	case nil:
		return ctx.JSON(http.StatusOK, service)
	default:
		return ctx.JSON(http.StatusInternalServerError, err)
	}
}

func (s *EchoServer) UpdateService(ctx echo.Context) error {
	serviceId := ctx.Param("id")
	service := new(models.Service)

	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if serviceId != service.ServiceID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match id on body")
	}

	fmt.Printf("Updating service with id %s\n\nUpdate values %+v\n", serviceId, service)
	service, err := s.DB.UpdateService(ctx.Request().Context(), service)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, service)
}
