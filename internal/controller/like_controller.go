package controller

import (
	"net/http"
	"strconv"
	"time"

	likeDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/like"
	errorDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/shared/error"
	successDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/shared/success"
	"github.com/EmersonRabelo/first-api-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LikeHandler struct {
	service service.LikeService
}

func NewLikeHandler(service service.LikeService) *LikeHandler {
	return &LikeHandler{service: service}
}

func (handler *LikeHandler) Create(context *gin.Context) {
	var req likeDTO.LikeCreateDTO

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	like, err := handler.service.Create(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Não foi possivel inserir o like", Details: err.Error()})
		return
	}

	context.JSON(http.StatusCreated, like)
}

func (handler *LikeHandler) Delete(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	if err := handler.service.Delete(&id); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Err ao remover like", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, successDTO.SuccessResponse{Message: "Like removido com sucesso"})
}

func (handler *LikeHandler) FindAll(context *gin.Context) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("page_size", "1"))

	startParam := context.Query("start")
	var start time.Time

	if startParam != "" {
		t, err := time.Parse("2006-01-02", startParam)
		if err != nil {
			context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Data inicial inválida", Details: err.Error()})
			return
		}
		// Adiciona quase 1 dia, mas tira 1 nanosegundo para pegar o último instante do dia
		start = t.AddDate(0, 0, 1).Add(-time.Nanosecond)
	}

	endParam := context.Query("end")
	var end time.Time

	if endParam != "" {
		t, err := time.Parse("2006-01-02", endParam)
		if err != nil {
			context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Data final inválida", Details: err.Error()})
			return
		}
		// Adiciona quase 1 dia, mas tira 1 nanosegundo para pegar o último instante do dia
		end = t.AddDate(0, 0, 1).Add(-time.Nanosecond)
	}

	var postId *uuid.UUID = nil
	postIdParam := context.Query("post_id")

	if postIdParam != "" {
		id, err := uuid.Parse(postIdParam)
		if err != nil {
			context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{
				Error:   "Id da postagem com formato inválido",
				Details: err.Error(),
			})
			return
		}
		postId = &id
	}

	likes, err := handler.service.FindAll(postId, start, end, page, pageSize)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{
			Error:   "Erro ao buscar a lista de curtidas",
			Details: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, likes)
}

func (handler *LikeHandler) FindById(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	like, err := handler.service.FindById(&id)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Erro ocorrido ao buscar o like", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, like)
}
