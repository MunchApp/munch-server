package munchserver

import (
	"munchserver/models"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", models.UserHandler)
	router.HandleFunc("/foodtrucks", models.FoodTruckHandler)
	router.HandleFunc("/reviews", models.ReviewHandler)
}
