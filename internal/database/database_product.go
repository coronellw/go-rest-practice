package database

import (
	"context"

	"github.com/coronellw/go-microservices/internal/models"
)

func (c Client) GetAllProducts(ctx context.Context, vendorID string) ([]models.Product, error) {
	var products []models.Product
	err := c.DB.WithContext(ctx).
		Where(models.Product{VendorID: vendorID}).
		Find(&products).Error

	return products, err
}
