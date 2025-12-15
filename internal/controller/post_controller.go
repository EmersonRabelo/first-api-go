package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	errorDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/error"
	postDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/post"
	successDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/success"
	service "github.com/EmersonRabelo/first-api-go/internal/service"
)

type PostHandler struct {
	service service.PostService
}

func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (handler *PostHandler) Create(context *gin.Context) {
	var req postDTO.PostCreateDTO

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	post, err := handler.service.Create(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Não foi possivel criar o usuário", Details: err.Error()})
		return
	}

	context.JSON(http.StatusCreated, post)
}

func (handler *PostHandler) Update(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	var req postDTO.PostUpdateDTO

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	post, err := handler.service.Update(&id, &req)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Não foi possivel atualizar o post", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, post)
}

func (handler *PostHandler) Delete(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	if err := handler.service.Delete(&id); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Err ao deletar post", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, successDTO.SuccessResponse{Message: "Post deletado com sucesso"})
}

func (handler *PostHandler) FindById(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	post, err := handler.service.FindById(&id)

	if err != nil {
		context.JSON(http.StatusNoContent, errorDTO.ErrorResponse{Error: "Erro ocorrido ao buscar o usuário", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, post)
}

func (handler *PostHandler) FindAll(context *gin.Context) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("page_size", "10"))

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

	var userId *uuid.UUID = nil
	userIdParam := context.Query("user_id")
	if userIdParam != "" {
		id, err := uuid.Parse(userIdParam)
		if err != nil {
			context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{
				Error:   "Id do usuário com formato inválido",
				Details: err.Error(),
			})
			return
		}
		userId = &id
	}

	posts, err := handler.service.FindAll(userId, start, end, page, pageSize)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{
			Error:   "Erro ao buscar a lista de postagens",
			Details: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, posts)
}
