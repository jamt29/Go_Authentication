package models

import (
	"Go_Authentication/database"
	"time"

	"github.com/jinzhu/gorm"
)

// User representa un usuario en la base de datos
type User struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	// Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

// LoginRequest representa la estructura de datos para el inicio de sesión del usuario
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// PasswordResetToken representa un token de restablecimiento de contraseña en la base de datos
type PasswordResetToken struct {
	UserID    uint   `gorm:"primaryKey"`
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}

// CreateUser crea un nuevo usuario en la base de datos
func CreateUser(user *User) error {
	// Crear el usuario en la base de datos
	if err := database.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// VerifyUserCredentials verifica las credenciales del usuario en la base de datos
func VerifyUserCredentials(email, password string) (*User, error) {
	var user User
	// Buscar el usuario por correo electrónico en la base de datos
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err // Usuario no encontrado o error en la consulta
	}

	// Verificar la contraseña del usuario
	if user.Password != password {
		return nil, nil // Contraseña incorrecta
	}

	return &user, nil // Credenciales verificadas correctamente
}

// GetUserByEmail busca un usuario por su dirección de correo electrónico
func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err // Usuario no encontrado o error en la consulta
	}
	return &user, nil // Usuario encontrado
}

// GetUserIDByPasswordResetToken busca un token de restablecimiento de contraseña en la base de datos y devuelve el ID de usuario asociado si el token es válido
func GetUserIDByPasswordResetToken(token string) (uint, error) {
	var passwordResetToken PasswordResetToken
	if err := database.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(&passwordResetToken).Error; err != nil {
		return 0, err // Token no encontrado o expirado
	}
	return passwordResetToken.UserID, nil
}

// UpdateUserPassword actualiza la contraseña del usuario en la base de datos
func UpdateUserPassword(userID uint, newPassword string) error {
	// Buscar al usuario por su ID
	var user User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err // Error al buscar al usuario
	}

	// Actualizar la contraseña del usuario
	user.Password = newPassword
	if err := database.DB.Save(&user).Error; err != nil {
		return err // Error al guardar la contraseña actualizada
	}

	return nil // Contraseña actualizada exitosamente
}
