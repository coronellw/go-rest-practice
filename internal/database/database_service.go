package database

import (
	"context"

	"github.com/coronellw/go-microservices/internal/models"
)

func (c Client) GetAllServices(ctx context.Context) ([]models.Service, error) {
	var services []models.Service
	err := c.DB.WithContext(ctx).Find(&services).Error

	return services, err
}
