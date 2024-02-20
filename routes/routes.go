package routes

import (
	"Go_Authentication/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	// Rutas para autenticación
	auth := router.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUpHandler)
		auth.POST("/login", controllers.LoginHandler)

		// Rutas para el restablecimiento de contraseña
		auth.POST("/forgot-password", controllers.ForgotPasswordHandler)
		auth.POST("/reset-password", controllers.ResetPasswordHandler)
	}
}
