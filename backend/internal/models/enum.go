package models

type Role string
type JobType string
type JobDomain string
type ApplicationStatus string

const (
	Admin        Role = "admin"
	CollegeAdmin Role = "college_admin"
	Student      Role = "student"
)

const (
	JobIntern    JobType = "INTERN"
	JobFTE       JobType = "FTE"
	JobInternPPO JobType = "INTERN_PPO"
)

const (
	DomainFrontend  JobDomain = "FRONTEND"
	DomainBackend   JobDomain = "BACKEND"
	DomainFullstack JobDomain = "FULLSTACK"
	DomainSDE       JobDomain = "SDE"
	DomainECE       JobDomain = "ECE"
	DomainAIML      JobDomain = "AIML"
	DomainOther     JobDomain = "OTHER"
)

const (
	AppApplied     ApplicationStatus = "APPLIED"
	AppShortlisted ApplicationStatus = "SHORTLISTED"
	AppInterview   ApplicationStatus = "INTERVIEW"
	AppOffered     ApplicationStatus = "OFFERED"
	AppRejected    ApplicationStatus = "REJECTED"
)
