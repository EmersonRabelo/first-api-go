package service

import (
	"errors"
	"time"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/like"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikeService interface {
	Create(like *dto.LikeCreateDTO) (*dto.LikeResponseDTO, error)
	FindById(id *uuid.UUID) (*dto.LikeResponseDTO, error)
	FindAll(page, pageSize int) (*dto.LikeResponseListDTO, error)
	Update(id *uuid.UUID, req *dto.LikeUpdateDTO) (*dto.LikeResponseDTO, error)
	Delete(id *uuid.UUID) error
}

type likeService struct {
	repository  repository.LikeRepository
	userService UserService
	postService PostService
}

func (l *likeService) Create(req *dto.LikeCreateDTO) (*dto.LikeResponseDTO, error) {

	if _, err := l.userService.FindById(&req.UserId); err != nil {
		return nil, errors.New("Usuário não econtrado")
	}

	if _, err := l.postService.FindById(&req.PostId); err != nil {
		return nil, errors.New("Postagem não econtrada")
	}

	like := &entity.Like{
		Id:        uuid.New(),
		UserId:    req.UserId,
		PostId:    req.PostId,
		Quantity:  0,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	if err := l.repository.Create(like); err != nil {
		return nil, err
	}

	return l.toPostResponse(like), nil
}

func (l *likeService) Delete(id *uuid.UUID) error {
	if _, err := l.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Like não encontrado")
		}

		return err
	}

	return l.repository.Delete(id)
}

func (l *likeService) FindAll(page int, pageSize int) (*dto.LikeResponseListDTO, error) {
	if page < 1 {
		return nil, errors.New("Pagina deve ser maior que 1")
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, errors.New("Tamanho da página deve ser maior que 1 e menor que 100")
	}

	likes, total, err := l.repository.FindAll(page, pageSize)

	if err != nil {
		return nil, err
	}

	var likeResponse []dto.LikeResponseDTO

	for _, like := range likes {
		likeResponse = append(likeResponse, *l.toPostResponse(&like))
	}

	totalPages := int(total) / pageSize

	if int(total)*pageSize > 0 {
		totalPages++
	}

	return &dto.LikeResponseListDTO{
		Data:       likeResponse,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}, nil

}

func (l *likeService) FindById(id *uuid.UUID) (*dto.LikeResponseDTO, error) {
	like, err := l.repository.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Like não encontrado")
		}

		return nil, err
	}

	return l.toPostResponse(like), nil
}

func (l *likeService) Update(id *uuid.UUID, req *dto.LikeUpdateDTO) (*dto.LikeResponseDTO, error) {
	return nil, errors.New("Não disponivel para uso")
}

func (l *likeService) toPostResponse(like *entity.Like) *dto.LikeResponseDTO {
	return &dto.LikeResponseDTO{
		Id:        like.Id,
		UserId:    like.UserId,
		PostId:    like.PostId,
		Quantity:  like.Quantity,
		CreatedAt: like.CreatedAt,
		UpdatedAt: like.UpdatedAt,
		DeletedAt: like.DeletedAt,
	}
}

func NewLikeService(repository repository.LikeRepository, userService UserService, postService PostService) LikeService {
	return &likeService{repository: repository, userService: userService, postService: postService}
}
