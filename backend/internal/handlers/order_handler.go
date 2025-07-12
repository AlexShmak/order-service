package handlers

import (
	"encoding/json"
	"github.com/AlexShmak/wb_test_task_l0/internal/kafka"
	"github.com/AlexShmak/wb_test_task_l0/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetOrderByIDHandler(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		h.Logger.Error("Unauthorized access attempt to create order")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Logger.Error("invalid order ID", "error", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.Storage.Orders.GetByID(c.Request.Context(), orderID, userId.(int64))
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

func (h *Handler) CreateOrderHandler(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		h.Logger.Error("Unauthorized access attempt to create order")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var orderRequest struct {
		DeliveryInfo struct {
			Name    string `json:"name" binding:"required"`
			Phone   string `json:"phone" binding:"required"`
			Zip     string `json:"zip" binding:"required"`
			City    string `json:"city" binding:"required"`
			Address string `json:"address" binding:"required"`
			Region  string `json:"region" binding:"required"`
			Email   string `json:"email" binding:"required,email"`
		} `json:"delivery" binding:"required"`
		Payment struct {
			Currency     string `json:"currency" binding:"required"`
			Provider     string `json:"provider" binding:"required"`
			Bank         string `json:"bank" binding:"required"`
			DeliveryCost int    `json:"delivery_cost" binding:"required"`
			GoodsTotal   int    `json:"goods_total" binding:"required"`
			CustomFee    int    `json:"custom_fee" binding:"gte=0"`
		} `json:"payment" binding:"required"`
		Items []struct {
			ChrtID     int    `json:"chrt_id" binding:"required"`
			Name       string `json:"name" binding:"required"`
			Price      int    `json:"price" binding:"required"`
			Sale       int    `json:"sale" binding:"required,min=0,max=100"`
			Size       string `json:"size" binding:"required"`
			Brand      string `json:"brand" binding:"required"`
			TotalPrice int    `json:"total_price" binding:"required"`
			NmID       int    `json:"nm_id" binding:"required"`
		} `json:"items" binding:"required,gt=0,dive"`
		Locale          string `json:"locale" binding:"required"`
		DeliveryService string `json:"delivery_service" binding:"required"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		h.Logger.Error("invalid order request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderUID := uuid.New().String()
	trackNumber := "WB" + uuid.New().String()[:10]
	now := time.Now().UTC()

	delivery := storage.Delivery{
		Name:    orderRequest.DeliveryInfo.Name,
		Phone:   orderRequest.DeliveryInfo.Phone,
		Zip:     orderRequest.DeliveryInfo.Zip,
		City:    orderRequest.DeliveryInfo.City,
		Address: orderRequest.DeliveryInfo.Address,
		Region:  orderRequest.DeliveryInfo.Region,
		Email:   orderRequest.DeliveryInfo.Email,
	}

	payment := storage.Payment{
		Transaction:  orderUID,
		RequestID:    "",
		Currency:     orderRequest.Payment.Currency,
		Provider:     orderRequest.Payment.Provider,
		Amount:       orderRequest.Payment.DeliveryCost + orderRequest.Payment.GoodsTotal,
		PaymentDt:    now.Unix(),
		Bank:         orderRequest.Payment.Bank,
		DeliveryCost: orderRequest.Payment.DeliveryCost,
		GoodsTotal:   orderRequest.Payment.GoodsTotal,
		CustomFee:    orderRequest.Payment.CustomFee,
	}

	items := make([]storage.Item, 0, len(orderRequest.Items))
	for _, item := range orderRequest.Items {
		items = append(items, storage.Item{
			ChrtID:      item.ChrtID,
			TrackNumber: trackNumber,
			Price:       item.Price,
			RID:         uuid.New().String(),
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NMID:        item.NmID,
			Brand:       item.Brand,
			Status:      202,
		})
	}

	order := &storage.Order{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             "WBIL",
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            orderRequest.Locale,
		InternalSignature: "",
		CustomerID:        strconv.FormatInt(userId.(int64), 10),
		DeliveryService:   orderRequest.DeliveryService,
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       now,
		OofShard:          "1",
	}

	//if err := h.Storage.Orders.Create(c.Request.Context(), order); err != nil {
	//	h.Logger.Error("failed to create order", "error", err)
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
	//	return
	//}

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		h.Logger.Error("failed to marshal order", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Send encoded order to kafka
	err = kafka.PushOrderToQueue(h.Config.Kafka.Topic, orderInBytes, h.Config)
	if err != nil {
		h.Logger.Error("failed to push order to kafka", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	h.Logger.Info("order created successfully", "info", orderUID)
	c.JSON(http.StatusCreated, gin.H{
		"message":      "Order placed successfully",
		"order_uid":    order.OrderUID,
		"track_number": order.TrackNumber,
	})
}
