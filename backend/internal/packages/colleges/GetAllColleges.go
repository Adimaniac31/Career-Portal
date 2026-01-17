package colleges

import (
	"iiitn-career-portal/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllColleges(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var colleges []models.College

		if err := db.
			Select("id, name, domain").
			Order("name asc").
			Find(&colleges).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch colleges",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": colleges,
		})
	}
}
