package server

import (
	"fmt"
	"net/http"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/labstack/echo"
)

func (s *EchoServer) GetAllVendors(ctx echo.Context) error {
	vendors, err := s.DB.GetAllVendors(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, vendors)
}

func (s *EchoServer) AddVendor(ctx echo.Context) error {
	vendor := new(models.Vendor)
	if err := ctx.Bind(vendor); err != nil {
		ctx.JSON(http.StatusUnsupportedMediaType, err)
	}
	vendor, err := s.DB.AddVendor(ctx.Request().Context(), vendor)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, vendor)
}

func (s *EchoServer) GetVendorById(ctx echo.Context) error {
	vendorId := ctx.Param("id")

	vendor, err := s.DB.GetVendorById(ctx.Request().Context(), vendorId)

	switch err.(type) {
	case *dberrors.NotFoundError:
		return ctx.JSON(http.StatusNotFound, err)
	case nil:
		return ctx.JSON(http.StatusOK, vendor)
	default:
		return ctx.JSON(http.StatusInternalServerError, err)
	}
}

func (s *EchoServer) UpdateVendor(ctx echo.Context) error {
	vendorId := ctx.Param("id")
	vendor := new(models.Vendor)
	if err := ctx.Bind(&vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if vendorId != vendor.VendorID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match id on body")
	}

	fmt.Printf("Updating vendor with id %s\n\nUpdate values %+v\n", vendorId, vendor)
	vendor, err := s.DB.UpdateVendor(ctx.Request().Context(), vendor)

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

	return ctx.JSON(http.StatusOK, vendor)
}

func (s *EchoServer) DeleteVendor(ctx echo.Context) error {
	vendorId := ctx.Param("id")
	err := s.DB.DeleteVendor(ctx.Request().Context(), vendorId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusResetContent, "")
}
