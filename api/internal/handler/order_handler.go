package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/database"
	"github.com/sqlc-dev/pqtype"
)


type OrderHandler struct {
	db *database.Queries
	conn *sql.DB
}


func NewOrderHandler(db *database.Queries, conn *sql.DB) *OrderHandler {
	return &OrderHandler{db: db, conn: conn}
}

type CreateOrderRequest struct {
	UserID   int64   `json:"user_id" binding:"required"`
	ProductID int64   `json:"product_id" binding:"required"`
	Amount int    `json:"amount" binding:"required"`

}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	// Ambil Idempotency-key dari header
	idempotencyKey := ctx.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : "Header Idempotency-Key is required"})
		return
	}

	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	// Mulai database transaction
	// untuk membungkus cek idempotency + logic order
	tx, err := h.conn.BeginTx(context.Background(), nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to begin transaction"})
		return
	}

	// Rollback jika terjadi panic atau error ditengah 
	defer tx.Rollback()

	qTx := h.db.WithTx(tx)

	//  POINT 1 : Cek apakah Idempotency-key sudah pernah dipakai
	existingKey, err := qTx.GetIdempotencyKey(ctx, database.GetIdempotencyKeyParams{
		Key : idempotencyKey,
		UserID : req.UserID, //idealnya ambil dari token jwt middleware
	})

	if err == nil {
		var body map[string]interface{}
		json.Unmarshal(existingKey.ResponseBody.RawMessage, &body)

		tx.Commit()

		ctx.JSON(int(existingKey.ResponseCode.Int32), body)
		return
	} else if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to query idempotency key"})
		return
	}

	// POINT 2 : RACE CONDITION HANDLING (INVENTORY)

	updateProduct, err := qTx.UpdateProductStock(ctx, database.UpdateProductStockParams{
		ID: req.ProductID,
		Amount : int32(req.Amount),
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusConflict, gin.H{"error" : "Insufficient product stock"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to update product stock"})
		return
	}
	newOrder, err := qTx.CreateOrder(ctx, database.CreateOrderParams{
		UserID : req.UserID,
		ProductID : req.ProductID,
		Amount : int32(req.Amount),
	})
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to create order"})
		return
	}

	// Point 3 : Save Idempotency-key beserta response nya
	responseMap := gin.H{
		"message": "Order created successfully",
		"order_id": newOrder.ID,
		"sisa_stock": updateProduct.Stock,
	}
	responseJSON, _ := json.Marshal(responseMap)

	_, err = qTx.CreateIdempotencyKey(ctx, database.CreateIdempotencyKeyParams{
		Key : idempotencyKey,
		UserID : req.UserID,

		// Simpan response code & body
		ResponseCode : sql.NullInt32{
			Int32 : int32(http.StatusOK),
			Valid : true,
		},
		// Simpan response body dalam bentuk JSON
		ResponseBody : pqtype.NullRawMessage{
			RawMessage : responseJSON,
			Valid : true,
		},
	})

	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to save idempotency key"})
		return
	}

	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to commit transaction"})
		return
	}

	ctx.JSON(http.StatusOK, responseMap)

}
		
	
