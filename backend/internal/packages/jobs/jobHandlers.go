package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateJob(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		var req CreateJobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		if strings.TrimSpace(req.Company) == "" {
			c.JSON(400, gin.H{"error": "company is required"})
			return
		}

		if !isValidJobType(req.JobType) {
			c.JSON(400, gin.H{"error": "invalid job_type"})
			return
		}

		if !isValidDomain(req.Domain) {
			c.JSON(400, gin.H{"error": "invalid domain"})
			return
		}

		if len(req.EligibleBatches) == 0 {
			c.JSON(400, gin.H{"error": "eligible_batches required"})
			return
		}

		if req.CTC == nil && req.Stipend == nil {
			c.JSON(400, gin.H{"error": "ctc or stipend required"})
			return
		}

		if req.RegistrationFormURL == nil {
			c.JSON(400, gin.H{"error": "registration_form_url is required"})
			return
		}

		url := strings.TrimSpace(*req.RegistrationFormURL)
		if url == "" {
			c.JSON(400, gin.H{"error": "registration_form_url cannot be empty"})
			return
		}

		if !isValidURL(url) {
			c.JSON(400, gin.H{"error": "registration_form_url must be a valid URL"})
			return
		}

		req.RegistrationFormURL = &url

		batchesJSON, err := json.Marshal(req.EligibleBatches)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to process batches"})
			return
		}

		job := models.Job{
			CollegeID:           *auth.CollegeID,
			Company:             strings.TrimSpace(req.Company),
			Title:               req.Title,
			JobType:             req.JobType,
			Domain:              req.Domain,
			EligibleBatches:     batchesJSON,
			CTC:                 req.CTC,
			Stipend:             req.Stipend,
			Description:         req.Description,
			RegistrationFormURL: req.RegistrationFormURL,
			IsActive:            true,
		}

		if err := db.Create(&job).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create job"})
			return
		}

		c.JSON(201, gin.H{
			"id":      job.ID,
			"message": "job created",
		})
	}
}

func GetJobs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		// -------- Pagination --------
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 50 {
			limit = 10
		}
		offset := (page - 1) * limit

		// -------- Filters --------
		q := strings.TrimSpace(c.Query("q"))
		jobType := c.Query("job_type")
		domain := c.Query("domain")
		sort := c.DefaultQuery("sort", "latest")

		minCTC, _ := strconv.ParseFloat(c.Query("min_ctc"), 64)
		maxCTC, _ := strconv.ParseFloat(c.Query("max_ctc"), 64)
		minStipend, _ := strconv.ParseFloat(c.Query("min_stipend"), 64)
		maxStipend, _ := strconv.ParseFloat(c.Query("max_stipend"), 64)

		batch, _ := strconv.Atoi(c.Query("batch"))

		// -------- Base query --------
		query := db.
			Table("jobs").
			Where("jobs.college_id = ?", auth.CollegeID).
			Where("jobs.is_active = true")

		// -------- Search (temporary DB search) --------
		if q != "" {
			query = query.Where(
				"(jobs.title ILIKE ? OR jobs.company ILIKE ? OR jobs.description ILIKE ?)",
				"%"+q+"%",
				"%"+q+"%",
				"%"+q+"%",
			)
		}

		// -------- Filters --------
		if jobType != "" {
			query = query.Where("jobs.job_type = ?", jobType)
		}
		if domain != "" {
			query = query.Where("jobs.domain = ?", domain)
		}
		if minCTC > 0 {
			query = query.Where("jobs.ctc >= ?", minCTC)
		}
		if maxCTC > 0 {
			query = query.Where("jobs.ctc <= ?", maxCTC)
		}
		if minStipend > 0 {
			query = query.Where("jobs.stipend >= ?", minStipend)
		}
		if maxStipend > 0 {
			query = query.Where("jobs.stipend <= ?", maxStipend)
		}
		if batch > 0 {
			query = query.Where(
				"jobs.eligible_batches @> ?",
				fmt.Sprintf("[%d]", batch),
			)
		}

		// -------- Count --------
		var total int64
		if err := query.Count(&total).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to count jobs"})
			return
		}

		// -------- Sorting --------
		switch sort {
		case "ctc_asc":
			query = query.Order("jobs.ctc ASC NULLS LAST")
		case "ctc_desc":
			query = query.Order("jobs.ctc DESC")
		case "stipend_asc":
			query = query.Order("jobs.stipend ASC NULLS LAST")
		case "stipend_desc":
			query = query.Order("jobs.stipend DESC")
		default:
			query = query.Order("jobs.created_at DESC")
		}

		// -------- Fetch --------
		var jobs []JobListItem
		if err := query.
			Select(`
				jobs.id,
				jobs.company,
				jobs.title,
				jobs.job_type,
				jobs.domain,
				jobs.ctc,
				jobs.stipend,
				jobs.description,
				jobs.created_at
			`).
			Limit(limit).
			Offset(offset).
			Scan(&jobs).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch jobs"})
			return
		}

		// -------- Response --------
		c.JSON(200, gin.H{
			"data": jobs,
			"meta": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	}
}

func GetJobByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		// 1️⃣ Parse job ID
		jobID, err := strconv.Atoi(c.Param("id"))
		if err != nil || jobID <= 0 {
			c.JSON(400, gin.H{"error": "invalid job id"})
			return
		}

		// 2️⃣ Fetch job (college-scoped + active)
		var job models.Job
		if err := db.
			Where("id = ?", jobID).
			Where("college_id = ?", auth.CollegeID).
			Where("is_active = true").
			First(&job).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(404, gin.H{"error": "job not found"})
				return
			}

			c.JSON(500, gin.H{"error": "failed to fetch job"})
			return
		}

		// 3️⃣ Decode eligible batches
		var batches []int
		if err := json.Unmarshal(job.EligibleBatches, &batches); err != nil {
			c.JSON(500, gin.H{"error": "failed to parse eligible batches"})
			return
		}

		// 4️⃣ Respond
		c.JSON(200, JobDetailResponse{
			ID:                  job.ID,
			Company:             job.Company,
			Title:               job.Title,
			JobType:             string(job.JobType),
			Domain:              string(job.Domain),
			EligibleBatches:     batches,
			CTC:                 job.CTC,
			Stipend:             job.Stipend,
			Description:         job.Description,
			RegistrationFormURL: job.RegistrationFormURL,
			CreatedAt:           job.CreatedAt,
		})
	}
}

func ApplyJobIntent(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.MustGet("auth").(*authorization.AuthContext)

		jobID, err := strconv.Atoi(c.Param("id"))
		if err != nil || jobID <= 0 {
			c.JSON(400, gin.H{"error": "invalid job id"})
			return
		}

		var job models.Job
		if err := db.
			Where("id = ?", jobID).
			Where("college_id = ?", auth.CollegeID).
			Where("is_active = true").
			First(&job).Error; err != nil {
			c.JSON(404, gin.H{"error": "job not found"})
			return
		}

		// eligible batch check
		var batches []int
		if err := json.Unmarshal(job.EligibleBatches, &batches); err != nil {
			c.JSON(500, gin.H{"error": "invalid job config"})
			return
		}

		var profile models.StudentProfile
		if err := db.
			Where("user_id = ?", auth.UserID).
			First(&profile).Error; err != nil {
			c.JSON(400, gin.H{"error": "complete profile first"})
			return
		}

		if profile.Batch == 0 {
			c.JSON(400, gin.H{"error": "complete profile before applying"})
			return
		}

		// if !profile.ProfileComplete {
		// 	c.JSON(400, gin.H{"error": "complete profile before applying"})
		// 	return
		// }

		if !containsInt(batches, profile.Batch) {
			c.JSON(400, gin.H{"error": "you are not eligible for this job"})
			return
		}

		// block if already applied
		var count int64
		db.Model(&models.Application{}).
			Where("job_id = ? AND student_id = ?", job.ID, auth.UserID).
			Count(&count)

		if count > 0 {
			c.JSON(409, gin.H{"error": "already applied"})
			return
		}

		// create or refresh intent
		intent := models.ApplicationIntent{
			JobID:     job.ID,
			StudentID: auth.UserID,
			CollegeID: *auth.CollegeID,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		if err := db.
			Where("job_id = ? AND student_id = ?", job.ID, auth.UserID).
			Assign(intent).
			FirstOrCreate(&intent).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create intent"})
			return
		}

		// notification payload
		payload, _ := json.Marshal(gin.H{
			"job_id":  job.ID,
			"title":   job.Title,
			"company": job.Company,
		})

		// DB notification
		_ = db.Create(&models.Notification{
			UserID:   auth.UserID,
			Type:     models.NotificationJobApplyIntent,
			TargetID: intent.ID,
			Payload:  payload,
		}).Error

		redisNotif, _ := json.Marshal(gin.H{
			"type":       models.NotificationJobApplyIntent,
			"target_id":  intent.ID,
			"payload":    json.RawMessage(payload),
			"is_read":    false,
			"created_at": time.Now(),
		})

		// Redis push
		rdb.LPush(
			context.Background(),
			fmt.Sprintf("notifications:user:%d", auth.UserID),
			redisNotif,
		)

		c.JSON(200, gin.H{
			"redirect_url": job.RegistrationFormURL,
			"message":      "complete application and confirm",
		})
	}
}

func UpdateJob(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		jobID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid job id",
			})
			return
		}

		var job models.Job
		if err := db.First(&job, jobID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "job not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		role := models.Role(c.GetString("role"))
		collegeID := c.GetUint("college_id")

		if !canMutateJob(role, job, collegeID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			return
		}

		var req UpdateJobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		updates := map[string]interface{}{}

		if req.Title != nil {
			updates["title"] = *req.Title
		}
		if req.JobType != nil {
			updates["job_type"] = *req.JobType
		}
		if req.Domain != nil {
			updates["domain"] = *req.Domain
		}
		if req.EligibleBatches != nil {
			updates["eligible_batches"] = *req.EligibleBatches
		}
		if req.CTC != nil {
			updates["ctc"] = req.CTC
		}
		if req.Stipend != nil {
			updates["stipend"] = req.Stipend
		}
		if req.RegistrationFormURL != nil {
			updates["registration_form_url"] = req.RegistrationFormURL
		}
		if req.Description != nil {
			updates["description"] = *req.Description
		}
		if req.IsActive != nil {
			updates["is_active"] = *req.IsActive
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no fields to update",
			})
			return
		}

		if err := db.Model(&job).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update job",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "job updated successfully",
		})
	}
}

func DeleteJob(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		jobID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid job id",
			})
			return
		}

		var job models.Job
		if err := db.First(&job, jobID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "job not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		role := models.Role(c.GetString("role"))
		collegeID := c.GetUint("college_id")

		if !canMutateJob(role, job, collegeID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			return
		}

		if err := db.Model(&job).
			Update("is_active", false).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete job",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "job deleted successfully",
		})
	}
}
