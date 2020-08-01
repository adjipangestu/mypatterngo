package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"mypatterngo/routes"
	"net/http"
	"os"
)

var server = Server{}

func Run()  {
	Routes := routes.Api{}

	// Check env data
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load env")
		panic(err)
	}

	// Init the database
	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	// Init the routes
	Routes.InitializeRoutes()


	//Seed DB
	//seed.Load(server.DB)

	//With Origin
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://local.detik.com:3000","http://local.detik.com:4000"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodPut},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "API-Key"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	handler := c.Handler(Routes.Router)

	// Run the server
	fmt.Println("Connected To Database")
	fmt.Println("Server started port 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
