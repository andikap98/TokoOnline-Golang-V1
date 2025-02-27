package middleware

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := os.Getenv("ADMIN_SECRET")
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Akses tidak diizinkan"})
			c.Abort() //menghentikan middleware
			return
		}

		if auth != key {
			log.Printf("Akses tidak diizinkan %s", auth)
			c.JSON(401, gin.H{"error": "Akses tidak diizinkan"})
			c.Abort() // menghentikan middleware
			return
		}

		c.Next() // Jika semua kondisi terpenuhi, lanjutkan ke middleware berikutnya
	}
}
