package controller

import (
	"net/http"
	"strconv"

	errorDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/error"
	replyDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/reply"
	successDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/success"
	"github.com/EmersonRabelo/first-api-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReplyHandler struct {
	service service.ReplyService
}

func NewReplyHandler(service service.ReplyService) *ReplyHandler {
	return &ReplyHandler{service: service}
}

func (handler *ReplyHandler) Create(context *gin.Context) {
	var req replyDTO.ReplyCreateDTO

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	like, err := handler.service.Create(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Não foi possivel inserir o reply", Details: err.Error()})
		return
	}

	context.JSON(http.StatusCreated, like)
}

func (handler *ReplyHandler) Update(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	var req replyDTO.ReplyUpdateDTO

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	post, err := handler.service.Update(&id, &req)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Não foi possivel atualizar o comentário", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, post)
}

func (handler *ReplyHandler) Delete(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	if err := handler.service.Delete(&id); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Err ao remover comentário", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, successDTO.SuccessResponse{Message: "Comentário removido com sucesso"})
}

func (handler *ReplyHandler) FindAll(context *gin.Context) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("page_size", "1"))

	likes, err := handler.service.FindAll(page, pageSize)
	if err != nil {
		context.JSON(http.StatusNoContent, errorDTO.ErrorResponse{Error: "Erro ao buscar a lista de comentários", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, likes)
}

func (handler *ReplyHandler) FindById(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	like, err := handler.service.FindById(&id)

	if err != nil {
		context.JSON(http.StatusNoContent, errorDTO.ErrorResponse{Error: "Erro ocorrido ao buscar o comentário", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, like)
}
