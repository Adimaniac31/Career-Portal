package applications

import (
	"context"
	"encoding/json"
	"iiitn-career-portal/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var validTransitions = map[models.ApplicationStatus][]models.ApplicationStatus{
	models.Applied: {
		models.Shortlisted,
		models.Rejected,
	},
	models.Shortlisted: {
		models.Interview,
		models.Rejected,
	},
	models.Interview: {
		models.Offered,
		models.Rejected,
	},
}

func isValidTransition(from, to models.ApplicationStatus) bool {
	nextStates, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range nextStates {
		if s == to {
			return true
		}
	}
	return false
}

func allowedSortColumn(col string) string {
	switch col {
	case "created_at", "status":
		return col
	default:
		return "created_at"
	}
}

func applyDefaults(q *ApplicationListQuery) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortDir != "asc" {
		q.SortDir = "desc"
	}
}

func buildApplicationQuery(
	db *gorm.DB,
	role models.Role,
	userID uint,
	collegeID uint,
	q ApplicationListQuery,
) *gorm.DB {

	query := db.Model(&models.Application{}).
		Preload("Job").
		Preload("Student")

	// Role scoping
	switch role {
	case models.Student:
		query = query.Where("student_id = ?", userID)

	case models.CollegeAdmin:
		query = query.Where("college_id = ?", collegeID)
	}

	// Application-level filters
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}

	if q.JobID != 0 {
		query = query.Where("job_id = ?", q.JobID)
	}

	// Job-level filters (require join)
	if q.JobDomain != "" || q.JobType != "" {
		query = query.Joins("JOIN jobs ON jobs.id = applications.job_id")
	}

	if q.JobDomain != "" {
		query = query.Where("jobs.domain = ?", q.JobDomain)
	}

	if q.JobType != "" {
		query = query.Where("jobs.job_type = ?", q.JobType)
	}

	// Search
	if q.Search != "" {
		search := "%" + q.Search + "%"
		query = query.
			Joins("JOIN users ON users.id = applications.student_id").
			Joins("JOIN jobs ON jobs.id = applications.job_id").
			Joins("JOIN companies ON companies.id = jobs.company_id").
			Where(`
				users.name ILIKE ? OR
				jobs.title ILIKE ? OR
				companies.name ILIKE ?
			`, search, search, search)
	}

	// Sorting
	sortCol := allowedSortColumn(q.SortBy)
	query = query.Order(sortCol + " " + q.SortDir)

	return query
}

func enqueueApplicationStatusNotifications(
	rdb *redis.Client,
	applications []models.Application,
	newStatus models.ApplicationStatus,
) error {

	ctx := context.Background()

	for _, app := range applications {
		payload, _ := json.Marshal(gin.H{
			"application_id": app.ID,
			"new_status":     newStatus,
		})

		msg, _ := json.Marshal(gin.H{
			"user_id":    app.StudentID,
			"type":       models.NotificationApplicationStatus,
			"target_id":  app.ID,
			"payload":    json.RawMessage(payload),
			"created_at": time.Now(),
		})

		if err := rdb.LPush(ctx, "notifications:queue", msg).Err(); err != nil {
			return err
		}
	}

	return nil
}
