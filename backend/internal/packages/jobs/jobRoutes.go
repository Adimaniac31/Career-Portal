package jobs

import (
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rc *redis.Client) {
	jobs := rg.Group("/jobs")
	{
		// Accessible to ALL authenticated users
		jobs.GET("", GetJobs(db))
		jobs.GET("/:id", GetJobByID(db))

		// College admin only
		jobs.POST(
			"",
			authorization.RequireRole(string(models.CollegeAdmin)),
			CreateJob(db),
		)

		// Student only
		jobs.POST(
			"/:id/apply",
			authorization.RequireRole(string(models.Student)),
			ApplyJobIntent(db, rc),
		)
	}
}
