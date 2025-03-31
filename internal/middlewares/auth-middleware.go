package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/anamalala/internal/services"
	"github.com/anamalala/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthMiddlewares struct {
	tokenUtil   utils.TokenUtil
	userservice services.UserService
}

func NewAuthMiddleware(tokenUtil utils.TokenUtil, userservice services.UserService) AuthMiddlewares {
	return AuthMiddlewares{
		tokenUtil:   tokenUtil,
		userservice: userservice,
	}
}

// AuthMiddleware verifica se o token JWT é válido e adiciona o ID do usuário e o role ao contexto
func (m *AuthMiddlewares) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Query("token")
		if authHeader == "" {
			authHeader = c.GetHeader("Authorization")
		}
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorização não fornecido"})
			c.Abort()
			return
		}
		// Verificar formato do token
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar token
		claims, err := m.tokenUtil.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido: " + err.Error()})
			c.Abort()
			return
		}

		user, err := m.userservice.GetUserByID(c, claims.UserID)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario inválido: " + err.Error()})
			c.Abort()
			return
		}

		if !user.IsLoggedIn {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario nao autenticado, entre de novo"})
			c.Abort()
			return
		}

		// Adicionar ID do usuário e role ao contexto
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		// Adicionar ao contexto para uso nos serviços
		ctx := context.WithValue(c.Request.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// OptionalAuthMiddleware verifica se há um token JWT, mas não rejeita a requisição se não houver
func (m *AuthMiddlewares) OptionalAuthMiddleware(tokenUtil *utils.TokenUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Sem token, continuar sem adicionar informações de usuário
			c.Next()
			return
		}

		// Verificar formato do token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Token mal formatado, continuar sem adicionar informações de usuário
			c.Next()
			return
		}

		tokenString := parts[1]

		// Validar token
		claims, err := tokenUtil.ValidateToken(tokenString)
		if err != nil {
			// Token inválido, continuar sem adicionar informações de usuário
			c.Next()
			return
		}

		// Adicionar ID do usuário e role ao contexto
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		// Adicionar ao contexto para uso nos serviços
		ctx := context.WithValue(c.Request.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetUserID retorna o ID do usuário do contexto
func getUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}

	return userID.(string), true
}

// GetUserRole retorna o role do usuário do contexto
func getUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("userRole")
	if !exists {
		return "", false
	}

	return userRole.(string), true
}

// IsAuthenticated verifica se o usuário está autenticado
func isAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("userID")
	return exists
}
