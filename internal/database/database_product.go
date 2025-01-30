package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
