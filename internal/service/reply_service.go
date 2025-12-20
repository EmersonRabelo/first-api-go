package service

import (
	"errors"
	"fmt"
	"time"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/reply"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	redisService "github.com/EmersonRabelo/first-api-go/internal/redis"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReplyService interface {
	Create(reply *dto.ReplyCreateDTO) (*dto.ReplyResponseDTO, error)
	FindById(id *uuid.UUID) (*dto.ReplyResponseDTO, error)
	FindAll(postID, userID *uuid.UUID, start, end time.Time, page int, pageSize int) (*dto.ReplyResponseListDTO, error)
	Update(id *uuid.UUID, req *dto.ReplyUpdateDTO) (*dto.ReplyResponseDTO, error)
	Delete(id *uuid.UUID) error
	incrementLike(id *uuid.UUID) (uint64, error)
	decrementLike(id *uuid.UUID) (uint64, error)
	setLike(id *uuid.UUID, value uint64) error
}

type replyService struct {
	repository  repository.ReplyRepository
	userService UserService
	postService PostService
	redisClient *redis.Client
}

func (r *replyService) Create(req *dto.ReplyCreateDTO) (*dto.ReplyResponseDTO, error) {

	if _, err := r.userService.FindById(&req.UserId); err != nil {
		return nil, errors.New("Usuário não encontrado")
	}

	if _, err := r.postService.FindById(&req.PostId); err != nil {
		return nil, errors.New("Postagem não encontrada")
	}

	quantity, err := r.incrementLike(&req.PostId)

	if err != nil {
		likesQuantity, err := r.repository.GetRepliesCountByPostID(&req.PostId)

		if err != nil {
			return nil, err
		}

		quantity = likesQuantity
		if err := r.setLike(&req.PostId, likesQuantity); err != nil {
			fmt.Println("Não foi possivel sincronizar o redis, ", err)
		}
	}

	reply := &entity.Reply{
		Id:        uuid.New(),
		UserId:    req.UserId,
		PostId:    req.PostId,
		Body:      req.Body,
		Quantity:  quantity,
		CreatedAt: time.Now(),
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

func (r *replyService) FindAll(postID, userID *uuid.UUID, start, end time.Time, page int, pageSize int) (*dto.ReplyResponseListDTO, error) {
	if page < 1 {
		return nil, errors.New("Pagina deve ser maior que 1")
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, errors.New("Tamanho da página deve ser maior que 1 e menor que 100")
	}

	if postID != nil {
		if _, err := r.postService.FindById(postID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Postagem não encontrada")
			}
			return nil, err
		}
	}

	if userID != nil {
		if _, err := r.userService.FindById(userID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Usuário não encontrado")
			}
			return nil, err
		}
	}

	replies, total, err := r.repository.FindAll(postID, userID, start, end, page, pageSize)

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

	reply.UpdatedAt = time

	if err := r.repository.Update(reply); err != nil {
		return nil, err
	}

	return r.toReplyResponse(reply), nil
}

func (r *replyService) toReplyResponse(reply *entity.Reply) *dto.ReplyResponseDTO {
	var deletedAt *time.Time
	if reply.DeletedAt.Valid {
		deletedAt = &reply.DeletedAt.Time
	}

	var updatedAt *time.Time
	if !reply.UpdatedAt.IsZero() {
		updatedAt = &reply.UpdatedAt
	}
	return &dto.ReplyResponseDTO{
		Id:        reply.Id,
		UserId:    reply.UserId,
		PostId:    reply.PostId,
		Body:      reply.Body,
		Quantity:  reply.Quantity,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func (l *replyService) incrementLike(id *uuid.UUID) (uint64, error) {
	count, err := redisService.IncrementCounter(l.redisClient, "reply:post", id.String())
	if err != nil {
		return 0, err
	}
	if count < 0 {
		// retorno erro caso o contador esteja negativo, para evitar overflow.
		return 0, fmt.Errorf("contador negativo: %d", count)
	}
	return uint64(count), nil
}

func (l *replyService) setLike(id *uuid.UUID, value uint64) error {
	if err := redisService.SetCounter(l.redisClient, "reply:post", id.String(), value); err != nil {
		return err
	}

	return nil
}

func (l *replyService) decrementLike(id *uuid.UUID) (uint64, error) {
	count, err := redisService.DecrementCounter(l.redisClient, "reply:post", id.String())
	if err != nil {
		return 0, err
	}
	if count < 0 {
		// retorno erro caso o contador esteja negativo, para evitar overflow.
		return 0, fmt.Errorf("contador negativo: %d", count)
	}
	return uint64(count), nil
}

func NewReplyService(repository repository.ReplyRepository, userService UserService, postService PostService, redisClient *redis.Client) ReplyService {
	return &replyService{
		repository:  repository,
		userService: userService,
		postService: postService,
		redisClient: redisClient,
	}
}
