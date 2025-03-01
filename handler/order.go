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

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: ambil id dari param
		id := c.Param("id")
		//TODO: baca request body
		var confirmeReq model.Confirm
		if err := c.BindJSON(&confirmeReq); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body : %v\n", err)
			c.JSON(500, gin.H{"error": "Data konfirmasi pembayaran tidak valid"})
			return
		}

		//TODO: ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat membaca data pesanan : %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}
		if order.Passcode == nil {
			log.Println("Passcode tidak valid")
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}

		//TODO: cocokan kata sandi pesanan
		if err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(confirmeReq.Passcode)); err != nil {
			log.Printf("Terjadi Kesalahan saat mencocokan kata sandi : %v\n", err)
			c.JSON(401, gin.H{"error": "Tidak diizinkan mengakses pesanan"})
			return
		}

		//TODO: pastikan pesanan belum dibayaran
		if order.PaidAt != nil {
			log.Println("Pesanan sudah dibayar")
			c.JSON(400, gin.H{"error": "Pesanan sudah dibayar"})
			return
		}

		//TODO: cocokan jumlah pembayaran
		if order.GrandTotal != confirmeReq.Amount {
			log.Printf("Jumlah pembayaran tidak sesuai")
			c.JSON(400, gin.H{"error": "Jumlah pembayaran tidak sesuai"})
			return
		}

		//TODO: ubah status pesanan menjadi dibayar
		current := time.Now()
		if err := model.UpdateOrderByID(db, id, confirmeReq, current); err != nil {
			log.Printf("Terjadi kesalahan saat memperbarui data pesanan : %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}

		order.Passcode = nil

		order.PaidAt = &current
		order.PaidBank = &confirmeReq.Bank
		order.PaidAccount = &confirmeReq.AccountNumber

		c.JSON(200, order)

	}
}
func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// TODO: ambil passcode dari query parameter
		passcode := c.Query("passcode")
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat membaca data pesanan : %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}

		if order.Passcode == nil {
			log.Println("Passcode tidak valid")
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode)); err != nil {
			log.Printf("Terjadi Kesalahan saat mencocokan kata sandi : %v\n", err)
			c.JSON(401, gin.H{"error": "Tidak diizinkan mengakses pesanan"})
			return
		}

		order.Passcode = nil
		c.JSON(200, order)

	}
}
