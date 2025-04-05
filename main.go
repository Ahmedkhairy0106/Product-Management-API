package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"Product-Management-API-v3/db"
	"Product-Management-API-v3/handlers"
	"Product-Management-API-v3/limiting"
)

func main() {
    // Initialize the database connection
    _, err := db.InitializeDB()
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.DB.Close()

    ch := chi.NewRouter()
    ch.Use(middleware.Logger)

    store := limiting.NewRateLimiterStore(1, 5)
    ch.Use(limiting.RateLimiterMiddleware(store))

	ch.Route("/api", func (api chi.Router){
		api.Get("/products", handlers.GetAllProduct)
		api.Get("/products/{id}", handlers.GetProductByID)

		api.Post("/products", handlers.AddProduct)

		api.Put("/products/{id}", handlers.UpdateProduct)
		api.Delete("/products/{id}", handlers.DeleteProduct)
	})

    log.Println("Server running on http://localhost:8081")
    log.Fatal(http.ListenAndServe(":8081", ch))
}