package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
)

type store struct {
	mtx               sync.RWMutex
	name              string
	supportedProducts map[ProductType]bool
	productsInStock   map[productID]Product
	processedOrders   map[orderID]*order
}

// newStore creates a new store.
func newStore(name string, supportedProducts ...ProductType) *store {
	store := &store{
		name:              name,
		productsInStock:   make(map[productID]Product),
		supportedProducts: make(map[ProductType]bool),
		processedOrders:   make(map[orderID]*order),
	}

	for _, productType := range supportedProducts {
		store.supportedProducts[productType] = true
	}

	return store
}

// updateProduct adds a new product if it does not already exist or
// updates an existing product.
func (s *store) updateProduct(products ...Product) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(products) == 0 {
		return errors.New("provide one or more products")
	}

	// Validate products.
	for _, product := range products {
		if product == nil {
			return errors.New("invalid product")
		}
		if supported, ok := s.supportedProducts[product.Type()]; !ok || !supported {
			return fmt.Errorf("product with ID %s is not supported", product.ID().String())
		}
		if !product.Validate() {
			return fmt.Errorf("product with ID %s is not valid or missing required fields", product.ID().String())
		}
	}

	for _, product := range products {
		s.productsInStock[product.ID()] = product
	}

	return nil
}

// sellProduct sells one or more product to a buyer.
func (s *store) sellProduct(order *order) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if order == nil || order.id.IsZero() || order.address == "" || order.amountPaid <= 0 || order.name == "" || len(order.products) == 0 {
		return errors.New("order is missing required fields")
	}

	for _, product := range order.products {
		if product == nil {
			return errors.New("invalid product")
		}
		if supported, ok := s.supportedProducts[product.Type()]; !ok || !supported {
			return fmt.Errorf("product with ID %s is not supported", product.ID().String())
		}
		if _, ok := s.productsInStock[product.ID()]; !ok {
			return fmt.Errorf("product with ID %s does not exist", product.ID().String())
		}
		if !product.Validate() {
			return fmt.Errorf("product is not valid, please select another") // this should not happen since we validate before we added to store.
		}
	}

	for _, product := range order.products {
		delete(s.productsInStock, product.ID())
	}

	s.processedOrders[order.id] = order
	return nil
}

// product returns a single product if it is found.
func (s *store) product(ID productID) Product {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	product, ok := s.productsInStock[ID]
	if !ok {
		return nil
	}
	return product
}

// availableProducts returns the available products matching the provided
// product type, and their total cost if they are in stock. If no
// product type is specified, all the products in the store, and
// their prices are returned.
func (s *store) availableProducts(productType ProductType) ([]Product, float64, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	var products []Product
	var totalCost float64

	// If a product type is specified, validate it.
	if productType != "" {
		if supported, ok := s.supportedProducts[productType]; !ok || !supported {
			return nil, 0, fmt.Errorf("product type %s is not supported", productType)
		}

		for _, product := range s.productsInStock {
			if product.Type() == productType {
				products = append(products, product)
				totalCost += product.Price()
			}
		}
	} else {
		for _, product := range s.productsInStock {
			products = append(products, product)
			totalCost += product.Price()
		}
	}

	return products, totalCost, nil
}

// soldProducts returns the sold products matching the provided product type,
// and their total cost. If no product type is specified, all the sold products
// in the store, and their prices are returned.
func (s *store) soldProducts(productType ProductType) ([]Product, float64, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	var products []Product
	var totalCost float64

	// If not empty validate it.
	if productType != "" {
		// Ensure this product type is valid.
		if _, ok := s.supportedProducts[productType]; !ok {
			return nil, 0, fmt.Errorf("there's no product type %s", productType)
		}

		for _, orders := range s.processedOrders {
			for _, product := range orders.products {
				if product.Type() == productType {
					products = append(products, product)
					totalCost += product.Price()
				}
			}
		}
	} else {
		for _, orders := range s.processedOrders {
			for _, product := range orders.products {
				products = append(products, product)
				totalCost += product.Price()
			}
		}
	}

	return products, totalCost, nil
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
// exist. In real life we should return an error.
func (s *store) deleteProducts(productIDs ...productID) (int, error) {
	if len(productIDs) == 0 {
		return 0, errors.New("provide one or more product IDs")
	}

	s.mtx.Lock()
	var deleted int
	for _, productID := range productIDs {
		if _, ok := s.productsInStock[productID]; ok {
			delete(s.productsInStock, productID)
			deleted++
		}
	}
	s.mtx.Unlock()

	return deleted, nil
}

// inStock checks if the specified product type is in this store and
// in stock.
func (s *store) inStock(productType ProductType) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	// Validate product type and ensure there is support for it.
	if supported, ok := s.supportedProducts[productType]; !ok || !supported {
		return false
	}

	for _, product := range s.productsInStock {
		if product.Type() == productType {
			return true
		}
	}

	return false
}

// updateProductsSupported adds or disables support for a product type.
func (s *store) updateProductsSupported(productType ProductType, enable bool) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.supportedProducts[productType]
	if !enable && !ok {
		return fmt.Errorf("cannot disable non-existent product type: %s", productType)
	}

	s.supportedProducts[productType] = enable
	return nil
}

func (s *store) allSupportedProducts() map[ProductType]bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	supportedProducts := make(map[ProductType]bool, 0)
	for pType, supported := range s.supportedProducts {
		supportedProducts[pType] = supported
	}
	return supportedProducts
}

// generateProductID generates a random ID for a product.
func generateProductID(product *product) {
	_, err := rand.Read(product.id[:])
	if err != nil {
		fmt.Println(err)
	}
}

// generateOrderID generates a random ID for an order.
func generateOrderID(order *order) {
	_, err := rand.Read(order.id[:])
	if err != nil {
		fmt.Println(err)
	}
}
