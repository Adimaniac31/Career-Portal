package applications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"
	"log"
	"net/http"
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
			Status:    models.Applied,
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
			UserID:   auth.UserID,
			Type:     models.NotificationJobApplied,
			TargetID: app.ID,
			Payload:  payload,
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

func BulkUpdateApplicationStatus(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req BulkStatusUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// extracted from auth middleware
		collegeID := c.GetUint("college_id")

		var applications []models.Application

		err := db.
			Joins("JOIN jobs ON jobs.id = applications.job_id").
			Where("applications.id IN ?", req.ApplicationIDs).
			Where("jobs.college_id = ?", collegeID).
			Find(&applications).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		// Security check: prevent cross-college updates
		if len(applications) != len(req.ApplicationIDs) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "one or more applications do not belong to your college",
			})
			return
		}

		// Status transition validation
		for _, app := range applications {
			if !isValidTransition(app.Status, req.NewStatus) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid status transition",
				})
				return
			}
		}

		// Transaction
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to start transaction",
			})
			return
		}

		if err := tx.
			Model(&models.Application{}).
			Where("id IN ?", req.ApplicationIDs).
			Update("status", req.NewStatus).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update application statuses",
			})
			return
		}

		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "transaction commit failed",
			})
			return
		}

		if err := enqueueApplicationStatusNotifications(
			redisClient,
			applications,
			req.NewStatus,
		); err != nil {
			// IMPORTANT: do NOT fail the request
			// log this instead
			log.Println("failed to enqueue notifications:", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"updated_count": len(req.ApplicationIDs),
			"new_status":    req.NewStatus,
		})
	}
}

func ListApplications(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var q ApplicationListQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		applyDefaults(&q)

		role := models.Role(c.GetString("role"))
		userID := c.GetUint("user_id")
		collegeID := c.GetUint("college_id")

		baseQuery := buildApplicationQuery(db, role, userID, collegeID, q)

		var total int64
		if err := baseQuery.Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to count applications",
			})
			return
		}

		var applications []models.Application
		if err := baseQuery.
			Limit(q.Limit).
			Offset((q.Page - 1) * q.Limit).
			Find(&applications).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch applications",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": applications,
			"meta": gin.H{
				"page":  q.Page,
				"limit": q.Limit,
				"total": total,
			},
		})
	}
}

func GetApplicationByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idParam := c.Param("id")
		appID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid application id",
			})
			return
		}

		var application models.Application

		err = db.
			Preload("Job").
			Preload("Student").
			First(&application, appID).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "application not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		role := models.Role(c.GetString("role"))
		userID := c.GetUint("user_id")
		collegeID := c.GetUint("college_id")

		// Authorization check
		switch role {

		case models.Student:
			if application.StudentID != userID {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "access denied",
				})
				return
			}

		case models.CollegeAdmin:
			if application.CollegeID != collegeID {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "access denied",
				})
				return
			}

		case models.Admin:
			// full access

		default:
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid role",
			})
			return
		}

		c.JSON(http.StatusOK, application)
	}
}
