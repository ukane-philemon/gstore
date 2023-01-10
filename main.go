package main

import (
	"fmt"
	"os"
)

func main() {
	autoShopSimulation()
}

/*
autoShopSimulation runs a simulation of the various functionalities implemented
by Store. Simulation of expected features and requirements are not in any
specific order.

User Story:

Build an auto shop inventory with the following features.

Requirements:
1. Shop owner needs to see the number of cars that are left to be sold
2. Shop owner needs to see the sum of the prices of the cars left
3. Shop owner needs to see the number of cars he has sold
4. Shop owner needs to see the sum total of the prices of cars sold
5. Shop owner needs to see a list of orders that for the sales made

The Store class should have attributes like:
1. Adding an Item to the store
2. Number of products in the store that are still up for sale
3. Listing all product items in the store
4. Sell an item
5. Show a list of sold items and the total price
*/
func autoShopSimulation() {
	// These are the supported product type for our Auto-Shop.
	productTypeCar, productTypeCarAccessory := "Car", "Car Accessory"

	// newStore creates a store that can sell different products. All product
	// prices in this store are denominated in the Nigerian Naira.
	autoShop := newStore("Auto Shop")

	item1 := &car{
		product: &product{
			name:        "Ford Ecosport",
			price:       5000000,
			productType: productTypeCar,
			category:    "Used Cars",
			description: "The EcoSport is easy to drive and spacious inside. The 1.0-litre petrol engine is a popular choice because of its efficiency.",
			images:      []string{"https://uks-cdn.pinewooddms.com/b04b90f8-2e99-463d-a023-7e3c771fb388/vehicles/1935a96a-3bb8-485e-affc-132707e733c1.jpg?", "https://uks-cdn.pinewooddms.com/b04b90f8-2e99-463d-a023-7e3c771fb388/vehicles/4cb99337-5c1b-4f0e-9bb7-3683f23520de.jpg?"},
			specifications: map[string][]string{
				"Key Features": {"Bluetooth", "Climate Control", "Air Conditioning", "Ask for a Test Drive Today", "24 Month Guarantee Available", "2 x Keys with car"},
				"Engine":       {"Auto", "Petrol"},
			},
		},
		color: "yellow",
		make:  "Ford",
		model: "1.5 Zetec 5dr",
		year:  "2016",
	}

	item2 := &car{
		product: &product{
			name:        "Honda HR-V SPORT",
			price:       7000000,
			productType: productTypeCar,
			category:    "Used Cars",
			description: "The Honda HR-V SPORT easy to drive and spacious inside. The automatic engine is a popular choice because of its efficiency.",
			images:      []string{"https://content.homenetiol.com/698/2163991/1920x1080/8ac0270d04d344b1ad58ae18e01c4c88.jpg", "https://content.homenetiol.com/698/2163991/1920x1080/ae3d1b14b4614451938dd3703a18222a.jpg"},
			specifications: map[string][]string{
				"Key Features": {"Bluetooth", "Cruise Control", "4 Doors", "Rear Defroster", "Climate Control", "Air Conditioning", "Ask for a Test Drive Today", "24 Month Guarantee Available", "2 x Keys with car"},
				"Engine":       {"Auto", "Petrol", "4 Cylinders 1.8L"},
			},
		},
		color: "black",
		make:  "Honda",
		model: "4 Cylinders 1.8L",
		year:  "2018",
	}

	item3 := &product{
		name:        "Toyota Shadow Logo Led Light (For 4 Doors)",
		price:       14000,
		productType: productTypeCarAccessory,
		category:    "Led Lights",
		description: "TOYOTA LED HOLOGRAM SAFETY LIGHTS(free batteries included): Stay safe at night when stepping out of your cars in poorly lit areas with our classy, elegant light emitting diode car door lights.",
		images:      []string{"https://ng.jumia.is/unsafe/fit-in/500x500/filters:fill(white)/product/74/552546/1.jpg?6525"},
		specifications: map[string][]string{
			"Key Features": {"Toyota LED Hologram Safety Lights, Free batteries included"},
		},
	}

	// Add different supported products to the store.
	// Store Feature 1.
	productIDs, err := autoShop.addProducts(item1, item2, item3)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, id := range productIDs {
		fmt.Printf("Successfully added product with ID(%s) to %s\n", id, autoShop.name)
	}

	// Store Feature 2 and 3.
	// Retrieve information for all products in the store.
	allAvailableProducts, totalCost := autoShop.availableProducts("")
	fmt.Printf("%s has %d products available that cost a total of %.2f NGN\n", autoShop.name, len(allAvailableProducts), totalCost)

	// Retrieve information for a specific product kind in the store.
	allAvailableProducts, totalCost = autoShop.availableProducts(productTypeCar)
	fmt.Printf("%s has %d %s's available that cost a total of %.2f NGN\n", autoShop.name, len(allAvailableProducts), productTypeCar, totalCost)

	// Store feature 4.
	order := &order{
		name:            "Philemon",
		amountPaid:      item1.price + item3.price,
		shippingAddress: "No 21 Alt_School Africa street, Banana Island, Lagos",
		products:        []Product{item1, item3},
	}

	orderID, err := autoShop.sellProduct(order)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s has processed order with ID(%s) successfully\n", autoShop.name, orderID)

	// Store Feature 5.
	allSoldProducts, totalCost := autoShop.soldProducts("")
	fmt.Printf("%s has sold a total of %d products for %.2f NGN\n", autoShop.name, len(allSoldProducts), totalCost)

	// Requirement 3 and 4.
	allSoldCars, totalCost := autoShop.soldProducts(productTypeCar)
	fmt.Printf("%s has sold %d %s for %.2f NGN\n", autoShop.name, len(allSoldCars), productTypeCar, totalCost)

	// Requirement 1 and 2.
	allAvailableCars, totalCost := autoShop.availableProducts(productTypeCar)
	fmt.Printf("%s has %d %s available that cost a total of %.2f NGN\n", autoShop.name, len(allAvailableCars), productTypeCar, totalCost)

	// Shop feature 5 and Requirement 5.
	processedOrders, totalPaid := autoShop.orders()
	fmt.Printf("%s has processed %d orders totalling %2.f NGN\n", autoShop.name, len(processedOrders), totalPaid)

	// Check that products are in stock.
	inStock := autoShop.inStock(productTypeCar)
	fmt.Printf("%s has a %s in stock: %v\n", autoShop.name, productTypeCar, inStock)

	inStock = autoShop.inStock(productTypeCarAccessory)
	fmt.Printf("%s has a %s in stock: %v\n", autoShop.name, productTypeCarAccessory, inStock)

	// Check product availability.
	product := autoShop.product(item1.id)
	fmt.Printf("Sold product with id %s is available: %v\n", item1.id, product != nil)

	// Delete products from store.
	deleted, err := autoShop.deleteProducts(productIDs...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Deleted %d product(s) from %s\n", deleted, autoShop.name)
}
