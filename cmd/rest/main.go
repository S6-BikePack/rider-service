package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"rider-service/internal/core/services/rabbitmq_service"
	"rider-service/internal/core/services/rider_service"
	"rider-service/internal/handlers"
	"rider-service/internal/repositories"
	"rider-service/pkg/rabbitmq"

	"github.com/gin-gonic/gin"

	_ "rider-service/docs"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/rider"

func main() {
	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	db, err := gorm.Open(postgres.Open(dbConn))

	if err != nil {
		panic(err)
	}

	riderRepository, err := repositories.NewCockroachDB(db)

	if err != nil {
		panic(err)
	}

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	riderService := rider_service.New(riderRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, riderService)

	router := gin.New()

	riderHandler := handlers.NewHTTPHandler(riderService, router)
	riderHandler.SetupEndpoints()
	riderHandler.SetupSwagger()

	port := GetEnvOrDefault("PORT", defaultPort)

	go rmqSubscriber.Listen("riderQueue")
	log.Fatal(router.Run(port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
