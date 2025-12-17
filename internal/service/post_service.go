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
	FindAll(id *uuid.UUID, start, end time.Time, page, pageSize int) (*dto.PostResponseListDTO, error)
	Update(id *uuid.UUID, req *dto.PostUpdateDTO) (*dto.PostResponseDTO, error)
	Delete(id *uuid.UUID) error
}

type postService struct {
	repository  repository.PostRepository
	userService UserService
}

func (p *postService) Create(req *dto.PostCreateDTO) (*dto.PostResponseDTO, error) {
	user, err := p.userService.FindById(&req.UserId)
	if err != nil {
		return nil, errors.New("Usuário não encontrado")
	}

	post := &entity.Post{
		Id:        uuid.New(),
		UserId:    user.Id,
		Body:      req.Body,
		IsActive:  true,
		CreatedAt: time.Now(),
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

func (p *postService) FindAll(id *uuid.UUID, start, end time.Time, page, pageSize int) (*dto.PostResponseListDTO, error) {

	if page < 1 {
		return nil, errors.New("Pagina deve ser maior que 1")
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, errors.New("Tamanho da página deve ser maior que 1 e menor que 100")
	}

	if id != nil {
		if _, err := p.userService.FindById(id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Usuário não encontrado")
			}
			return nil, err
		}
	}

	posts, total, err := p.repository.FindAll(id, start, end, page, pageSize)

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

	now := time.Now()

	if !req.IsActive {
		post.DeletedAt.Time = now
	}

	post.UpdatedAt = now
	post.IsActive = req.IsActive

	if err := p.repository.Update(post); err != nil {
		return nil, err
	}

	return p.toPostResponse(post), nil

}

func (p *postService) toPostResponse(post *entity.Post) *dto.PostResponseDTO {
	var deletedAt *time.Time
	if post.DeletedAt.Valid {
		deletedAt = &post.DeletedAt.Time
	}

	var updatedAt *time.Time
	if !post.UpdatedAt.IsZero() {
		updatedAt = &post.UpdatedAt
	}
	return &dto.PostResponseDTO{
		Id:        post.Id,
		UserId:    post.UserId,
		Body:      post.Body,
		IsActive:  post.IsActive,
		CreatedAt: post.CreatedAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func NewPostService(repository repository.PostRepository, userService UserService) PostService {
	return &postService{
		repository:  repository,
		userService: userService,
	}
}
