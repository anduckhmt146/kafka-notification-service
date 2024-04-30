package handlers

import (
	"errors"
	"net/http"

	"github.com/anduckhmt146/kakfa-consumer/internal/constants"
	"github.com/anduckhmt146/kakfa-consumer/internal/dtos"
	"github.com/anduckhmt146/kakfa-consumer/internal/services"
	"github.com/gin-gonic/gin"
)

type INotificationHandler interface {
	SendNotification(ctx *gin.Context)
}

type notiHandler struct {
	notiService services.INotificationService
}

func NewNotificationHandler(notiService services.INotificationService) INotificationHandler {
	return &notiHandler{notiService: notiService}
}

func (n *notiHandler) SendNotification(ctx *gin.Context) {
	var messageDtos dtos.MessageDtos
	if err := ctx.BindJSON(&messageDtos); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	err := n.notiService.SendNotification(messageDtos.FromID, messageDtos.ToID, messageDtos.Message)
	if errors.Is(err, constants.ErrUserNotFoundInProducer) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": messageDtos.Message,
	})
}
