package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-service/internal/model"
)

type ProductServiceClient interface {
	GetProductByID(productID string) (*model.Product, error)
}

type productServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewProductServiceClient(baseURL string) ProductServiceClient {
	return &productServiceClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *productServiceClient) GetProductByID(productID string) (*model.Product, error) {
	url := fmt.Sprintf("%s/products/%s", c.baseURL, productID)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product with id %s not found", productID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get product: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var product model.Product
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, err
	}

	return &product, nil
}
