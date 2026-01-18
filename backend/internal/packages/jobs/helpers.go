package jobs

import "iiitn-career-portal/internal/models"

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
