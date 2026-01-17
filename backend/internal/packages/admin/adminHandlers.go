package admin

import (
	"iiitn-career-portal/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCollege(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name   string `json:"name" binding:"required"`
			Domain string `json:"domain" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		college := models.College{
			Name:   req.Name,
			Domain: req.Domain,
		}

		if err := db.Create(&college).Error; err != nil {
			c.JSON(400, gin.H{"error": "college already exists"})
			return
		}

		c.JSON(201, college)
	}
}
