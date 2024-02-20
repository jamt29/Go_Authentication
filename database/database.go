package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Configurar la conexión a la base de datos PostgreSQL
	dsn := "host=localhost user=postgres password=0000 dbname=dbformulario port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error al conectar a la base de datos: " + err.Error())
	}

	// Asignar la instancia de la conexión de base de datos a la variable global DB
	DB = db
}
