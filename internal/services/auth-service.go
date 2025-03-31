package services

import (
	"context"
	"errors"

	"time"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/repositories/interfaces"
	"github.com/anamalala/internal/utils"
)

type AuthService struct {
	userRepo  interfaces.UserRepository
	tokenUtil utils.TokenUtil
	//sms *sms.Service
}

func (s AuthService) VerifyResetPasswordCode(ctx context.Context, contact string, token string) bool {

	user, err := s.userRepo.FindByContact(ctx, contact)
	if err != nil {
		return false
	}

	// Verificar se o código existe e está válido
	if user.ResetCode != token {
		return false
	}

	// Verificar se o código não expirou
	if time.Now().After(user.ResetCodeExpiry) {
		return false
	}

	return true

}

func NewAuthService(userRepo interfaces.UserRepository, tokenUtil utils.TokenUtil) AuthService {
	return AuthService{
		userRepo:  userRepo,
		tokenUtil: tokenUtil,
	}
}

func (s *AuthService) Register(ctx context.Context, user models.User) (models.User, error) {
	// Verificar se o usuário já existe
	_, err := s.userRepo.FindByContact(ctx, user.Contact)
	if err == nil {
		return models.User{}, errors.New("user esxists")
	}
	// Hash da senha
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	user.Password = hashedPassword

	user.Role = "user" // Por padrão, todos os novos registros são usuários normais
	// Salvar usuário
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, contact, password string) (models.User, string, error) {
	// Buscar usuário pelo contacto
	user, err := s.userRepo.FindByContact(ctx, contact)
	if err != nil {
		return models.User{}, "", errors.New("credenciais inválidas")
	}

	// Verificar senha
	if !utils.CheckPasswordHash(password, user.Password) {
		return models.User{}, "", errors.New("credenciais inválidas")
	}

	// Gerar token JWT
	token, err := s.tokenUtil.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return models.User{}, "", err
	}
	user.IsLoggedIn = true
	er := s.userRepo.Update(ctx, user)
	user.Password = ""
	if er != nil {
		return models.User{}, "", err
	}
	return user, token, nil
}

func (s *AuthService) Logout(ctx context.Context, contact string) error {
	// Buscar usuário pelo contacto
	user, err := s.userRepo.FindByContact(ctx, contact)
	if err != nil {
		return errors.New("credenciais inválidas: ")
	}
	user.IsLoggedIn = false
	er := s.userRepo.Update(ctx, user)
	if er != nil {
		return er
	}
	return nil
}

func (s *AuthService) ResetPasswordRequest(ctx context.Context, contact string) error {
	// Verificar se o usuário existe
	user, err := s.userRepo.FindByContact(ctx, contact)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Gerar código de recuperação (6 dígitos)
	resetCode := utils.GenerateResetCode()

	// Armazenar código no usuário
	user.ResetCode = resetCode
	user.ResetCodeExpiry = time.Now().Add(15 * time.Minute) // Expira em 15 minutos

	// Atualizar usuário
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Enviar SMS com código de recuperação (mock)
	// Na implementação real, integrar com serviço de SMS

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, contact, resetCode, newPassword string) error {
	// Buscar usuário
	user, err := s.userRepo.FindByContact(ctx, contact)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Verificar se o código existe e está válido
	if user.ResetCode != resetCode {
		return errors.New("código de recuperação inválido")
	}

	// Verificar se o código não expirou
	if time.Now().After(user.ResetCodeExpiry) {
		return errors.New("código de recuperação expirado")
	}

	// Hash da nova senha
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Atualizar senha e limpar códigos de recuperação
	user.Password = hashedPassword
	user.ResetCode = ""
	user.ResetCodeExpiry = time.Time{}
	user.UpdatedAt = time.Now()

	// Salvar usuário
	return s.userRepo.Update(ctx, user)
}
