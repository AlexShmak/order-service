package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetOrderByIDHandler(c *gin.Context) {

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Logger.Error("invalid order ID", "error", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.Storage.Orders.GetByID(c.Request.Context(), orderID)
	if err != nil {
		h.Logger.Error("failed to get order", "error", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to get order"})
		return
	}
	if order == nil {
		h.Logger.Error("order not found", "id", orderID)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, order)
}
