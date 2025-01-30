package database

import (
	"context"
	"errors"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (c Client) GetAllVendors(ctx context.Context) ([]models.Vendor, error) {
	var vendors []models.Vendor

	result := c.DB.WithContext(ctx).
		Find(&vendors)

	return vendors, result.Error
}

func (c Client) AddVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	vendor.VendorID = uuid.NewString()
	err := c.DB.WithContext(ctx).Create(&vendor).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, err
	}

	return vendor, nil
}
