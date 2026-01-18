package applications

import (
	"context"
	"encoding/json"
	"fmt"
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ConfirmApplication(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		// 1️⃣ Parse intent ID
		intentID, err := strconv.Atoi(c.Param("intent_id"))
		if err != nil || intentID <= 0 {
			c.JSON(400, gin.H{"error": "invalid intent id"})
			return
		}

		// 2️⃣ Fetch intent
		var intent models.ApplicationIntent
		if err := db.
			Where("id = ?", intentID).
			Where("student_id = ?", auth.UserID).
			First(&intent).Error; err != nil {

			c.JSON(404, gin.H{"error": "application intent not found"})
			return
		}

		// 3️⃣ Check expiry
		if time.Now().After(intent.ExpiresAt) {
			_ = db.Delete(&intent)
			c.JSON(400, gin.H{"error": "application intent expired"})
			return
		}

		// 4️⃣ Fetch job (safety)
		var job models.Job
		if err := db.
			Where("id = ?", intent.JobID).
			Where("college_id = ?", intent.CollegeID).
			First(&job).Error; err != nil {

			c.JSON(404, gin.H{"error": "job not found"})
			return
		}

		tx := db.Begin()

		// 5️⃣ Create application
		app := models.Application{
			JobID:     intent.JobID,
			StudentID: intent.StudentID,
			CollegeID: intent.CollegeID,
			Status:    models.AppApplied,
		}

		if err := tx.Create(&app).Error; err != nil {
			tx.Rollback()
			c.JSON(409, gin.H{"error": "application already exists"})
			return
		}

		// 6️⃣ Delete intent
		if err := tx.Delete(&intent).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "failed to finalize application"})
			return
		}

		// 7️⃣ Create notification (DB)
		payload, _ := json.Marshal(gin.H{
			"job_id":  job.ID,
			"title":   job.Title,
			"company": job.Company,
		})

		if err := tx.Create(&models.Notification{
			UserID:  auth.UserID,
			Type:    models.NotificationJobApplied,
			Payload: payload,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "failed to notify"})
			return
		}

		if err := tx.Commit().Error; err != nil {
			c.JSON(500, gin.H{"error": "transaction failed"})
			return
		}

		// 8️⃣ Push Redis notification (non-blocking)
		redisPayload, _ := json.Marshal(gin.H{
			"type":    models.NotificationJobApplied,
			"job_id":  job.ID,
			"title":   job.Title,
			"company": job.Company,
		})

		_ = rdb.LPush(
			context.Background(),
			fmt.Sprintf("notifications:user:%d", auth.UserID),
			redisPayload,
		).Err()

		// 9️⃣ Respond
		c.JSON(200, gin.H{
			"message": "application confirmed",
		})
	}
}
