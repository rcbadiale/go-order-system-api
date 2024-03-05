package main

import (
	"clean_architecture/configs"
	"clean_architecture/internal/events/handler"
	"clean_architecture/internal/infra/graph"
	"clean_architecture/internal/infra/grpc/pb"
	"clean_architecture/internal/infra/grpc/service"
	"clean_architecture/internal/infra/web/webserver"
	"clean_architecture/pkg/dispatcher"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Read configs from .env with priority for environment variables
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// Open DB connection
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?multiStatements=true",
			configs.DBUser,
			configs.DBPassword,
			configs.DBHost,
			configs.DBPort,
			configs.DBName,
		),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Migrate DB (if needed)
	log.Println("migrating database")
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+configs.DBMigrationFolder,
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			panic(err)
		}
		log.Println("no changes made on DB")
	}

	// Setup RabbitMQ connection
	rabbitMQChannel := getRabbitMQChannel(
		configs.RabbitMQUser,
		configs.RabbitMQPassword,
		configs.RabbitMQHost,
		configs.RabbitMQPort,
	)

	// Create event dispatcher
	eventDispatcher := dispatcher.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	// Creating usecases
	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	// Start webserver
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("POST /order", webOrderHandler.Create)
	webserver.AddHandler("GET /order", webOrderHandler.ReadAll)
	log.Println("starting webserver on port", configs.WebServerPort)
	go webserver.Start()

	// Start gRPC
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(
		*createOrderUseCase,
		*listOrdersUseCase,
	)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)
	log.Println("starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", ":"+configs.GRPCServerPort)
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// Start GraphQL server
	srv := graphql_handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{
				CreateOrderUseCase: *createOrderUseCase,
				ListOrdersUseCase:  *listOrdersUseCase,
			}},
		))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Println("starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(user, passwd, host, port string) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		user,
		passwd,
		host,
		port,
	))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
