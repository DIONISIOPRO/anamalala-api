package handlers

import (
	"net/http"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/services"
	"github.com/anamalala/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
	validator   utils.Validator
}

func NewAuthHandler(authService services.AuthService, validator utils.Validator) AuthHandler {
	return AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, "Dados de registro inválidos,"+err.Error())
		return
	}
	user.Contact = h.validator.FormatPhoneNumber(user.Contact)

	// Validar campos obrigatórios
	if user.Name == "" || user.Province == "" || user.Contact == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, "Todos os campos são obrigatórios")
		return
	}

	ok := h.validator.ValidatePhoneNumber(user.Contact)

	if !ok {
		c.JSON(http.StatusBadRequest, "Numero invalido")
		return
	}

	err := h.validator.ValidatePassword(user.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Palavra passe invalida")
		return
	}

	ok = h.validator.ValidateProvince(user.Province)

	if !ok {
		c.JSON(http.StatusBadRequest, "Provincia invalida")
		return
	}
	saveuser := models.User{}

	saveuser.Contact = user.Contact
	saveuser.Password = user.Password
	saveuser.Name = user.Name
	saveuser.Province = user.Province

	registeredUser, err := h.authService.Register(c, saveuser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao registrar usuário")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Usuário registrado com sucesso",
		"data": gin.H{
			"user": registeredUser,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest struct {
		Contact  string `json:"contact" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, "Dados de login inválidos")
		return
	}

	loginRequest.Contact = h.validator.FormatPhoneNumber(loginRequest.Contact)

	user, token, err := h.authService.Login(c, loginRequest.Contact, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Credenciais inválidas")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": user,
	})
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"contact" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Número de telefone inválido")
		return
	}
	// Gerar token de redefinição e enviar SMS
	err := h.authService.ResetPasswordRequest(c, request.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao enviar token de redefinição")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token de redefinição enviado com sucesso",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"contact" binding:"required"`
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	err := h.authService.ResetPassword(c, request.PhoneNumber, request.Token, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Falha ao redefinir senha")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Senha redefinida com sucesso",
	})
}

func (h *AuthHandler) VerifyResetPasswordToken(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"contact" binding:"required"`
		Token       string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	ok := h.authService.VerifyResetPasswordCode(c, request.PhoneNumber, request.Token)
	if !ok {
		c.JSON(http.StatusBadRequest, "codigo invalido")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"token": request.Token,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var request struct{
		Phone string `json:"contact"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Número de telefone inválido")
		return
	}

	err := h.authService.Logout(c, request.Phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} 
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout realizado com sucesso",
	})
}
