package clients

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type InventoryClient struct {
	baseURL string
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{baseURL: baseURL}
}

type InventoryClientInterface interface {
	CheckStock(productID string, quantity int) (bool, error)
}

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
