package usecase

import (
	"clean_architecture/internal/entity"
	"clean_architecture/internal/infra/database"
	"clean_architecture/pkg/dispatcher"
)

type CreateOrderUseCase struct {
	OrderRepository database.OrderRepositoryInterface
	OrderCreated    dispatcher.EventInterface
	EventDispatcher dispatcher.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	orderRepository database.OrderRepositoryInterface,
	orderCreated dispatcher.EventInterface,
	eventDispatcher dispatcher.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		OrderCreated:    orderCreated,
		EventDispatcher: eventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (*OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	err := order.CalculateFinalPrice()
	if err != nil {
		return nil, err
	}

	err = c.OrderRepository.Save(&order)
	if err != nil {
		return nil, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)
	return &dto, nil
}
