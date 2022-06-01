package main

import (
	"context"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"rider-service/config"
	"rider-service/internal/core/services"
	"rider-service/internal/handlers"
	"rider-service/internal/repositories"
	"rider-service/pkg/logging"
	"rider-service/pkg/rabbitmq"
	"rider-service/pkg/tracing"

	"github.com/gin-gonic/gin"

	_ "rider-service/docs"
)

const defaultConfig = "./config/local.config"

func main() {
	cfgPath := GetEnvOrDefault("config", defaultConfig)
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Logging and Tracing
	//--------------------------------------------------------------------------------------

	logger, err := logging.NewSimpleLogger(cfg)

	if err != nil {
		panic(err)
	}

	tracer, err := tracing.NewOpenTracing(cfg.Server.Service, cfg.Tracing.Host, cfg.Tracing.Port)

	if err != nil {
		logger.Panic(context.Background(), err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Database
	//--------------------------------------------------------------------------------------

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database, cfg.Database.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		logger.Panic(context.Background(), err)
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	if err = db.Use(otelgorm.NewPlugin(otelgorm.WithTracerProvider(tracer))); err != nil {
		panic(err)
	}

	if err != nil {
		logger.Panic(context.Background(), err)
	}

	serviceAreaRepository, err := repositories.NewServiceAreaRepository(db)

	if err != nil {
		logger.Panic(context.Background(), err)
	}

	riderRepository, err := repositories.NewRiderRepository(db)

	if err != nil {
		logger.Panic(context.Background(), err)
	}

	//--------------------------------------------------------------------------------------
	// Setup RabbitMQ
	//--------------------------------------------------------------------------------------

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		logger.Panic(context.Background(), err)
	}

	rmqPublisher := services.NewRabbitMQPublisher(rmqServer, tracer, cfg)

	//--------------------------------------------------------------------------------------
	// Setup Services
	//--------------------------------------------------------------------------------------

	serviceAreaService := services.NewServiceAreaService(serviceAreaRepository)
	riderService := services.NewRiderService(riderRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, riderService, serviceAreaService, cfg)

	//--------------------------------------------------------------------------------------
	// Setup HTTP server
	//--------------------------------------------------------------------------------------

	router := gin.New()
	router.Use(otelgin.Middleware(cfg.Server.Service, otelgin.WithTracerProvider(tracer)))

	riderHandler := handlers.NewHTTPHandler(riderService, router, logger, cfg)
	riderHandler.SetupEndpoints()
	riderHandler.SetupSwagger()

	go rmqSubscriber.Listen()
	logger.Fatal(context.Background(), router.Run(cfg.Server.Port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
