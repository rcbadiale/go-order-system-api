package database

import (
	"database/sql"
	"fmt"
	"testing"

	"clean_architecture/internal/entity"

	"github.com/stretchr/testify/suite"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", "file::memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestGivenMultipleOrders_WhenReadAll_ThenShouldReturnAllOrders() {
	// Create Repository
	repo := NewOrderRepository(suite.Db)

	var ordersCreated []*entity.Order
	for i := range 10 {
		// Create order
		order, err := entity.NewOrder(
			fmt.Sprintf("%d", i),
			float64(i+1),
			float64(i+1),
		)
		suite.NoError(err)
		suite.NoError(order.CalculateFinalPrice())
		// Save order
		err = repo.Save(order)
		suite.NoError(err)
		// Register data for assertion later
		ordersCreated = append(ordersCreated, order)
	}

	// Read all orders
	var ordersResult []entity.Order
	ordersResult, err := repo.ReadAll()
	suite.NoError(err)
	suite.Equal(10, len(ordersResult))

	// Assert order response with registered data
	for i := range 10 {
		suite.Equal(ordersCreated[i].ID, ordersResult[i].ID)
		suite.Equal(ordersCreated[i].Price, ordersResult[i].Price)
		suite.Equal(ordersCreated[i].Tax, ordersResult[i].Tax)
		suite.Equal(ordersCreated[i].FinalPrice, ordersResult[i].FinalPrice)
	}
}
