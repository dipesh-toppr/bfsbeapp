package main

import (
	"fmt"

	"github.com/dipesh-toppr/bfsbeapp/config"
	"github.com/dipesh-toppr/bfsbeapp/routes"
)

// Application starts here.
func main() {
	db := config.Database

	routes.LoadRoutes()

	db.Close()
	fmt.Printf("Database close")
}
