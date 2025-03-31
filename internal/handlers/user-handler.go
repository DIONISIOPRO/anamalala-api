package handlers

import (
	"net/http"
	"strconv"

	"github.com/anamalala/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	user, err := h.userService.GetUserByID(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar perfil")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Perfil obtido com sucesso",
		"data":    user,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	var updateData struct {
		Name     string `json:"name"`
		Province string `json:"province"`
		Contact  string `json:"contact"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	user, err := h.userService.GetUserByID(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar usuário")
		return
	}

	// Atualizar apenas os campos fornecidos
	if updateData.Name != "" {
		user.Name = updateData.Name
	}

	if updateData.Password != "" {
		user.Password = updateData.Password
	}

	if updateData.Contact != "" {
		user.Contact = updateData.Contact
	}
	if updateData.Province != "" {
		user.Province = updateData.Province
	}


	updatedUser, err := h.userService.UpdateUser(c, user.ID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao atualizar perfil")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Perfil atualizado com sucesso",
		"data":    updatedUser,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, "ID do usuário não fornecido")
		return
	}

	user, err := h.userService.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Usuário não encontrado")
		return
	}

	// Remover informações sensíveis para visualização pública
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuário obtido com sucesso",
		"data":    user,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	// Parâmetros para paginação
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil {
		limit = 10
	}
	// Filtros opcionais

	users, total, err := h.userService.GetAllUsers(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar usuários")
		return
	}

	// Remover senhas da resposta
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuários obtidos com sucesso",
		"data": gin.H{
			"users":      users,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}
func (h *UserHandler) GetTotalOnline(c *gin.Context) {
	users, total, err := h.userService.GetAllUsers(c, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar usuários")
		return
	}

	// Remover senhas da resposta
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuários obtidos com sucesso",
		"data": gin.H{
			"total": total,
		},
	})
}

func (h *UserHandler) GetUsersByProvince(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	province := c.Param("province")
	if province == "" {
		c.JSON(http.StatusBadRequest, "Província não fornecida")
		return
	}

	// Parâmetros para paginação
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil {
		limit = 10
	}

	users, total, err := h.userService.GetUsersByProvince(c, province, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar usuários por província")
		return
	}

	// Remover senhas da resposta
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuários obtidos com sucesso",
		"data": gin.H{
			"province":   province,
			"users":      users,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}
