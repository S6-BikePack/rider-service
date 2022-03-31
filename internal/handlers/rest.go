package handlers

import (
	"github.com/google/uuid"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/ports"
	"rider-service/pkg/dto"

	ginSwagger "github.com/swaggo/gin-swagger"
	"rider-service/docs"
	_ "rider-service/docs"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	riderService ports.RiderService
	router       *gin.Engine
}

func NewHTTPHandler(riderService ports.RiderService, router *gin.Engine) *HTTPHandler {
	return &HTTPHandler{
		riderService: riderService,
		router:       router,
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
	docs.SwaggerInfo.Title = "Rider service API"
	docs.SwaggerInfo.Description = "The rider service manages all riders for the BikePack system."

	handler.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetAll godoc
// @Summary  get all riders
// @Schemes
// @Description  gets all riders in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.Rider
// @Router       /api/riders [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	riders, err := handler.riderService.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, riders)
}

// Get godoc
// @Summary  get rider
// @Schemes
// @Param        id     path  string           true  "Rider id"
// @Description  gets a rider from the system by its ID
// @Produce      json
// @Success      200  {object}  domain.Rider
// @Router       /api/riders/{id} [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	rider, err := handler.riderService.Get(uid)

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, rider)
}

// Create godoc
// @Summary  create rider
// @Schemes
// @Description  creates a new rider
// @Accept       json
// @Param        rider  body  BodyCreate  true  "Add rider"
// @Produce      json
// @Success      200  {object}  ResponseCreate
// @Router       /api/riders [post]
func (handler *HTTPHandler) Create(c *gin.Context) {
	body := dto.BodyCreate{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	rider, err := handler.riderService.Create(body.Name, body.Status)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseCreate(rider))
}

// UpdateRider godoc
// @Summary  update rider
// @Schemes
// @Description  updates a rider's information
// @Accept       json
// @Param        rider  body  BodyUpdate  true  "Update rider"
// @Param        id     path  string      true  "Rider id"
// @Produce      json
// @Success      200  {object}  ResponseUpdate
// @Router       /api/riders/{id} [put]
func (handler *HTTPHandler) UpdateRider(c *gin.Context) {
	body := dto.BodyUpdate{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	rider, err := handler.riderService.Update(uid, body.Name, body.Status)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseUpdate(rider))
}

// UpdateLocation godoc
// @Summary  update rider location
// @Schemes
// @Description  updates a rider's location
// @Accept       json
// @Param        rider  body  domain.Location  true  "Update rider"
// @Param        id  path  string  true  "Rider id"
// @Produce      json
// @Success      200  {object}  ResponseUpdate
// @Router       /api/riders/{id}/location [put]
func (handler *HTTPHandler) UpdateLocation(c *gin.Context) {
	body := domain.Location{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	uid, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	rider, err := handler.riderService.UpdateLocation(uid, body)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseUpdate(rider))
}
