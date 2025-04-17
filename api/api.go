package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	querier repo.Querier
}

func NewMessageHandler(querier repo.Querier) *MessageHandler {
	return &MessageHandler{
		querier: querier,
	}
}

func (h *MessageHandler) WireHttpHandler() http.Handler {

	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	})) //prevents the server from crashing if an error occurs in any route

	r.POST("/message", h.handleCreateMessage)
	r.GET("/message/:id", h.handleGetMessage)
	r.GET("/thread/:id/messages", h.handleGetThreadMessages)
	r.GET("/thread/:id", h.handleGetThread)
	r.DELETE("/message/:id", h.handleDeleteMessage)
	r.PATCH("/message", h.handleUpdateMessage)
	r.POST("/threads", h.handleCreateThread)
	r.POST("/order", h.handleCreateOrder)

	return r
}

func (h *MessageHandler) handleCreateOrder(c *gin.Context) {
	var req repo.CreateOrderParams
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.querier.CreateOrder(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiKey := os.Getenv("CAMPAY_API_KEY")

	ref, err := requestPayment(req.Number, req.Amount, "request payment", apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order":                 order,
		"transaction_reference": ref,
	})
}

func (h *MessageHandler) handleCreateMessage(c *gin.Context) {
	var req repo.CreateMessageParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err = h.querier.GetThread(c, req.ThreadID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no result for this thread, please create a thread before creating a message"})
			return
		}
		c.JSON(http.StatusBadRequest, "No thread found, please create a thread before creating the message")
		return
	}

	message, err := h.querier.CreateMessage(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

type ThreadParam struct {
	Title string `json:"title"`
}

func (h *MessageHandler) handleCreateThread(c *gin.Context) {
	var req ThreadParam
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thread, err := h.querier.CreateThread(c, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, thread)
}

func (h *MessageHandler) handleGetThread(c *gin.Context) {
	id := c.Param("id")

	_, err := h.querier.GetThread(c, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no result for this thread, please create a thread before creating a message"})
			return
		}
	}

	c.JSON(http.StatusOK, "thread found")
}

func (h *MessageHandler) handleGetMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	message, err := h.querier.GetMessageByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) handleGetThreadMessages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	messages, err := h.querier.GetMessagesByThread(c, repo.GetMessagesByThreadParams{
		ThreadID: id,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(messages) == 0 {
		c.JSON(http.StatusNotFound, "No message for this thread")
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) handleDeleteMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.querier.DeleteMessage(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete message"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted successfully"})

}

func (h *MessageHandler) handleUpdateMessage(c *gin.Context) {
	var req repo.UpdateMessageParams
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.querier.UpdateMessage(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Message updated successfully"})
}
