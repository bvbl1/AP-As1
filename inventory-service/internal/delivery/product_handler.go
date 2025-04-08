package delivery

import (
	"Assignment1_AbylayMoldakhmet/inventory-service/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase domain.ProductUsecase
}

func NewProductHandler(r *gin.Engine, uc domain.ProductUsecase) {
	handler := &ProductHandler{usecase: uc}

	products := r.Group("/products")
	{
		products.POST("", handler.Create)
		products.GET("/:id", handler.GetByID)
		products.PATCH("/:id", handler.Update)
		products.DELETE("/:id", handler.Delete)
		products.GET("", handler.List)
	}
	r.GET("/api/products/:id/check-stock", handler.CheckStock)

}

func (h *ProductHandler) Create(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.usecase.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var updated domain.Product
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	updated.ID = domain.StringToObjectID(id)

	if err := h.usecase.Update(&updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.usecase.List(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CheckStock(c *gin.Context) {
	productID := c.Param("id")
	quantity, _ := strconv.Atoi(c.Query("quantity"))

	available, err := h.usecase.CheckStock(productID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": available})
}
