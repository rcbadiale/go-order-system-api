//go:build wireinject
// +build wireinject

package main

import (
	"clean_architecture/internal/events"
	"clean_architecture/internal/infra/database"
	"clean_architecture/internal/infra/web"
	"clean_architecture/internal/usecase"
	"clean_architecture/pkg/dispatcher"
	"database/sql"

	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(database.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	dispatcher.NewEventDispatcher,
	events.NewOrderCreated,
	wire.Bind(new(dispatcher.EventInterface), new(*events.OrderCreated)),
	wire.Bind(new(dispatcher.EventDispatcherInterface), new(*dispatcher.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	events.NewOrderCreated,
	wire.Bind(new(dispatcher.EventInterface), new(*events.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher dispatcher.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher dispatcher.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
