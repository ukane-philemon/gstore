package main

import (
	"encoding/hex"
	"fmt"
	"time"
)

type (
	// ProductType is a type of product.
	ProductType string
	// productID is the unique product ID of a product.
	productID [32]byte
	// orderID is the unique order ID of an order.
	orderID [32]byte

	// Product is a product in a Store.
	Product interface {
		// ID returns the unique ID of the product.
		ID() productID
		// Type returns the product type.
		Type() ProductType
		// DisplayName returns the display name of the product.
		DisplayName() string
		// Description returns brief information about the product.
		Description() string
		// Price returns the price of the product.
		Price() float64
		// Category returns the category of the product.
		Category() string
		// Display prints information about  product.
		Display()
		// Images returns a list of image urls of the product.
		Images() []string
		// Validate validates a product and returns true if it is valid.
		Validate() bool
		// CreatedAt returns when this product was created.
		CreatedAt() *time.Time
		// LastUpdated returns the date this product was last updated.
		LastUpdated() *time.Time
	}

	// Vehicle is a type of product.
	Vehicle interface {
		Product
		Drive() bool
		Start() bool
	}

	// buyer can buy products from Store.
	buyer struct {
		name       string
		amountPaid float64
		address    string
	}

	// order is a sell request from a buyer with a unique ID.
	order struct {
		id orderID
		*buyer
		products []Product
	}
)

var zeroProductID productID
var zeroOrderID orderID

func (pi productID) String() string {
	return hex.EncodeToString(pi[:])
}

func (pi productID) IsZero() bool {
	return pi == zeroProductID
}

func (oi orderID) String() string {
	return hex.EncodeToString(oi[:])
}

func (oi orderID) IsZero() bool {
	return oi == zeroOrderID
}

// product implements the Product interface.
type product struct {
	id             productID
	name           string
	price          float64
	productType    ProductType
	category       string
	lastUpdated    *time.Time
	createdAt      *time.Time
	description    string
	images         []string
	specifications map[string][]string
}

// ID returns the unique ID of the product.
func (p *product) ID() productID {
	return p.id
}

// Type returns the product type.
func (p *product) Type() ProductType {
	return p.productType
}

// DisplayName returns the display name of the product.
func (p *product) DisplayName() string {
	return p.name
}

// Description returns brief information about the product.
func (p *product) Description() string {
	return p.description
}

// Price returns the price of the product.
func (p *product) Price() float64 {
	return p.price
}

// Category returns the category of the product.
func (p *product) Category() string {
	return p.category
}

// Display prints information about the product.
func (p *product) Display() {
	fmt.Println("Name: ", p.name)
	fmt.Println("Description: ", p.description)
	fmt.Println("Price: ", p.price)
}

// Images returns a list of image urls of the product.
func (p *product) Images() []string {
	return p.images
}

// Validate validates a product and returns true if it is valid.
func (p *product) Validate() bool {
	now := time.Now()
	return p != nil && p.id != zeroProductID && p.price > 0 && p.name != "" && len(p.images) != 0 && p.createdAt != nil &&
		!p.createdAt.After(now) && p.lastUpdated != nil &&
		!p.lastUpdated.After(now) && p.productType != "" && len(p.specifications) != 0
}

// CreatedAt returns when this product was created.
func (p *product) CreatedAt() *time.Time {
	return p.createdAt
}

// LastUpdated returns the date this product was last updated.
func (p *product) LastUpdated() *time.Time {
	return p.lastUpdated
}

type car struct {
	*product
	color string
	make  string
	model string
}

// Drive implements part of the vehicle interface for car.
func (c *car) Drive() bool {
	fmt.Printf("You are driving %s\n", c.DisplayName())
	return true
}

// // Starts implements part of the vehicle interface for car.
func (c *car) Start() bool {
	fmt.Printf("Your %s car has started\n", c.DisplayName())
	return true
}

// Display implements part of the Product interface for car.
func (c *car) Display() {
	fmt.Println("Name: ", c.DisplayName())
	fmt.Println("Make and Model: ", c.make, c.model)
	fmt.Println("Specifications:")
	for specTitle, specInfo := range c.specifications {
		fmt.Println(specTitle)
		for _, specDesc := range specInfo {
			fmt.Println(specDesc)
		}
	}
}

// Validate implements part of the product interface for car.
func (c *car) Validate() bool {
	return c.product != nil && c.product.Validate() && c.make != "" &&
		c.model != "" && c.color != ""
}

// carAccessory is a product that implements the Product interface.
type carAccessory struct {
	*product
	specifications map[string][]string
}

// Display implements part of the Product interface for carAccessory.
func (c *carAccessory) Display() {
	fmt.Println("Name: ", c.DisplayName())
	fmt.Println("Specifications:")
	for specTitle, specInfo := range c.specifications {
		fmt.Println(specTitle)
		for _, specDesc := range specInfo {
			fmt.Println(specDesc)
		}
	}
}

// Validate implements part of the product interface for carAccessory.
func (c *carAccessory) Validate() bool {
	return c.product != nil && c.product.Validate()
}
