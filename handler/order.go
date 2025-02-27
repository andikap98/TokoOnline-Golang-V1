package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
	"toko_online/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil data dari pesanan request
		var checkoutOrder model.Checkout
		if err := c.BindJSON(&checkoutOrder); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request data : %v\n", err)
			c.JSON(500, gin.H{"error": "Data produk tidak valid"})
			return
		}

		ids := []string{}
		orderQt := make(map[string]int64)
		for _, o := range checkoutOrder.Products {
			ids = append(ids, o.ID)
			orderQt[o.ID] = o.Quantity
		}

		// TODO: ambil produk data dari database
		products, err := model.SelectProductIn(db, ids)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil product : %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}

		// c.JSON(200, products)
		//TODO: buat kata sandi
		passcode := generatedPasscode(5)
		// TODO: hash kata sandi
		hashCode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)
		if err != nil {
			log.Printf("Terjadi Kesalahan Saat Membuat Hash: %v", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}
		hashCodeString := string(hashCode)

		// TODO: buat order & detail
		order := model.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			Passcode:   &hashCodeString,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}
		for _, p := range products {
			total := p.Price * int64(orderQt[p.ID])

			detail := model.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQt[p.ID],
				Price:     p.Price,
				Total:     total,
			}
			details = append(details, detail)
			order.GrandTotal += total
		}

		if err := model.CreateOrder(db, order, details); err != nil {
			log.Printf("Error saving order: %v", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		orderWithDetail := model.OrderWithDetail{
			Order:   order,
			Details: details,
		}
		orderWithDetail.Order.Passcode = &passcode
		c.JSON(200, orderWithDetail)

	}
}

func generatedPasscode(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[randomGenerator.Intn(len(charset))]
	}

	return string(code)
}

func CheckoutConfirm(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
