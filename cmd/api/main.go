package main

import (
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org
// @host localhost:8081 // Replace with your actual host
// @BasePath /
func main() {
	loadconfig, locaconfigerr := configs.LoadConfig()
	if locaconfigerr != nil {
		panic(locaconfigerr)
	}
	app, err := InitializeApp(loadconfig)
	// Initialize Gin router
	if err != nil {
		panic(err)
	}
	log.Println("Starting gRPC server on :50051")
	if err := app.GRPCServer.Run(":50051"); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	port := fmt.Sprintf(":%v", loadconfig.Server.Port)
	app.Server.Run(port)

}
