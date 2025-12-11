package service

import (
	"errors"
	"time"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/post"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService interface {
	Create(post *dto.PostCreateDTO) (*dto.PostResponseDTO, error)
	FindById(id *uuid.UUID) (*dto.PostResponseDTO, error)
	FindAll(page, pageSize int) (*dto.PostResponseListDTO, error)
	Update(id *uuid.UUID, req *dto.PostUpdateDTO) (*dto.PostResponseDTO, error)
	Delete(id *uuid.UUID) error
}

type postService struct {
	repository  repository.PostRepository
	userService *userService
}

func (p *postService) Create(req *dto.PostCreateDTO) (*dto.PostResponseDTO, error) {
	user, err := p.userService.repository.FindById(&req.UserId)
	if err != nil {
		return nil, errors.New("Usuário não encontrado")
	}

	post := &entity.Post{
		Id:        uuid.New(),
		UserId:    user.Id,
		User:      *user,
		Body:      req.Body,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	if err := p.repository.Create(post); err != nil {
		return nil, err
	}

	return p.toPostResponse(post), nil
}

func (p *postService) Delete(id *uuid.UUID) error {
	if _, err := p.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Post não encontrado")
		}
		return err
	}

	return p.repository.Delete(id)
}

func (p *postService) FindAll(page int, pageSize int) (*dto.PostResponseListDTO, error) {
	if page < 1 {
		return nil, errors.New("Pagina deve ser maior que 1")
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, errors.New("Tamanho da página deve ser maior que 1 e menor que 100")
	}

	posts, total, err := p.repository.FindAll(page, pageSize)

	if err != nil {
		return nil, err
	}

	var postResponse []dto.PostResponseDTO

	for _, post := range posts {
		postResponse = append(postResponse, *p.toPostResponse(&post))
	}

	totalPage := int(total) / pageSize

	if int(total)*pageSize > 0 {
		totalPage++
	}

	return &dto.PostResponseListDTO{
		Data:       postResponse,
		Total:      total,
		Page:       page,
		TotalPages: totalPage,
	}, nil
}

func (p *postService) FindById(id *uuid.UUID) (*dto.PostResponseDTO, error) {
	post, err := p.repository.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Post não encontrado")
		}

		return nil, err
	}

	return p.toPostResponse(post), err
}

func (p *postService) Update(id *uuid.UUID, req *dto.PostUpdateDTO) (*dto.PostResponseDTO, error) {
	post, err := p.repository.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Post não encontrado")
		}

		return nil, err
	}

	if req.Body != "" {
		post.Body = req.Body
	}

	time := time.Now()

	if !req.IsActive {
		post.DeletedAt = &time
	}

	post.UpdatedAt = &time
	post.IsActive = req.IsActive

	if err := p.repository.Update(post); err != nil {
		return nil, err
	}

	return p.toPostResponse(post), nil

}

func (p *postService) toPostResponse(post *entity.Post) *dto.PostResponseDTO {
	return &dto.PostResponseDTO{
		Id:        post.Id,
		UserId:    post.UserId,
		Body:      post.Body,
		IsActive:  post.IsActive,
		CreatedAt: post.CreatedAt,
		UpdatedAt: *post.UpdatedAt,
		DeletedAt: *post.DeletedAt,
	}
}

func NewPostService(repository repository.PostRepository) PostService {
	return &postService{repository: repository}
}
