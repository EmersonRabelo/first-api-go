package controller

import (
	"net/http"

	errorDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/error"
	userDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/user"
	service "github.com/EmersonRabelo/first-api-go/internal/service"
	"github.com/gin-gonic/gin"
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
