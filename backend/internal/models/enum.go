package models

type Role string
type JobType string
type JobDomain string
type ApplicationStatus string
type NotificationType string

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
	Applied     ApplicationStatus = "APPLIED"
	Shortlisted ApplicationStatus = "SHORTLISTED"
	Interview   ApplicationStatus = "INTERVIEW"
	Offered     ApplicationStatus = "OFFERED"
	Rejected    ApplicationStatus = "REJECTED"
)

const (
	// Jobs
	NotificationNewJob NotificationType = "NEW_JOB"

	// Applications
	NotificationJobApplyIntent    NotificationType = "JOB_APPLY_INTENT"
	NotificationJobApplied        NotificationType = "JOB_APPLIED"
	NotificationApplicationStatus NotificationType = "APPLICATION_STATUS_UPDATE"

	// Discussions
	NotificationDiscussionUpdate NotificationType = "DISCUSSION_UPDATE"
)
