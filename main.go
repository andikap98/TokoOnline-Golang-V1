package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"toko_online/handler"
	"toko_online/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		fmt.Printf("Gagal Membuat koneksi database %v", err)
		os.Exit(1)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Gagal Ping ke database %v", err)
		os.Exit(1)
	}
	fmt.Println("Koneksi database berhasil")

	if _, err = migrate(db); err != nil {
		fmt.Printf("Failed to run migration %v", err)
		os.Exit(1)
	}

	r := gin.Default()
	r.GET("api/v1/products", handler.ListProducts(db))
	r.GET("api/v1/products/:id", handler.GetProduct(db))
	r.POST("api/v1/checkout", handler.CheckoutOrder(db))

	r.POST("api/v1/orders/:id/confirm", handler.ConfirmOrder(db))
	r.GET("api/v1/orders/:id", handler.GetOrder(db))

	r.POST("admin/products", middleware.AdminOnly(), handler.CreateProduct(db))
	r.PUT("admin/products/:id", middleware.AdminOnly(), handler.UpdateProduct(db))
	r.DELETE("admin/products/:id", middleware.AdminOnly(), handler.DeleteProduct(db))

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Gagal menjalankan server %v", err)
		os.Exit(1)
	}
}
