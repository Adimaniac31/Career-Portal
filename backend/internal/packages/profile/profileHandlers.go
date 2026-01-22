package profile

import (
	"context"
	"fmt"
	"iiitn-career-portal/internal/config"
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type UpdateProfileRequest struct {
	Name       *string  `json:"name"`
	Batch      *int     `json:"batch"`
	CGPA       *float32 `json:"cgpa"`
	LinkedinID *string  `json:"linkedin_id"`
}

func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		// 1️⃣ Fetch user
		var user models.User
		if err := db.
			Select("id, name, email, role, college_id").
			Where("id = ?", auth.UserID).
			First(&user).Error; err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		// 2️⃣ Fetch student profile (optional)
		var profile models.StudentProfile
		err := db.Where("user_id = ?", auth.UserID).First(&profile).Error

		var profileResp gin.H
		if err == gorm.ErrRecordNotFound {
			profileResp = gin.H{
				"batch":            nil,
				"cgpa":             nil,
				"resume_url":       nil,
				"linkedin_id":      nil,
				"profile_complete": false,
			}
		} else if err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch profile"})
			return
		} else {
			profileResp = gin.H{
				"batch":            profile.Batch,
				"cgpa":             profile.CGPA,
				"resume_url":       profile.ResumeURL,
				"linkedin_id":      profile.LinkedinID,
				"profile_complete": profile.ProfileComplete,
			}
		}

		// 3️⃣ Respond
		c.JSON(200, gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"college_id": user.CollegeID,
			"profile":    profileResp,
		})
	}
}

func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		// 🔒 Only students can update student profile
		if auth.Role != string(models.Student) {
			c.JSON(403, gin.H{"error": "profile not applicable for this role"})
			return
		}

		var req UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		// Validate inputs
		if req.CGPA != nil && (*req.CGPA < 0 || *req.CGPA > 10) {
			c.JSON(400, gin.H{"error": "invalid cgpa"})
			return
		}
		if req.Batch != nil && *req.Batch < 2000 {
			c.JSON(400, gin.H{"error": "invalid batch"})
			return
		}

		tx := db.Begin()

		// 1️⃣ Update user name
		if req.Name != nil {
			if err := tx.
				Model(&models.User{}).
				Where("id = ?", auth.UserID).
				Update("name", *req.Name).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "failed to update name"})
				return
			}
		}

		// 2️⃣ Upsert student profile
		var profile models.StudentProfile
		err := tx.Where("user_id = ?", auth.UserID).
			First(&profile).Error

		if err == gorm.ErrRecordNotFound {
			profile = models.StudentProfile{
				UserID: auth.UserID,
			}
		} else if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "failed to load profile"})
			return
		}

		if req.Batch != nil {
			profile.Batch = *req.Batch
		}
		if req.CGPA != nil {
			profile.CGPA = req.CGPA
		}
		if req.LinkedinID != nil {
			profile.LinkedinID = *req.LinkedinID
		}

		// 3️⃣ Compute profile completeness
		profile.ProfileComplete =
			profile.Batch != 0 &&
				profile.ResumeURL != nil

		if err := tx.Save(&profile).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "failed to update profile"})
			return
		}

		if err := tx.Commit().Error; err != nil {
			c.JSON(500, gin.H{"error": "transaction failed"})
			return
		}

		c.JSON(200, gin.H{"message": "profile updated"})
	}
}

func uploadResume(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		file, err := c.FormFile("resume")
		if err != nil {
			c.JSON(400, gin.H{"error": "resume file required"})
			return
		}

		// size limit: 2MB
		if file.Size > 2*1024*1024 {
			c.JSON(400, gin.H{"error": "resume too large"})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to open file"})
			return
		}
		defer f.Close()

		// detect MIME
		header := make([]byte, 512)
		_, _ = f.Read(header)
		contentType := http.DetectContentType(header)

		if contentType != "application/pdf" {
			c.JSON(400, gin.H{"error": "only PDF resumes allowed"})
			return
		}

		_, _ = f.Seek(0, 0)

		// virus scan hook
		if err := scanForVirus(f); err != nil {
			c.JSON(400, gin.H{"error": "virus detected"})
			return
		}

		_, _ = f.Seek(0, 0)

		// MinIO client
		minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
			Secure: cfg.MinioUseSSL,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "storage unavailable"})
			return
		}

		objectPath := fmt.Sprintf(
			"resumes/%d/%d/resume.pdf",
			*auth.CollegeID,
			auth.UserID,
		)

		// upload (overwrite-safe)
		_, err = minioClient.PutObject(
			context.Background(),
			cfg.MinioBucket,
			objectPath,
			f,
			file.Size,
			minio.PutObjectOptions{
				ContentType: "application/pdf",
			},
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "upload failed"})
			return
		}

		resumeURL := fmt.Sprintf(
			"%s/%s/%s",
			cfg.MinioPublicURL,
			cfg.MinioBucket,
			objectPath,
		)

		// DB transaction
		tx := db.Begin()

		if err := tx.Model(&models.StudentProfile{}).
			Where("user_id = ?", auth.UserID).
			Updates(map[string]interface{}{
				"resume_url":       resumeURL,
				"profile_complete": true,
			}).Error; err != nil {

			tx.Rollback()

			// cleanup MinIO on DB failure
			_ = minioClient.RemoveObject(
				context.Background(),
				cfg.MinioBucket,
				objectPath,
				minio.RemoveObjectOptions{},
			)

			c.JSON(500, gin.H{"error": "failed to save resume"})
			return
		}

		tx.Commit()

		c.JSON(200, gin.H{
			"message":    "resume uploaded successfully",
			"resume_url": resumeURL,
		})
	}
}
