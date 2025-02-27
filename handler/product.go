package handler

import (
	"database/sql"
	"errors"
	"log"
	"toko_online/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO ambil dari database
		products, err := model.SelectProduct(db)
		if err != nil {
			log.Printf("Gagal mengambil data dari database %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return

		}
		// TODO return data
		c.JSON(200, products)
	}

}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// baca id dari URL
		id := c.Param("id")
		// TODO ambil dari database
		product, err := model.SelectProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Gagal mengambil data dari database %v", err)
				c.JSON(404, gin.H{"error": "Product tidak ditemukan"})
				return
			}
			log.Printf("Gagal mengambil data dari database %v", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}
		// TODO berikan response
		c.JSON(200, product)
	}

}

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		if err := c.Bind(&product); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body %v", err)
			c.JSON(400, gin.H{"error": "Data Produk tidak valid"})
			return
		}

		product.ID = uuid.New().String()
		if err := model.InsertProduct(db, product); err != nil {
			log.Printf("Terjadi kesalahan saat membuat produk %v", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}
		// TODO berikan response
		c.JSON(201, product)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product model.Product
		if err := c.Bind(&product); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body %v", err)
			c.JSON(400, gin.H{"error": "Data Produk tidak valid"})
			return
		}

		productExists, err := model.SelectProductByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil product : %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}
		if productExists.Name != "" {
			productExists.Name = product.Name
		}
		if productExists.Price != 0 {
			productExists.Price = product.Price
		}
		if err := model.UpdateProduct(db, *productExists); err != nil {
			log.Printf("Terjadi kesalahan saat memperbarui produk : %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}
		// TODO berikan response
		c.JSON(201, productExists)
	}
}

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := model.DeleteProduct(db, id); err != nil {
			log.Printf("Terjadi kesalahan saat menghapus produk : %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahn pada server"})
			return
		}
		c.JSON(204, nil)
	}
}
