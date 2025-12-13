package controller

import (
	"net/http"
	"strconv"

	errorDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/error"
	likeDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/like"
	successDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/success"
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

	likes, err := handler.service.FindAll(page, pageSize)
	if err != nil {
		context.JSON(http.StatusNoContent, errorDTO.ErrorResponse{Error: "Erro ao buscar a lista de curtidas", Details: err.Error()})
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
		context.JSON(http.StatusNoContent, errorDTO.ErrorResponse{Error: "Erro ocorrido ao buscar o like", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, like)
}
