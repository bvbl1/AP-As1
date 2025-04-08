package http

import (
	"Assignment1_AbylayMoldakhmet/order-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/order-service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	uc *usecase.OrderUsecase
}

func RegisterOrderRoutes(r *gin.Engine, uc *usecase.OrderUsecase) {
	handler := &OrderHandler{uc: uc}

	protected := r.Group("/orders")
	{
		protected.POST("", handler.Create)
		protected.GET("/:id", handler.GetByID)
		protected.PATCH("/:id", handler.UpdateStatus)
		protected.GET("", handler.List)
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if userID, exists := c.Get("userID"); exists {
		order.UserID = userID.(string)
	}

	if err := h.uc.Create(&order); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrNotEnoughStock {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if userID, exists := c.Get("userID"); exists && order.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		Status domain.OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	switch request.Status {
	case domain.StatusPaid, domain.StatusCancelled:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	if err := h.uc.UpdateStatus(id, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *OrderHandler) List(c *gin.Context) {
	if id, exists := c.Get("userID"); exists {
		userID := id.(string)
		orders, err := h.uc.GetByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": orders})
		return
	}

	orders, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}
