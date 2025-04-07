package clients

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// InventoryClient - структура для работы с Inventory Service
type InventoryClient struct {
	baseURL string
}

// NewInventoryClient создает новый клиент для взаимодействия с Inventory Service
func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{baseURL: baseURL}
}

// InventoryClientInterface интерфейс для работы с сервисом инвентаря
type InventoryClientInterface interface {
	CheckStock(productID string, quantity int) (bool, error)
}

// CheckStock проверяет наличие товара на складе
func (c *InventoryClient) CheckStock(productID string, quantity int) (bool, error) {
	resp, err := http.Get(c.baseURL + "/api/products/" + productID + "/check-stock?quantity=" + strconv.Itoa(quantity))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Available bool `json:"available"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Available, nil
}
