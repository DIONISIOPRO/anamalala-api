package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type TokenUtil struct {
	secretKey []byte
	expiresIn time.Duration
}

func NewTokenUtil(secretKey string, expiresIn time.Duration) TokenUtil {
	return TokenUtil{
		secretKey: []byte(secretKey),
		expiresIn: expiresIn,
	}
}

// GenerateToken cria um novo token JWT para o usuário
func (t TokenUtil) GenerateToken(userID, role string) (string, error) {
	// Definir tempo de expiração

	// Criar claims

	// Criar token com claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,                           // Subject (user identifier)
		"iss": role,                             // Issuer
		"aud": role,                             // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),
	})

	// Assinar token com chave secreta
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida um token JWT e retorna as claims
func (t TokenUtil) ValidateToken(tokenString string) (Claims, error) {
	claims := Claims{}
	// Parse do token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return t.secretKey, nil
	})
	claims.UserID, _ = token.Claims.GetSubject()
	claims.Role, _ = token.Claims.GetIssuer()
	return claims, nil
}
