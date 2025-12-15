package service

import (
	"errors"
	"time"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/reply"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReplyService interface {
	Create(reply *dto.ReplyCreateDTO) (*dto.ReplyResponseDTO, error)
	FindById(id *uuid.UUID) (*dto.ReplyResponseDTO, error)
	FindAll(page, pageSize int) (*dto.ReplyResponseListDTO, error)
	Update(id *uuid.UUID, req *dto.ReplyUpdateDTO) (*dto.ReplyResponseDTO, error)
	Delete(id *uuid.UUID) error
}

type replyService struct {
	repository  repository.ReplyRepository
	userService UserService
	postService PostService
}

func (r *replyService) Create(req *dto.ReplyCreateDTO) (*dto.ReplyResponseDTO, error) {

	if _, err := r.userService.FindById(&req.UserId); err != nil {
		return nil, errors.New("Usuário não encontrado")
	}

	if _, err := r.postService.FindById(&req.PostId); err != nil {
		return nil, errors.New("Postagem não encontrada")
	}

	reply := &entity.Reply{
		Id:        uuid.New(),
		UserId:    req.UserId,
		PostId:    req.PostId,
		Body:      req.Body,
		Quantity:  0,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	if err := r.repository.Create(reply); err != nil {
		return nil, err
	}

	return r.toReplyResponse(reply), nil
}

func (r *replyService) Delete(id *uuid.UUID) error {
	if _, err := r.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Like não encontrado")
		}
		return err
	}

	return r.repository.Delete(id)
}

func (r *replyService) FindAll(page int, pageSize int) (*dto.ReplyResponseListDTO, error) {
	if page < 1 {
		return nil, errors.New("Pagina deve ser maior que 1")
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, errors.New("Tamanho da página deve ser maior que 1 e menor que 100")
	}

	replies, total, err := r.repository.FindAll(page, pageSize)

	if err != nil {
		return nil, err
	}

	var replyResponse []dto.ReplyResponseDTO

	for _, reply := range replies {
		replyResponse = append(replyResponse, *r.toReplyResponse(&reply))
	}

	totalPages := int(total) / pageSize

	if int(total)*pageSize > 0 {
		totalPages++
	}

	return &dto.ReplyResponseListDTO{
		Data:       replyResponse,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}, nil

}

func (r *replyService) FindById(id *uuid.UUID) (*dto.ReplyResponseDTO, error) {
	reply, err := r.repository.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Like não encontrado")
		}

		return nil, err
	}

	return r.toReplyResponse(reply), nil

}

func (r *replyService) Update(id *uuid.UUID, req *dto.ReplyUpdateDTO) (*dto.ReplyResponseDTO, error) {
	reply, err := r.repository.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Post não encontrado")
		}
		return nil, err
	}

	if req.Body != "" {
		reply.Body = req.Body
	}

	time := time.Now()

	reply.UpdatedAt = &time

	if err := r.repository.Update(reply); err != nil {
		return nil, err
	}

	return r.toReplyResponse(reply), nil
}

func (r *replyService) toReplyResponse(reply *entity.Reply) *dto.ReplyResponseDTO {
	return &dto.ReplyResponseDTO{
		Id:        reply.Id,
		UserId:    reply.UserId,
		PostId:    reply.PostId,
		Body:      reply.Body,
		Quantity:  reply.Quantity,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
		DeletedAt: reply.DeletedAt,
	}
}

func NewReplyService(repository repository.ReplyRepository, userService UserService, postService PostService) ReplyService {
	return &replyService{
		repository:  repository,
		userService: userService,
		postService: postService,
	}
}
