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
	AppApplied     ApplicationStatus = "APPLIED"
	AppShortlisted ApplicationStatus = "SHORTLISTED"
	AppInterview   ApplicationStatus = "INTERVIEW"
	AppOffered     ApplicationStatus = "OFFERED"
	AppRejected    ApplicationStatus = "REJECTED"
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
