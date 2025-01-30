package database

import (
	"context"
	"errors"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (c Client) GetAllCustomers(ctx context.Context, emailAddress string) ([]models.Customer, error) {
	var customers []models.Customer
	result := c.DB.WithContext(ctx).
		Where(models.Customer{Email: emailAddress}).
		Find(&customers)

	return customers, result.Error
}

func (c Client) AddCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer.CustomerID = uuid.NewString()
	result := c.DB.WithContext(ctx).
		Create(&customer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}

		return nil, result.Error
	}

	return customer, nil
}

func (c Client) GetCustomerById(ctx context.Context, customerId string) (*models.Customer, error) {
	var customer models.Customer
	err := c.DB.WithContext(ctx).
		Where(models.Customer{CustomerID: customerId}).
		Find(&customer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "customer", ID: customerId}
		}

		return nil, err
	}

	// No error but query execution may return empty which means not found.
	if (customer == models.Customer{}) {
		return nil, &dberrors.NotFoundError{Entity: "customer", ID: customerId}
	}

	return &customer, nil
}
