package utils

import (
	"Go_Authentication/models"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateJWTToken genera un token JWT con los claims proporcionados
func GenerateJWTToken(userID uint) (string, error) {
	// Definir los claims del token JWT
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expira en 24 horas
		// Puedes agregar otros claims según tus necesidades, como el nombre de usuario, roles, etc.
	}

	// Generar el token JWT con los claims y la firma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con una clave secreta y obtener la representación de cadena del token
	authToken, err := token.SignedString([]byte("tu_clave_secreta"))
	if err != nil {
		return "", err // Error al firmar el token
	}

	return authToken, nil // Token JWT generado correctamente
}

// GeneratePasswordResetToken genera un token único para restablecer la contraseña
func GeneratePasswordResetToken() (string, error) {
	// Generar un token aleatorio seguro
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err // Error al generar el token
	}

	// Convertir los bytes aleatorios en una cadena base64
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	return token, nil
}

// SendPasswordResetEmail envía un correo electrónico al usuario con el enlace de restablecimiento de contraseña
func SendPasswordResetEmail(email, token string) error {
	// Configurar los detalles del servidor de correo electrónico de Gmail
	smtpServer := "smtp.gmail.com"
	smtpPort := "587" // Puerto TLS para Gmail

	// Tu dirección de correo electrónico de Gmail y contraseña
	senderEmail := "tucorreo@gmail.com"
	senderPassword := "tupassword"

	// Construir el mensaje de correo electrónico
	to := []string{email}
	subject := "Restablecimiento de contraseña"
	body := "Hola,\n\nPara restablecer tu contraseña, haz clic en el siguiente enlace:\n\nhttp://tuapp.com/reset-password?token=" + token + "\n\nSi no solicitaste un restablecimiento de contraseña, ignora este correo electrónico.\n\nSaludos,\nTu aplicación"

	// Autenticarse en el servidor de correo electrónico de Gmail
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	// Enviar el correo electrónico utilizando TLS
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, to, []byte("Subject:"+subject+"\r\n\r\n"+body))
	if err != nil {
		return err // Error al enviar el correo electrónico
	}

	return nil
}

// VerifyPasswordResetToken verifica si un token de restablecimiento de contraseña es válido
func VerifyPasswordResetToken(token string) (uint, error) {
	// Verificar si el token existe en la base de datos y si aún no ha expirado
	userID, err := models.GetUserIDByPasswordResetToken(token)
	if err != nil {
		return 0, err // Error al buscar el token
	}
	if userID == 0 {
		return 0, errors.New("Token de restablecimiento de contraseña inválido o expirado")
	}
	return userID, nil
}
