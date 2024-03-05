package usecase

import (
	"clean_architecture/internal/infra/database"
)

type ListOrdersUseCase struct {
	OrderRepository database.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	orderRepository database.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := c.OrderRepository.ReadAll()
	if err != nil {
		return nil, err
	}

	var ordersDto []OrderOutputDTO
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		ordersDto = append(ordersDto, dto)
	}
	return ordersDto, nil
}
