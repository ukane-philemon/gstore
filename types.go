package main

import (
	"encoding/hex"
	"fmt"
	"time"
)

type (
	// Product is a product in a Store.
	Product interface {
		// ID returns the unique ID of the product.
		ID() productID
		// Type returns the product type.
		Type() string
		// Product returns the underlying product.
		Product() *product
		// DisplayName returns the display name of the product.
		DisplayName() string
		// Price returns the price of the product.
		Price() float64
		// Display prints information about product.
		Display()
		// Images returns a list of image urls of the product.
		Images() []string
		// IsValid checks if a product is valid and returns true if it is valid.
		IsValid() bool
	}

	// order is a buy request from a buyer.
	order struct {
		id              orderID
		name            string
		amountPaid      float64
		shippingAddress string
		products        []Product
	}
)

// productID is the unique ID of a product.
type productID [16]byte

var zeroProductID productID

func (pi productID) String() string {
	return hex.EncodeToString(pi[:])
}

func (pi productID) IsZero() bool {
	return pi == zeroProductID
}

// orderID is the unique ID of an order.
type orderID [12]byte

var zeroOrderID orderID

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
	productType    string
	category       string
	description    string
	images         []string
	specifications map[string][]string
	lastUpdated    *time.Time
	createdAt      *time.Time
}

// ID returns the unique ID of the product.
func (p *product) ID() productID {
	return p.id
}

// Type returns the product type.
func (p *product) Type() string {
	return p.productType
}

// Product returns the underlying product.
func (p *product) Product() *product {
	return p
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
	fmt.Println("Specifications:")
	for specTitle, specInfo := range p.specifications {
		fmt.Println(specTitle)
		for _, specDesc := range specInfo {
			fmt.Println(specDesc)
		}
	}
}

// Images returns a list of image urls of the product.
func (p *product) Images() []string {
	return p.images
}

// IsValid checks if a product is valid and returns true if it is valid.
func (p *product) IsValid() bool {
	return p != nil && p.name != "" && p.productType != "" && p.description != "" &&
		p.price > 0 && len(p.images) != 0 && len(p.specifications) != 0
}

// CreatedAt returns when this product was created.
func (p *product) CreatedAt() *time.Time {
	return p.createdAt
}

// LastUpdated returns the date this product was last updated.
func (p *product) LastUpdated() *time.Time {
	return p.lastUpdated
}

// car is a store product, embeddeds the product struct and re-implements
// several methods defined by the Product interface.
type car struct {
	*product
	color string
	make  string
	model string
	year  string
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

// IsValid implements part of the product interface for car.
func (c *car) IsValid() bool {
	return c.product != nil && c.product.IsValid() && c.make != "" &&
		c.model != "" && c.color != ""
}
