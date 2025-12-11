package controller

import (
	"net/http"
	"strconv"

	errorDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/error"
	successDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/success"
	userDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/user"
	service "github.com/EmersonRabelo/first-api-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) Create(context *gin.Context) {
	var req userDTO.CreateDTO

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	user, err := handler.service.Create(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Não foi possivel criar o usuário.", Details: err.Error()})
		return
	}

	context.JSON(http.StatusCreated, user)
}

func (handler *UserHandler) Update(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	var req userDTO.UpdateDTO

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	user, err := handler.service.Update(&id, &req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

func (handler *UserHandler) Delete(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	if err := handler.service.Delete(&id); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Erro ao deletar usuário", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, successDTO.SuccessResponse{Message: "Usuário deletado com sucesso"})

}

func (handler *UserHandler) FindAll(context *gin.Context) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("page_size", "10"))

	users, err := handler.service.FindAll(page, pageSize)
	if err != nil {
		context.JSON(http.StatusNoContent, errorDTO.ErrorResponse{Error: "Erro ao buscar a lista de usuários", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, users)
}

func (handler *UserHandler) FindById(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	user, err := handler.service.FindById(&id)

	if err != nil {
		context.JSON(http.StatusNoContent, errorDTO.ErrorResponse{Error: "Erro ocorrido ao buscar o usuário", Details: err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}
