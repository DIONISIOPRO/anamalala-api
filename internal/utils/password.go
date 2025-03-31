package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword gera um hash bcrypt para uma senha
func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("a senha deve ter pelo menos 6 caracteres")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPasswordHash verifica se a senha corresponde ao hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateResetCode gera um código de reset de 6 dígitos
func GenerateResetCode() string {
	// Gerar número aleatório entre 100000 e 999999
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	n = n.Add(n, big.NewInt(100000))
	
	return fmt.Sprintf("%d", n)
}

// GenerateRandomString gera uma string aleatória de comprimento específico
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	return hex.EncodeToString(bytes), nil
}
