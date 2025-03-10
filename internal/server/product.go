package server

import (
	"fmt"
	"net/http"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/labstack/echo"
)

func (s *EchoServer) GetAllProducts(ctx echo.Context) error {
	vendorId := ctx.QueryParam("vendorId")
	products, err := s.DB.GetAllProducts(ctx.Request().Context(), vendorId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, products)
}

func (s *EchoServer) AddProduct(ctx echo.Context) error {
	product := new(models.Product)
	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	product, err := s.DB.AddProduct(ctx.Request().Context(), product)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, product)
}

func (s *EchoServer) GetProductById(ctx echo.Context) error {
	productId := ctx.Param("id")
	product, err := s.DB.GetProductById(ctx.Request().Context(), productId)

	switch err.(type) {
	case *dberrors.NotFoundError:
		return ctx.JSON(http.StatusNotFound, err)
	case nil:
		return ctx.JSON(http.StatusOK, product)
	default:
		return ctx.JSON(http.StatusInternalServerError, err)
	}
}

func (s *EchoServer) UpdateProduct(ctx echo.Context) error {
	productId := ctx.Param("id")
	product := new(models.Product)
	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	if productId != product.ProductID {
		return ctx.JSON(http.StatusBadRequest, "id on path does not match id on body")
	}

	fmt.Printf("Updating product with id %s\n\nUpdate values %+v\n", productId, product)
	product, err := s.DB.UpdateProduct(ctx.Request().Context(), product)

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

	return ctx.JSON(http.StatusOK, product)
}

func (s *EchoServer) DeleteProduct(ctx echo.Context) error {
	productId := ctx.Param("id")
	err := s.DB.DeleteProduct(ctx.Request().Context(), productId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusResetContent, "")
}
