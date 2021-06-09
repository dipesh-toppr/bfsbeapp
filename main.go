package main

import (
	"fmt"

	"github.com/dipesh-toppr/bfsbeapp/managers"
	"github.com/dipesh-toppr/bfsbeapp/routes"
)

// Application starts here.
func main() {
	db := managers.Database

	routes.LoadRoutes()

	db.Close()
	fmt.Printf("Database close")
}
