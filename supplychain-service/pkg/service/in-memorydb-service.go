package service

import (
	"errors"
	"sync"
	"supplychain-service/pkg/models"
)

type InMemoryDB struct {
	data map[string]*models.CoffeeProduct
	mu   sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]*models.CoffeeProduct),
	}
}

func (db *InMemoryDB) Write(product *models.CoffeeProduct) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the productID already exists
	if _, exists := db.data[string(product.ID)]; exists {
		return errors.New("product with the same ID already exists")
	}

	// Add the product to the in-memory database
	db.data[string(product.ID)] = product
	return nil

}

func (db *InMemoryDB) Read(productID string) (*models.CoffeeProduct, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	product, exists := db.data[productID]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}