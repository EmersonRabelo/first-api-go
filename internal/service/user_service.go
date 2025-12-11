package service

import (
	"errors"
	"time"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/user"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Create(user *dto.CreateDTO) (*dto.UserResponseDTO, error)
	FindById(id *uuid.UUID) (*dto.UserResponseDTO, error)
	FindAll(page, pageSize int) (*dto.UserResponseListDTO, error)
	Update(id *uuid.UUID, req *dto.UpdateDTO) (*dto.UserResponseDTO, error)
	Delete(id *uuid.UUID) error
}

type userService struct {
	repository repository.UserRepository
}

func (u *userService) Create(req *dto.CreateDTO) (*dto.UserResponseDTO, error) {
	if existingUser, err := u.repository.FindByEmail(&req.Email); err != nil && existingUser != nil {
		return nil, errors.New("Email already in use.")
	}

	if existingUser, err := u.repository.FindByName(&req.Name); err != nil && existingUser != nil {
		return nil, errors.New("Name already in use.")
	}

	user := &entity.User{
		Id:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	if err := u.repository.Create(user); err != nil {
		return nil, err
	}

	return u.toUserReponseDto(user), nil

}

func (u *userService) toUserReponseDto(user *entity.User) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAT: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func (u *userService) Delete(id *uuid.UUID) error {

	if _, err := u.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Usuário não encontrado")
		}

		return err
	}

	return u.repository.Delete(id)
}

func (u *userService) FindAll(page int, pageSize int) (*dto.UserResponseListDTO, error) {
	if page < 1 {
		return nil, errors.New("Pagina deve ser maior que 1")
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, errors.New("Tamanho da página deve ser maior que 1 e menor que 100")
	}

	users, total, err := u.repository.FindAll(page, pageSize)

	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponseDTO

	for _, user := range users {
		userResponses = append(userResponses, *u.toUserReponseDto(&user))
	}

	totalPages := int(total) / pageSize
	if int(total)*pageSize > 0 {
		totalPages++
	}

	return &dto.UserResponseListDTO{
		Data:       userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (u *userService) FindById(id *uuid.UUID) (*dto.UserResponseDTO, error) {
	user, err := u.repository.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Usuário não encontrado")
		}

		return nil, err
	}

	return u.toUserReponseDto(user), nil
}

func (u *userService) Update(id *uuid.UUID, req *dto.UpdateDTO) (*dto.UserResponseDTO, error) {
	user, err := u.repository.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email {
		if existingUser, err := u.repository.FindByEmail(&req.Email); err != nil {
			return nil, err
		} else if existingUser != nil {
			return nil, errors.New("E-mail já cadastrado")
		}

		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	time := time.Now()

	if !req.IsActive {
		user.DeletedAt = &time
	}

	user.UpdatedAt = &time
	user.IsActive = req.IsActive

	if err := u.repository.Update(user); err != nil {
		return nil, err
	}

	return u.toUserReponseDto(user), nil

}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repository: repo}
}
