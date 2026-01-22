package jobs

import (
	"iiitn-career-portal/internal/models"
	"net/url"
)

func isValidJobType(j models.JobType) bool {
	switch j {
	case models.JobIntern, models.JobFTE, models.JobInternPPO:
		return true
	default:
		return false
	}
}

func isValidDomain(d models.JobDomain) bool {
	switch d {
	case models.DomainFrontend, models.DomainBackend, models.DomainFullstack,
		models.DomainSDE, models.DomainECE, models.DomainAIML, models.DomainOther:
		return true
	default:
		return false
	}
}

func canMutateJob(role models.Role, job models.Job, collegeID uint) bool {
	if role == models.CollegeAdmin && job.CollegeID == collegeID {
		return true
	}
	return false
}

func isValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}
