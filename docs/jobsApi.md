POST /api/admin/jobs
RequireAuth
RequireRole(COLLEGE_ADMIN)

GET /api/jobs
q=backend
page=1
limit=10
job_type=INTERN
domain=BACKEND
min_ctc=6
max_ctc=20

POST /api/jobs/:job_id/apply
RequireAuth
RequireRole(STUDENT)

PATCH /api/admin/applications/:id/status
RequireAuth
RequireRole(COLLEGE_ADMIN)

GET /api/notifications

both
GET    /api/jobs
GET    /api/jobs/:id

college only
POST   /api/jobs
PATCH  /api/jobs/:id
DELETE /api/jobs/:id

student only
POST   /api/jobs/:id/apply

