package applications

import (
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	applications := rg.Group("/applications")
	applications.POST(
		"/:intent_id/confirm",
		authorization.RequireRole(string(models.Student)),
		ConfirmApplication(db, redisClient),
	)
	applications.PATCH(
		"/status/bulk",
		authorization.RequireRole(string(models.CollegeAdmin)),
		BulkUpdateApplicationStatus(db, redisClient),
	)
	applications.GET(
		"",
		authorization.RequireRole(
			string(models.Student),
			string(models.CollegeAdmin),
		),
		ListApplications(db),
	)
	applications.GET(
		"/:id",
		authorization.RequireRole(
			string(models.Student),
			string(models.CollegeAdmin),
		),
		GetApplicationByID(db),
	)

}
