package handlers

import (
	"net/http"
	"services-api/logger"
	"services-api/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandlers struct {
	repo   *repository.ServiceRepository
	logger *logger.Logger
}

func NewServiceHandlers(repo *repository.ServiceRepository, logger *logger.Logger) *ServiceHandlers {
	return &ServiceHandlers{repo: repo, logger: logger}
}

func RegisterServiceHandlers(r *gin.Engine, repo *repository.ServiceRepository, logger *logger.Logger) {
	h := NewServiceHandlers(repo, logger)
	r.GET("/services", h.GetServices)
	r.GET("/services/:id", h.GetService)
	r.POST("/services", h.CreateService)
}

func (h *ServiceHandlers) GetServices(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Query("name")
	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	services, err := h.repo.GetServices(ctx, name, sort, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func (h *ServiceHandlers) GetService(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	service, err := h.repo.GetService(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandlers) CreateService(c *gin.Context) {
	ctx := c.Request.Context()
	var request struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Versions    []string `json:"versions"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error(ctx, "Invalid request payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.repo.CreateServiceRaw(ctx, request.Name, request.Description, request.Versions); err != nil {
		h.logger.Error(ctx, "Failed to create service", "name", request.Name, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info(ctx, "Service created", "name", request.Name)
	c.JSON(http.StatusCreated, gin.H{"name": request.Name})
}
