package jobs

import (
	"iiitn-career-portal/internal/models"
	"time"

	"gorm.io/datatypes"
)

type CreateJobRequest struct {
	Company             string           `json:"company" binding:"required"`
	Title               string           `json:"title" binding:"required"`
	JobType             models.JobType   `json:"job_type" binding:"required"`
	Domain              models.JobDomain `json:"domain" binding:"required"`
	EligibleBatches     []int            `json:"eligible_batches" binding:"required"`
	CTC                 *float64         `json:"ctc"`
	Stipend             *float64         `json:"stipend"`
	Description         string           `json:"description"`
	RegistrationFormURL *string          `json:"registration_form_url" binding:"required"`
}

type JobListItem struct {
	ID          uint      `json:"id"`
	Company     string    `json:"company"`
	Title       string    `json:"title"`
	JobType     string    `json:"job_type"`
	Domain      string    `json:"domain"`
	CTC         *float64  `json:"ctc"`
	Stipend     *float64  `json:"stipend"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type JobDetailResponse struct {
	ID                  uint      `json:"id"`
	Company             string    `json:"company"`
	Title               string    `json:"title"`
	JobType             string    `json:"job_type"`
	Domain              string    `json:"domain"`
	EligibleBatches     []int     `json:"eligible_batches"`
	CTC                 *float64  `json:"ctc"`
	Stipend             *float64  `json:"stipend"`
	Description         string    `json:"description"`
	RegistrationFormURL *string   `json:"registration_form_url"`
	CreatedAt           time.Time `json:"created_at"`
}

type UpdateJobRequest struct {
	Title   *string           `json:"title"`
	JobType *models.JobType   `json:"job_type"`
	Domain  *models.JobDomain `json:"domain"`

	EligibleBatches *datatypes.JSON `json:"eligible_batches"`

	CTC     *float64 `json:"ctc"`
	Stipend *float64 `json:"stipend"`

	RegistrationFormURL *string `json:"registration_form_url"`
	Description         *string `json:"description"`

	IsActive *bool `json:"is_active"`
}
