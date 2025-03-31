package handlers

import (
	"net/http"
	"strconv"

	"github.com/anamalala/internal/services"
	"github.com/anamalala/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) AdminHandler {
	return AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		utils.ForbiddenResponse(c, "Somente para Administradores")
		return
	}

	stats, err := h.adminService.GetSystemStats(c)

	if err != nil {
		utils.InternalServerErrorResponse(c, "Falha ao obter estatísticas")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Estatísticas obtidas com sucesso",
		"data":    stats,
	})
}

func (h *AdminHandler) PromoteToAdmin(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		utils.ForbiddenResponse(c, "Acesso restrito a administradores")
		return
	}

	userID := c.Param("id")
	if userID == "" {
		utils.BadRequestResponse(c, "ID do usuário não fornecido")
		return
	}
	// Validar role

	err := h.adminService.PromoteToAdmin(c, userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Falha ao atualizar função do usuário")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Função do usuário atualizada com sucesso",
		"data": gin.H{
			"user_id": userID,
			"role":    "admin",
		},
	})
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		utils.ForbiddenResponse(c, "Acesso restrito a administradores")
		return
	}

	userID := c.Param("id")
	if userID == "" {
		utils.BadRequestResponse(c, "ID do usuário não fornecido")
		return
	}


	err := h.adminService.BanUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao banir usuário")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuário banido com sucesso",
		"data": gin.H{
			"user_id": userID,
		},
	})
}

func (h *AdminHandler) UnbanUser(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, "ID do usuário não fornecido")
		return
	}

	err := h.adminService.UnbanUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao desbanir usuário")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuário desbanido com sucesso",
		"data": gin.H{
			"user_id": userID,
		},
	})
}

func (h *AdminHandler) GetBannedUsers(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	// Parâmetros para paginação

	 page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if  err != nil{
		page = 1
	 }
	limit,_ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil{
		limit = 10
	 }

	users, total, err := h.adminService.GetBannedUsers(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar usuários banidos")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuários banidos obtidos com sucesso",
		"data": gin.H{
			"users":      users,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

func (h *AdminHandler) GetAdminLogs(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	// Parâmetros para paginação
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if  err != nil{
		page = 1
	 }
	limit,_ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil{
		limit = 20
	 }


	// Filtros opcionais
	action := c.Query("action")

	logs, err := h.adminService.GetAdminLogs(page, limit, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar logs administrativos",)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logs obtidos com sucesso",
		"data": gin.H{
			"logs":       logs,
			"total":      len(logs),
			"page":       page,
			"limit":      limit,
			"totalPages": (len(logs) + limit - 1) / limit,
		},
	})
}

func (h *AdminHandler) SendMassMessage(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden,"Acesso restrito a administradores")
		return
	}

	var request struct {
		Message   string   `json:"message" binding:"required"`
		Provinces []string `json:"provinces"` // Se vazio, envia para todos
		Title     string   `json:"title"`
	}

	var err error

	if err := c.ShouldBindJSON(&request); err != nil {
		 c.JSON(http.StatusBadRequest,"Dados inválidos")
		return
	}

	if len(request.Provinces) == 0{
		_, err = h.adminService.SendSMSToAllUsers(c, request.Message)
	}

	for _, province := range request.Provinces{
		_, err = h.adminService.SendSMSToProvince(c, request.Message, province)
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao enviar mensagem em massa")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mensagem em massa enviada com sucesso",

	})
}

func (h *AdminHandler) DeleteAccount(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, "ID do usuário não fornecido")
		return
	}

	var request struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	err := h.adminService.DeleteUserAccount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,"Falha ao excluir conta")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Conta excluída com sucesso",
	})
}
