package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"rider-service/internal/core/ports"
	"rider-service/internal/core/services/rabbitmq_service"
	"rider-service/internal/core/services/rider_service"
	"rider-service/internal/graph"
	"rider-service/internal/graph/generated"
	"rider-service/internal/handlers"
	"rider-service/internal/repositories"
	"rider-service/pkg/rabbitmq"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/rider"

func main() {
	port := GetEnvOrDefault("PORT", defaultPort)

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	db, err := gorm.Open(postgres.Open(dbConn))

	if err != nil {
		panic(err)
	}

	riderRepository, err := repositories.NewCockroachDB(db)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	riderService := rider_service.New(riderRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, riderService)

	router := gin.Default()
	router.POST("/query", graphqlHandler(riderService))
	router.GET("/", playgroundHandler())

	go rmqSubscriber.Listen("riderQueue")
	log.Fatal(router.Run(port))
}

func graphqlHandler(riderService ports.RiderService) gin.HandlerFunc {

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{RiderService: riderService}}))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	srv := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
