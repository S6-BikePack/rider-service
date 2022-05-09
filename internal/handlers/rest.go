package handlers

import (
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/interfaces"
	"rider-service/pkg/authorization"
	"rider-service/pkg/dto"
	"rider-service/pkg/logging"

	ginSwagger "github.com/swaggo/gin-swagger"
	"rider-service/docs"
	_ "rider-service/docs"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	riderService interfaces.RiderService
	router       *gin.Engine
	logger       logging.Logger
	config       *config.Config
}

func NewHTTPHandler(riderService interfaces.RiderService, router *gin.Engine, logger logging.Logger, config *config.Config) *HTTPHandler {
	return &HTTPHandler{
		riderService: riderService,
		router:       router,
		logger:       logger,
		config:       config,
	}
}

func (handler *HTTPHandler) SetupEndpoints() {
	api := handler.router.Group("/api")
	api.GET("/riders", handler.GetAll)
	api.GET("/riders/:id", handler.Get)
	api.POST("/riders", handler.Create)
	api.PUT("/riders/:id", handler.UpdateRider)
	api.PUT("/riders/:id/location", handler.UpdateLocation)
}

func (handler *HTTPHandler) SetupSwagger() {
	docs.SwaggerInfo.Title = handler.config.Server.Service + " API"
	docs.SwaggerInfo.Description = handler.config.Server.Description

	handler.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetAll godoc
// @Summary  get all riders
// @Schemes
// @Description  gets all riders in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.RiderListResponse
// @Router       /api/riders [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if authorization.NewRest(c).AuthorizeAdmin() {
		riders, err := handler.riderService.GetAll(ctx)

		if err != nil {
			handler.logger.Error(ctx, err.Error(), "error", err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, dto.CreateServiceAreaListResponse(riders))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)

}

// Get godoc
// @Summary  get rider
// @Schemes
// @Param        id     path  string           true  "Rider id"
// @Description  gets a rider from the system by its ID
// @Produce      json
// @Success      200  {object}  dto.RiderResponse
// @Router       /api/riders/{id} [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		rider, err := handler.riderService.Get(ctx, c.Param("id"))

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, dto.CreateRiderResponse(rider))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// Create godoc
// @Summary  create rider
// @Schemes
// @Description  creates a new rider
// @Accept       json
// @Param        rider  body  dto.BodyCreateRider  true  "Add rider"
// @Produce      json
// @Success      200  {object}  dto.RiderResponse
// @Router       /api/riders [post]
func (handler *HTTPHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	body := dto.BodyCreateRider{}
	err := c.BindJSON(&body)

	if err != nil || body.ID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(body.ID) {

		rider, err := handler.riderService.Create(ctx, body.ID, body.ServiceArea, domain.Dimensions(body.Capacity))

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			handler.logger.Error(ctx, err.Error())
			return
		}

		c.JSON(http.StatusOK, dto.CreateRiderResponse(rider))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// UpdateRider godoc
// @Summary  update rider
// @Schemes
// @Description  updates a rider's information
// @Accept       json
// @Param        rider  body  dto.BodyCreateRider  true  "Update rider"
// @Param        id     path  string      true  "Rider id"
// @Produce      json
// @Success      200  {object}  dto.RiderResponse
// @Router       /api/riders/{id} [put]
func (handler *HTTPHandler) UpdateRider(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	body := dto.BodyCreateRider{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := authorization.NewRest(c)
	riderId := c.Param("id")

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(riderId) {

		handler.logger.Info(ctx, "Updating rider position", "rider", riderId, "body", body)

		rider, err := handler.riderService.Update(ctx, riderId, body.Status, body.ServiceArea, domain.Dimensions(body.Capacity))

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			handler.logger.Error(ctx, err.Error())
			return
		}

		c.JSON(http.StatusOK, dto.CreateRiderResponse(rider))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// UpdateLocation godoc
// @Summary  update rider location
// @Schemes
// @Description  updates a rider's location
// @Accept       json
// @Param        rider  body  dto.BodyLocation  true  "Update rider"
// @Param        id  path  string  true  "Rider id"
// @Produce      json
// @Success      200  {object}  dto.RiderResponse
// @Router       /api/riders/{id}/location [put]
func (handler *HTTPHandler) UpdateLocation(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	body := dto.BodyLocation{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		id := c.Param("id")

		rider, err := handler.riderService.UpdateLocation(ctx, id, domain.Location(body))

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			handler.logger.Error(ctx, err.Error())
			return
		}

		c.JSON(http.StatusOK, dto.CreateRiderResponse(rider))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
