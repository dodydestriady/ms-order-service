package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-service/config"
	"order-service/internal/model"
	"os"
)

type ProductServiceClient interface {
	GetProductByID(productID string) (*model.Product, error)
}

type productServiceClient struct {
	baseURL     string
	client      *http.Client
	mockProduct bool
}

func NewProductServiceClient(baseURL string, cfg *config.Config) ProductServiceClient {
	return &productServiceClient{
		baseURL:     baseURL,
		client:      &http.Client{},
		mockProduct: cfg.MockProduct,
	}
}

func (c *productServiceClient) GetProductByID(productID string) (*model.Product, error) {
	if c.mockProduct {
		fmt.Println("Using mock product data from JSON file!")

		mockData, err := os.ReadFile("product-response.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read mock product file: %w", err)
		}

		var mockProduct model.Product
		if err := json.Unmarshal(mockData, &mockProduct); err != nil {
			return nil, fmt.Errorf("failed to parse mock product JSON: %w", err)
		}

		return &mockProduct, nil
	}

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
