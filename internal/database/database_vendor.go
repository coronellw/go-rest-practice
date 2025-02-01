package database

import (
	"context"
	"errors"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (c Client) GetVendorById(ctx context.Context, vendorId string) (*models.Vendor, error) {
	var vendor models.Vendor

	err := c.DB.WithContext(ctx).
		Where(models.Vendor{VendorID: vendorId}).
		Find(&vendor).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "vendor", ID: vendorId}
		}

		return nil, err
	}

	if (vendor == models.Vendor{}) {
		return nil, &dberrors.NotFoundError{Entity: "vendor", ID: vendorId}
	}

	return &vendor, nil
}

func (c Client) UpdateVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	var vendors []models.Vendor
	result := c.DB.WithContext(ctx).
		Model(&vendors).
		Clauses(clause.Returning{}).
		Where(models.Vendor{VendorID: vendor.VendorID}).
		Updates(models.Vendor{
			Name:    vendor.Name,
			Contact: vendor.Contact,
			Phone:   vendor.Phone,
			Email:   vendor.Email,
			Address: vendor.Address,
		})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}

		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "vendor", ID: vendor.VendorID}
	}

	return &vendors[0], nil
}
func (c Client) DeleteVendor(ctx context.Context, vendorId string) error {
	return c.DB.WithContext(ctx).Delete(models.Vendor{VendorID: vendorId}).Error
}
