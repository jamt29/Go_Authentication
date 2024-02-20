package main

import (
	"Go_Authentication/database"
	"Go_Authentication/models"
	"Go_Authentication/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conectar a la base de datos
	database.Connect()
	// Auto-migrar el esquema de la base de datos
	database.DB.AutoMigrate(&models.User{})

	// Crear una instancia del router Gin
	router := gin.Default()

	// Configurar las rutas
	routes.Setup(router)

	// Iniciar el servidor Gin
	router.Run(":8080")
}
