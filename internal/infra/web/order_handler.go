package web

import (
	"clean_architecture/internal/infra/database"
	"clean_architecture/internal/usecase"
	"clean_architecture/pkg/dispatcher"
	"encoding/json"
	"net/http"
)

type WebOrderHandler struct {
	EventDispatcher   dispatcher.EventDispatcherInterface
	OrderRepository   database.OrderRepositoryInterface
	OrderCreatedEvent dispatcher.EventInterface
}

func NewWebOrderHandler(
	eventDispatcher dispatcher.EventDispatcherInterface,
	orderRepository database.OrderRepositoryInterface,
	orderCreatedEvent dispatcher.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   eventDispatcher,
		OrderRepository:   orderRepository,
		OrderCreatedEvent: orderCreatedEvent,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Proccess input
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute use case
	createOrder := usecase.NewCreateOrderUseCase(
		h.OrderRepository,
		h.OrderCreatedEvent,
		h.EventDispatcher,
	)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Proccess output
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	// Execute use case
	listOrders := usecase.NewListOrdersUseCase(
		h.OrderRepository,
	)
	output, err := listOrders.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Proccess output
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
