package database

import "clean_architecture/internal/entity"

type OrderRepositoryInterface interface {
	Save(order *entity.Order) error
	// GetTotal() (int, error)
	ReadAll() ([]entity.Order, error)
}
