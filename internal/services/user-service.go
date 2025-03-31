package services

import (
	"context"
	"errors"
	"time"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/repositories/interfaces"
	"github.com/anamalala/internal/utils"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) UserService {
	return UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	
	// Não retornar a senha e o código de reset
	user.Password = ""
	user.ResetCode = ""
	
	return user, nil
}


func (s *UserService) GetUserByContact(ctx context.Context, contact string) (models.User, error) {
	user, err := s.userRepo.FindByContact(ctx, contact)
	if err != nil {
		return models.User{}, err
	}
	
	// Não retornar a senha e o código de reset
	user.Password = ""
	user.ResetCode = ""
	return user, nil
}


func (s *UserService) UpdateUser(ctx context.Context, id string, updateData models.User) (models.User, error) {

	// Obter usuário atual
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	
	// Atualizar campos
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	
	if updateData.Province != "" {
		user.Province = updateData.Province
	}
	
	// Verificar se o contato está sendo alterado e se já está em uso
	if updateData.Contact != "" && updateData.Contact != user.Contact {
		existingUser, _ := s.userRepo.FindByContact(ctx, updateData.Contact)
		if existingUser.ID != "" {
			return models.User{}, errors.New("este contato já está em uso")
		}
		user.Contact = updateData.Contact
	}
	
	// Atualizar senha se fornecida
	if updateData.Password != "" {
		hashedPassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			return models.User{}, err
		}
		user.Password = hashedPassword
	}
	
	user.UpdatedAt = time.Now()
	
	// Salvar alterações
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	
	// Não retornar a senha e o código de reset
	user.Password = ""
	user.ResetCode = ""
	
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, page, limit int) (models.Users, int, error) {
	users, total, err := s.userRepo.List(ctx, int64(page), int64(limit))
	if err != nil {
		return nil, 0, err
	}
	
	// Remover informações sensíveis
	for _, user := range users {
		user.Password = ""
		user.ResetCode = ""
	}
	
	return users, int(total), nil
}

func (s *UserService) GetUsersByProvince(ctx context.Context, province string, page, limit int) (models.Users, int, error) {
	users, total, err := s.userRepo.ListByProvince(ctx, province, int64(page), int64(limit))
	if err != nil {
		return nil, 0, err
	}
	
	// Remover informações sensíveis
	for i, user := range users {
		user.Password = ""
		user.ResetCode = ""
		users[i] = user
	}
	
	return users, int(total), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
