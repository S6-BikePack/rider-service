package main

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"rider-service/internal/core/services"
	"rider-service/internal/core/services/rabbitmq_service"
	"rider-service/internal/core/services/rider_service"
	"rider-service/internal/handlers"
	"rider-service/internal/repositories"
	"rider-service/pkg/logging"
	"rider-service/pkg/rabbitmq"
	"rider-service/pkg/tracing"

	"github.com/gin-gonic/gin"

	_ "rider-service/docs"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/rider"
const defaultTracingUrl = "localhost"
const defaultTracingPort = "6831"

func main() {
	logger, err := logging.NewSugaredOtelZap()
	defer logger.Close()

	if err != nil {
		panic(err)
	}

	logger.Logger.Info(context.Background(), "Starting rider-service")

	tracingUrl := GetEnvOrDefault("TRACING_URL", defaultTracingUrl)
	tracingPort := GetEnvOrDefault("TRACING_PORT", defaultTracingPort)

	tracer, err := tracing.NewOpenTracing("rider-service", tracingUrl, tracingPort)

	if err != nil {
		logger.Logger.Panic(err)
	}

	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	db, err := gorm.Open(postgres.Open(dbConn))

	if err != nil {
		logger.Logger.Panic(err)
	}

	if err = db.Use(otelgorm.NewPlugin(otelgorm.WithTracerProvider(tracer))); err != nil {
		panic(err)
	}

	if err != nil {
		logger.Logger.Panic(err)
	}

	serviceAreaRepository, err := repositories.NewServiceAreaRepository(db)

	if err != nil {
		logger.Logger.Panic(err)
	}

	riderRepository, err := repositories.NewCockroachDB(db)

	if err != nil {
		logger.Logger.Panic(err)
	}

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		logger.Logger.Panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer, tracer)

	serviceAreaService := services.NewServiceAreaService(serviceAreaRepository)
	riderService := rider_service.New(riderRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, riderService, serviceAreaService)

	router := gin.New()
	router.Use(otelgin.Middleware("rider-service", otelgin.WithTracerProvider(tracer)))

	riderHandler := handlers.NewHTTPHandler(riderService, router, logger)
	riderHandler.SetupEndpoints()
	riderHandler.SetupSwagger()

	port := GetEnvOrDefault("PORT", defaultPort)

	go rmqSubscriber.Listen("riderQueue")
	logger.Logger.Fatal(router.Run(port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
