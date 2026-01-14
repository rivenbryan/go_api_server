package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type dollars float32             // Float in Dollars
type database map[string]dollars // A database is a map of string keys and values dollars

// Two methods for database. One is list and another one is price
func (db database) list(w http.ResponseWriter, r *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %0.2f\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%0.2f\n", price)
}

func (db database) create(w http.ResponseWriter, r *http.Request) {
	// Extract out the key and ID from request
	item := r.URL.Query().Get("item")
	priceStr := r.URL.Query().Get("price")

	if item == "" || priceStr == "" {
		http.Error(w, "item and price are required", http.StatusBadRequest)
		return
	}

	if _, ok := db[item]; ok {
		http.Error(w, "item already exists", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		http.Error(w, "invalid price", http.StatusBadRequest)
		return
	}
	db[item] = dollars(price)

}

func (db database) update(w http.ResponseWriter, r *http.Request) {
	// Extract out the key and ID from request
	item := r.URL.Query().Get("item")
	priceStr := r.URL.Query().Get("price")

	if item == "" || priceStr == "" {
		http.Error(w, "item and price are required", http.StatusBadRequest)
		return
	}

	if _, ok := db[item]; !ok {
		http.Error(w, "item does not exists", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		http.Error(w, "invalid price", http.StatusBadRequest)
		return
	}

	db[item] = dollars(price)
}

func (db database) delete(w http.ResponseWriter, r *http.Request) {
	// Extract out the key and ID from request
	item := r.URL.Query().Get("item")

	if item == "" {
		http.Error(w, "item is required", http.StatusBadRequest)
		return
	}

	if _, ok := db[item]; !ok {
		http.Error(w, "item does not exists", http.StatusBadRequest)
		return
	}

	delete(db, item)
}
func main() {
	fmt.Println("Hello, World!")
	db := database{ // Creates a database type called DB
		"shoes": 50,
		"socks": 5,
	}

	// db.list and db.price are two methods of function type
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
