API CONTRACTS — IIITN Career Portal
Global conventions (lock these)
Base
/api/v1

Auth

Authorization: Bearer <access_token>

Token issued via OIDC

Backend never trusts frontend role claims

Response shape
{
  "success": true,
  "data": {},
  "error": null
}


Error:

{
  "success": false,
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "You are not allowed to perform this action"
  }
}

Pagination (every list)
{
  "items": [],
  "page": 1,
  "page_size": 20,
  "total": 312
}

1. Auth & Profile
Get current user
GET /auth/me


Response

{
  "id": "uuid",
  "name": "Aditya",
  "email": "bt22cse044@iiitn.ac.in",
  "role": "STUDENT",
  "college_id": "uuid",
  "linkedin_url": "...",
  "settings": {
    "email_notifications": true
  }
}


Backend concepts

OIDC token validation

User bootstrap

Config-based email domain check

2. Colleges (ADMIN only)
Create college
POST /colleges

{
  "name": "IIIT Nagpur",
  "allowed_domains": ["iiitn.ac.in"]
}

List colleges
GET /colleges


RBAC

ADMIN only

3. Jobs
Create job (College_Admin)
POST /jobs

{
  "company": "Google",
  "role": "Intern",
  "job_type": "SDE",
  "batch_eligible": [2025, 2026],
  "stipend": 80000,
  "ctc": null,
  "on_campus": true,
  "registration_url": "https://...",
  "rounds": "OA → Interviews",
  "description": "..."
}

List jobs (Students)
GET /jobs?search=goog&page=1&stipend_min=20000


Lifecycle

Client → Cache → Typesense → DB → Cache

4. Applications (core workflow)
Apply for job
POST /jobs/{job_id}/apply


Response

{
  "application_id": "uuid",
  "status": "APPLIED"
}


Backend guarantees

Idempotent (no duplicate applies)

DB transaction

Async notification

My applications
GET /applications/me

Bulk status update (College_Admin)
PATCH /applications/bulk-status

{
  "job_id": "uuid",
  "from_status": "APPLIED",
  "to_status": "SHORTLISTED"
}


Rules

Valid state transitions enforced

Audit log created

Notifications async

5. Experiences & Discussions
Create experience
POST /experiences

{
  "company": "Amazon",
  "content": "OA experience..."
}

List experiences
GET /experiences?company=Amazon&page=1

Add comment
POST /experiences/{id}/comments

{
  "parent_id": null,
  "content": "Thanks for sharing"
}


Design choice

Flat list + parent_id

Pagination mandatory

6. Company History & Analytics
Company overview
GET /companies/{name}/overview


Response

{
  "company": "Amazon",
  "years": [
    {
      "year": 2024,
      "hired": 12,
      "ppo": 5,
      "avg_ctc": 18
    }
  ]
}


Backend

Aggregation queries

Cached heavily

Cron-precomputed stats

7. Alumni Search
Search alumni
GET /alumni?company=Microsoft


Response

[
  {
    "name": "John Doe",
    "batch": 2022,
    "linkedin_url": "...",
    "email": "visible_if_allowed"
  }
]


Rules

Internal DB only

Privacy filters enforced

Rate-limited

8. Notifications
List notifications
GET /notifications

Mark read
POST /notifications/{id}/read

9. Files (Resumes)
Upload resume
POST /files/resume


Multipart upload

Stored in MinIO

Metadata in DB

10. Observability hooks (important)

Every API must emit:

Trace ID (OpenTelemetry)

Structured logs

Metrics

Example:

POST /jobs/{id}/apply
→ applications_created_total++
→ trace: apply-flow