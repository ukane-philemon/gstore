package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// store is the keeps track of all the existing and sold products.
type store struct {
	name            string
	mtx             sync.RWMutex
	products        map[productID]Product
	processedOrders map[orderID]*order
}

// newStore creates a new store.
func newStore(name string) *store {
	store := &store{
		name:            name,
		products:        make(map[productID]Product),
		processedOrders: make(map[orderID]*order),
	}

	return store
}

// addProducts adds new product(s) and returns an array of product IDs.
func (s *store) addProducts(products ...Product) ([]productID, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(products) == 0 {
		return nil, errors.New("provide one or more products")
	}

	// Validate products.
	for _, product := range products {
		if product == nil {
			return nil, errors.New("invalid product")
		}

		if !product.IsValid() {
			return nil, fmt.Errorf("product with ID %s is not valid or missing required fields", product.ID().String())
		}
	}

	now := time.Now()
	productIDs := make([]productID, len(products))
	for i, p := range products {
		product := p.Product()

		// Generate a new ID for this product.
		s.generateProductID(product)

		// Set essential product dates.
		product.createdAt = &now
		product.lastUpdated = &now

		// Add product to store products map and also add the product ID to
		// return to callers.
		productID := p.ID()
		s.products[productID] = p
		productIDs[i] = productID
	}

	return productIDs, nil
}

// sellProduct sells one or more product to a buyer and returns the order ID.
func (s *store) sellProduct(order *order) (orderID, error) {
	if order == nil || order.shippingAddress == "" || order.amountPaid <= 0 || order.name == "" || len(order.products) == 0 {
		return zeroOrderID, errors.New("order is missing required fields")
	}

	var totalProductCost float64
	for _, p := range order.products {
		if p == nil {
			return zeroOrderID, errors.New("invalid product")
		}

		if _, ok := s.products[p.ID()]; !ok {
			return zeroOrderID, fmt.Errorf("product with ID %s does not exist", p.ID().String())
		}

		if !p.IsValid() {
			return zeroOrderID, fmt.Errorf("product with ID(%s) is not valid", p.ID())
		}

		totalProductCost += p.Price()
	}

	// Check if buyer paid enough.
	if order.amountPaid < totalProductCost {
		return zeroOrderID, fmt.Errorf("order amount paid is not enough, need %f but paid %f", totalProductCost, order.amountPaid)
	}

	s.mtx.Lock()
	for _, p := range order.products {
		delete(s.products, p.ID())
	}

	// Generate new order ID.
	s.generateOrderID(order)
	s.processedOrders[order.id] = order
	s.mtx.Unlock()

	return order.id, nil
}

// product returns a single product if it is found.
func (s *store) product(ID productID) Product {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	product, ok := s.products[ID]
	if !ok {
		return nil
	}
	return product
}

// availableProducts returns the available products matching the provided
// product type, and their total cost if they are in stock. If no product type
// is specified, all the products in the store, and their prices are returned.
func (s *store) availableProducts(productType string) ([]Product, float64) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	var products []Product
	var totalCost float64

	if productType == "" {
		for _, product := range s.products {
			products = append(products, product)
			totalCost += product.Price()
		}
		return products, totalCost
	}

	for _, product := range s.products {
		if product.Type() == productType {
			products = append(products, product)
			totalCost += product.Price()
		}
	}

	return products, totalCost
}

// soldProducts returns the sold products matching the provided product type,
// and their total cost. If no product type is specified, all the sold products
// in the store, and their prices are returned.
func (s *store) soldProducts(productType string) ([]Product, float64) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	var products []Product
	var totalCost float64

	if productType == "" {
		for _, orders := range s.processedOrders {
			for _, product := range orders.products {
				products = append(products, product)
				totalCost += product.Price()
			}
		}
		return products, totalCost
	}

	for _, orders := range s.processedOrders {
		for _, product := range orders.products {
			if product.Type() == productType {
				products = append(products, product)
				totalCost += product.Price()
			}
		}
	}

	return products, totalCost
}

// orders returns a list of processed orders.
func (s *store) orders() ([]*order, float64) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	var orders []*order
	var totalPaid float64
	for _, order := range s.processedOrders {
		orders = append(orders, order)
		totalPaid += order.amountPaid
	}
	return orders, totalPaid
}

// deleteProducts removes one or more available product from the store and
// return the number of products deleted. It will be a no-op if product does not
// exist.
func (s *store) deleteProducts(productIDs ...productID) (int, error) {
	if len(productIDs) == 0 {
		return 0, errors.New("provide one or more product IDs")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	var deleted int
	for _, productID := range productIDs {
		if _, ok := s.products[productID]; ok {
			delete(s.products, productID)
			deleted++
		}
	}

	return deleted, nil
}

// inStock checks if the specified product type is in this store and
// in stock.
func (s *store) inStock(productType string) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, product := range s.products {
		if product.Type() == productType {
			return true
		}
	}

	return false
}

// generateProductID generates a random ID for a product.
func (s *store) generateProductID(product *product) {
	_, err := rand.Read(product.id[:])
	if err != nil {
		log.Println(err)
	}
}

// generateOrderID generates a random ID for an order.
func (s *store) generateOrderID(product *order) {
	_, err := rand.Read(product.id[:])
	if err != nil {
		log.Println(err)
	}
}
