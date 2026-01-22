package applications

import "iiitn-career-portal/internal/models"

type BulkStatusUpdateRequest struct {
	ApplicationIDs []uint                   `json:"application_ids" binding:"required,min=1"`
	NewStatus      models.ApplicationStatus `json:"new_status" binding:"required"`
}

type ApplicationListQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`

	Status models.ApplicationStatus `form:"status"`
	JobID  uint                     `form:"job_id"`

	JobDomain models.JobDomain `form:"job_domain"`
	JobType   models.JobType   `form:"job_type"`

	Search string `form:"search"`

	SortBy  string `form:"sort_by"`  // created_at, status
	SortDir string `form:"sort_dir"` // asc, desc
}
