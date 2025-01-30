package database

import (
	"context"
	"errors"

	"github.com/coronellw/go-microservices/internal/dberrors"
	"github.com/coronellw/go-microservices/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (c Client) GetAllServices(ctx context.Context) ([]models.Service, error) {
	var services []models.Service
	err := c.DB.WithContext(ctx).Find(&services).Error

	return services, err
}

func (c Client) AddService(ctx context.Context, service *models.Service) (*models.Service, error) {
	service.ServiceID = uuid.NewString()
	err := c.DB.WithContext(ctx).Create(&service).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, err
	}

	return service, nil
}

func (c Client) GetServiceById(ctx context.Context, serviceId string) (*models.Service, error) {
	var service models.Service
	err := c.DB.WithContext(ctx).
		Where(models.Service{ServiceID: serviceId}).
		Find(&service).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "service", ID: serviceId}
		}

		return nil, err
	}

	if (models.Service{} == service) {
		return nil, &dberrors.NotFoundError{Entity: "service", ID: serviceId}
	}

	return &service, nil
}
