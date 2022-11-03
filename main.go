package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// These are the supported product type for our Auto-Shop.
	productTypeCar := ProductType("Car")
	productTypeCarAccessory := ProductType("Car Accessory")
	now := time.Now()

	// newStore creates a store that can sell different products.
	autoShop := newStore("Auto Shop", productTypeCar, productTypeCarAccessory)
	item1 := &car{
		product: &product{
			name:        "Ford Ecosport",
			price:       5000000,
			productType: productTypeCar,
			category:    "Used Cars",
			lastUpdated: &now,
			createdAt:   &now,
			description: "The EcoSport is easy to drive and spacious inside. The 1.0-litre petrol engine is a popular choice because of its efficiency.",
			images:      []string{"https://uks-cdn.pinewooddms.com/b04b90f8-2e99-463d-a023-7e3c771fb388/vehicles/1935a96a-3bb8-485e-affc-132707e733c1.jpg?", "https://uks-cdn.pinewooddms.com/b04b90f8-2e99-463d-a023-7e3c771fb388/vehicles/4cb99337-5c1b-4f0e-9bb7-3683f23520de.jpg?"},
			specifications: map[string][]string{
				"Key Features": {"Bluetooth", "Climate Control", "Air Conditioning", "Ask for a Test Drive Today", "24 Month Guarantee Available", "2 x Keys with car"},
				"Engine":       {"Auto", "Petrol"},
			},
		},
		color: "yellow",
		make:  "Ford",
		model: "1.5 Zetec 5dr 2016",
	}
	generateProductID(item1.product)

	item2 := &car{
		product: &product{
			name:        "Ford Ecosport",
			price:       5000000,
			productType: productTypeCar,
			category:    "Used Cars",
			lastUpdated: &now,
			createdAt:   &now,
			description: "The EcoSport is easy to drive and spacious inside. The 1.0-litre petrol engine is a popular choice because of its efficiency.",
			images:      []string{"https://uks-cdn.pinewooddms.com/b04b90f8-2e99-463d-a023-7e3c771fb388/vehicles/1935a96a-3bb8-485e-affc-132707e733c1.jpg?", "https://uks-cdn.pinewooddms.com/b04b90f8-2e99-463d-a023-7e3c771fb388/vehicles/4cb99337-5c1b-4f0e-9bb7-3683f23520de.jpg?"},
			specifications: map[string][]string{
				"Key Features": {"Bluetooth", "Climate Control", "Air Conditioning", "Ask for a Test Drive Today", "24 Month Guarantee Available", "2 x Keys with car"},
				"Engine":       {"Auto", "Petrol"},
			},
		},
		color: "black",
		make:  "Ford",
		model: "1.5 Zetec 5dr 2016",
	}
	generateProductID(item2.product)

	item3 := &carAccessory{
		product: &product{
			name:        "Toyota Shadow Logo Led Light (For 4 Doors)",
			price:       14000,
			productType: productTypeCarAccessory,
			category:    "Led Lights",
			lastUpdated: &now,
			createdAt:   &now,
			description: "TOYOTA LED HOLOGRAM SAFETY LIGHTS(free batteries included): Stay safe at night when stepping out of your cars in poorly lit areas with our classy, elegant light emitting diode car door lights.",
			images:      []string{"https://ng.jumia.is/unsafe/fit-in/500x500/filters:fill(white)/product/74/552546/1.jpg?6525"},
		},
		specifications: map[string][]string{
			"Key Features": {"Toyota LED Hologram Safety Lights, Free batteries included"},
		},
	}
	generateProductID(item3.product)

	// Add different supported products to the store.
	// Store Feature 1.
	err := autoShop.updateProduct(item1, item2, item3)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Store Feature 2 and 3.
	allAvailableProducts, totalCost, err := autoShop.availableProducts("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s has %d products available that cost a total of %.2f NGN\n", autoShop.name, len(allAvailableProducts), totalCost)

	// Store feature 4.
	order := &order{
		buyer: &buyer{
			name:       "Philemon",
			amountPaid: item1.price + item3.price,
			address:    "No 21 Alt_School Africa street, Banana Island, Lagos",
		},
		products: []Product{item1, item3},
	}
	generateOrderID(order)

	err = autoShop.sellProduct(order)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Store Feature 5 and Requirement 3.
	allSoldProducts, totalCost, err := autoShop.soldProducts("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s has sold %d products at %.2f NGN\n", autoShop.name, len(allSoldProducts), totalCost)

	allSoldCars, totalCost, err := autoShop.soldProducts(productTypeCar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s has sold %d %s at %.2f NGN\n", autoShop.name, len(allSoldCars), productTypeCar, totalCost)

	// Requirement 1 and 2.
	allAvailableCars, totalCost, err := autoShop.availableProducts(productTypeCar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s has %d %s available that cost a total of %.2f NGN\n", autoShop.name, len(allAvailableCars), productTypeCar, totalCost)

	processedOrders, totalPaid := autoShop.orders()
	fmt.Printf("%s has processed %d orders totalling %2.f NGN\n", autoShop.name, len(processedOrders), totalPaid)

	// Check that products are in stock.
	inStock := autoShop.inStock(productTypeCar)
	fmt.Printf("%s has a %s in stock: %v\n", autoShop.name, productTypeCar, inStock)

	inStock = autoShop.inStock(productTypeCarAccessory)
	fmt.Printf("%s has a %s in stock: %v\n", autoShop.name, productTypeCarAccessory, inStock)

	// Check product availability.
	product := autoShop.product(item1.id)
	fmt.Printf("Product with id %s is available: %v\n", item1.id.String(), product != nil)

	// Update store settings.
	err = autoShop.updateProductsSupported("Spare Parts", true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = autoShop.updateProductsSupported(productTypeCarAccessory, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Retrieve supported product kind and there status.
	supportedProducts := autoShop.allSupportedProducts()
	for pType, isSupported := range supportedProducts {
		fmt.Printf("%s currently has support for %s: %v\n", autoShop.name, pType, isSupported)
	}

	// Delete products from store.
	deleted, err := autoShop.deleteProducts(item1.id, item2.id, item3.id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Deleted %d products from %s\n", deleted, autoShop.name)
}
