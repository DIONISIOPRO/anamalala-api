package middlewares

import (
	"net/http"

	"github.com/anamalala/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminMiddlewares struct {
	tokenUtil utils.TokenUtil
}

func NewAdminMiddleware(tokenUtil utils.TokenUtil) AdminMiddlewares {
	return AdminMiddlewares{
		tokenUtil: tokenUtil,
	}
}

// AdminMiddleware verifica se o usuário é um administrador
func (md *AdminMiddlewares) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o usuário está autenticado
		if !isAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			c.Abort()
			return
		}

		// Obter o role do usuário
		userRole, exists := getUserRole(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role não encontrado no token"})
			c.Abort()
			return
		}

		// Verificar se o usuário é um administrador
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "acesso restrito a administradores"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// BannedUserMiddleware verifica se o usuário está banido e impede o acesso
func (md *AdminMiddlewares) BannedUserMiddleware(userService interface {
	IsBanned(userID string) (bool, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o usuário está autenticado
		if !isAuthenticated(c) {
			c.Next()
			return
		}

		// Obter o ID do usuário
		userID, exists := getUserID(c)
		if !exists {
			c.Next()
			return
		}

		// Verificar se o usuário está banido
		banned, err := userService.IsBanned(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao verificar status do usuário"})
			c.Abort()
			return
		}

		if banned {
			c.JSON(http.StatusForbidden, gin.H{"error": "usuário está banido"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ContentOwnerOrAdminMiddleware verifica se o usuário é o proprietário do conteúdo ou um administrador
func (md *AdminMiddlewares) ContentOwnerOrAdminMiddleware(getOwnerIDFunc func(contentID string) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o usuário está autenticado
		if !isAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
			c.Abort()
			return
		}

		// Obter o ID do usuário e o role
		userID, _ := getUserID(c)
		userRole, _ := getUserRole(c)

		// Se o usuário for administrador, permitir acesso
		if userRole == "admin" {
			c.Next()
			return
		}

		// Obter o ID do conteúdo (postagem, comentário, etc.)
		contentID := c.Param("id")
		if contentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conteúdo não fornecido"})
			c.Abort()
			return
		}

		// Obter o ID do proprietário do conteúdo
		ownerID, err := getOwnerIDFunc(contentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao verificar propriedade do conteúdo"})
			c.Abort()
			return
		}

		// Verificar se o usuário é o proprietário
		if userID != ownerID {
			c.JSON(http.StatusForbidden, gin.H{"error": "não autorizado a modificar este conteúdo"})
			c.Abort()
			return
		}

		c.Next()
	}
}
