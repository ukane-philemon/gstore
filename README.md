# gstore
This a an exam for the Alt School second term exam.

John has just opened up his car selling shop, to sell different cars. He gets
the cars he needs to sell from different people and they all bring it to him. He
needs to manage the list of cars he has, attach a price to them, and put them on
display to be sold, basically John needs an inventory to manage cars & to manage
sales. For instance,

Requirements:

    1. He needs to see the number of cars that are left to be sold
    2. He needs to see the sum of the prices of the cars left
    3. He needs to see the number of cars he has sold Sum total of the prices of cars he has sold
    4. A list of orders that for the sales he made

Using the knowledge of OOP in Go, Build simple classes for the following
“objects”

Car, Product, Store

The Car class can have any car attributes you can think of.

The Product class should have attributes of a product i.e (the product, quantity
of the product in stock, price of the product). A car is a product of the store,
but there can be other products so the attribute of the product can be promoted
to the car. The Product class should have methods to display a product, and a
method to display the status of a product if it is still in stock or not.

The Store class should have attributes like:

    1. Adding an Item to the store
    2. Number of products in the store that are still up for sale
    3. Listing all product items in the store
    4. Sell an item
    5. Show a list of sold items and the total price

To run, first clone this repo using `git clone  https://github.com/ukane-philemon/gstore.git`, run `cd gstore` on your terminal, and `go run .` . Or `go install` and run `gstore` on your terminal.

**Note**: Ensure to have an Go 1.18 and above installed.

