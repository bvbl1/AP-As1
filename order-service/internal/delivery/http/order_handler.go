package http

import (
	"Assignment1_AbylayMoldakhmet/order-service/internal/domain"
	"Assignment1_AbylayMoldakhmet/order-service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	uc *usecase.OrderUsecase // изменили на указатель
}

func RegisterOrderRoutes(r *gin.Engine, uc *usecase.OrderUsecase) {
	handler := &OrderHandler{uc: uc} // передаем указатель

	// Protected routes (require JWT)
	protected := r.Group("/orders")
	// protected.Use(middleware.JwtAuthMiddleware()) // Раскомментировать после добавления middleware
	{
		protected.POST("", handler.Create)
		protected.GET("/:id", handler.GetByID)
		protected.PATCH("/:id", handler.UpdateStatus)
		protected.GET("", handler.List)
	}
}

// Create - Создание нового заказа
func (h *OrderHandler) Create(c *gin.Context) {
	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Установка userID из JWT токена (пример)
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

// GetByID - Получение заказа по ID
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Проверка прав доступа (пример)
	if userID, exists := c.Get("userID"); exists && order.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateStatus - Обновление статуса заказа
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		Status domain.OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	// Проверка допустимых статусов
	switch request.Status {
	case domain.StatusPaid, domain.StatusCancelled:
		// Valid statuses
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

// List - Получение списка заказов пользователя
func (h *OrderHandler) List(c *gin.Context) {
	userID := "" // По умолчанию
	if id, exists := c.Get("userID"); exists {
		userID = id.(string)
	}

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid user ID"})
		return
	}

	orders, err := h.uc.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}
