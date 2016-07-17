package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Burger struct {
	Toppings []string `json:"toppings"`
	Price    int      `json:"price"`
	Calories int      `json:"calories"`
}

type Burgers []Burger

type BurgersResponse struct {
	Burgers Burgers `json:"burgers"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	burgers := BurgersResponse{
		Burgers: Burgers{
			Burger{
				Toppings: []string{"cheddar cheese", "lettuce", "mushrooms"},
				Price:    100,
				Calories: 1000,
			},
			Burger{
				Toppings: []string{"bacon", "peanut butter"},
				Price:    50,
				Calories: 500,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(burgers)
}

func main() {
	http.HandleFunc("/", indexHandler)

	fmt.Println("Listening on port :8080")
	http.ListenAndServe(":8080", nil)
}
