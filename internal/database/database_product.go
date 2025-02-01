package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) GetAllProducts(ctx context.Context, vendorID string) ([]models.Product, error) {
	var products []models.Product
	err := c.DB.WithContext(ctx).
		Where(models.Product{VendorID: vendorID}).
		Find(&products).Error

	return products, err
}

func (c Client) AddProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.ProductID = uuid.NewString()
	fmt.Printf("%+v\n", product)
	err := c.DB.WithContext(ctx).Create(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}

		return nil, err
	}

	return product, nil
}

func (c Client) GetProductById(ctx context.Context, productId string) (*models.Product, error) {
	var product models.Product
	err := c.DB.WithContext(ctx).
		Where(models.Product{ProductID: productId}).
		Find(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "product", ID: productId}
		}

		return nil, err
	}

	if (product == models.Product{}) {
		return nil, &dberrors.NotFoundError{Entity: "product", ID: productId}
	}

	return &product, nil
}

func (c Client) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	var products []models.Product
	result := c.DB.WithContext(ctx).
		Model(&products).
		Clauses(clause.Returning{}).
		Where(models.Product{ProductID: product.ProductID}).
		Updates(models.Product{
			Name:     product.Name,
			Price:    product.Price,
			VendorID: product.VendorID,
		})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}

		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "product", ID: product.ProductID}
	}

	return &products[0], nil
}

func (c Client) DeleteProduct(ctx context.Context, productId string) error {
	return c.DB.WithContext(ctx).Delete(models.Product{ProductID: productId}).Error
}
