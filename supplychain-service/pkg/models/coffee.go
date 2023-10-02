package models

import (
    "time"
)

// CoffeeProduct represents a coffee bean product.
type CoffeeProduct struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    State       string    `json:"state"`
    CreatedAt   time.Time `json:"created_at"`
}

func CreateCoffeeProduct(name, description string, price float64) *CoffeeProduct {
	return &CoffeeProduct{
		Name:        name,
		Description: description,
		Price:       price,
	}
}