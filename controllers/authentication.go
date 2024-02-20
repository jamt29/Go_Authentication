package controllers

import (
	"Go_Authentication/models"
	"Go_Authentication/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignUpHandler maneja el registro de nuevos usuarios
func SignUpHandler(c *gin.Context) {
	// Bind JSON request to User model
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash de la contraseña del usuario utilizando bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
		return
	}

	// Reemplazar la contraseña original con el hash
	user.Password = string(hashedPassword)

	// Crear el usuario en la base de datos
	if err := models.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Usuario registrado exitosamente"})
}

// LoginHandler maneja el inicio de sesión de usuarios existentes
func LoginHandler(c *gin.Context) {
	// Bind JSON request to LoginData model
	var loginData models.LoginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el usuario por su correo electrónico
	user, err := models.GetUserByEmail(loginData.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Correo electrónico o contraseña incorrectos"})
		return
	}

	// Verificar si la contraseña proporcionada coincide con la contraseña almacenada (encriptada) en la base de datos
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Correo electrónico o contraseña incorrectos"})
		return
	}

	// Si las credenciales son válidas, generar un token JWT
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token JWT"})
		return
	}

	// Respuesta exitosa con el token JWT
	c.JSON(http.StatusOK, gin.H{"authToken": token, "message": "Inicio de sesion exitoso"})
}

// ForgotPasswordHandler maneja la solicitud de restablecimiento de contraseña
func ForgotPasswordHandler(c *gin.Context) {
	// Obtener la dirección de correo electrónico del cuerpo de la solicitud
	var requestData struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dirección de correo electrónico no válida"})
		return
	}

	// Verificar si el usuario existe
	user, err := models.GetUserByEmail(requestData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario"})
		return
	}
	if user == nil {
		// No revelar si el usuario existe o no por motivos de seguridad
		c.JSON(http.StatusOK, gin.H{"message": "Se ha enviado un correo electrónico de restablecimiento de contraseña si la dirección de correo electrónico coincide con una cuenta en nuestro sistema"})
		return
	}

	// Generar un token de restablecimiento de contraseña
	token, err := utils.GeneratePasswordResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token de restablecimiento de contraseña"})
		return
	}

	// Enviar el correo electrónico de restablecimiento de contraseña con el token
	err = utils.SendPasswordResetEmail(user.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo electrónico de restablecimiento de contraseña"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Se ha enviado un correo electrónico de restablecimiento de contraseña"})
}

// ResetPasswordHandler maneja el restablecimiento de contraseña
func ResetPasswordHandler(c *gin.Context) {
	// Obtener el token y la nueva contraseña del cuerpo de la solicitud
	var requestData struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de restablecimiento de contraseña no válidos"})
		return
	}

	// Verificar si el token de restablecimiento de contraseña es válido
	userID, err := utils.VerifyPasswordResetToken(requestData.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token de restablecimiento de contraseña inválido o expirado"})
		return
	}

	// Actualizar la contraseña del usuario
	err = models.UpdateUserPassword(userID, requestData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al restablecer la contraseña"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Contraseña restablecida exitosamente"})
}
