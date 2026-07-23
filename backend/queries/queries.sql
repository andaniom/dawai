-- name: CreateUser :one
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserRoles :many
SELECT role_id FROM user_roles WHERE user_id = $1;

-- name: GetUserSchools :many
SELECT school_id FROM user_schools WHERE user_id = $1;

-- name: CreateSchool :one
INSERT INTO schools (name) VALUES ($1) RETURNING *;

-- name: GetSchool :one
SELECT * FROM schools WHERE id = $1;

-- name: CreateSubject :one
INSERT INTO subjects (school_id, name, description)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetSubjectsBySchool :many
SELECT * FROM subjects WHERE school_id = $1 ORDER BY name;

-- name: GetRubricComponentsBySubject :many
SELECT * FROM rubric_components WHERE subject_id = $1 ORDER BY name;

-- name: CreateRubricComponent :one
INSERT INTO rubric_components (school_id, subject_id, name, description, scale_min, scale_max, weight)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: CreateStudent :one
INSERT INTO students (school_id, user_id, class)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetStudentsBySchool :many
SELECT * FROM students WHERE school_id = $1 ORDER BY created_at;

-- name: GetStudentByID :one
SELECT * FROM students WHERE id = $1;

-- name: CreateAssessment :one
INSERT INTO assessments (school_id, subject_id, student_id, teacher_id, feedback)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetAssessmentsByStudent :many
SELECT * FROM assessments WHERE student_id = $1 ORDER BY submitted_at DESC;

-- name: GetAssessmentByID :one
SELECT * FROM assessments WHERE id = $1;

-- name: CreateAssessmentComponent :one
INSERT INTO assessment_components (assessment_id, rubric_component_id, score)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetAssessmentComponents :many
SELECT * FROM assessment_components WHERE assessment_id = $1;

-- name: BlacklistJWT :one
INSERT INTO jwt_blacklist (jti, expires_at)
VALUES ($1, $2) RETURNING *;

-- name: IsJWTBlacklisted :one
SELECT jti FROM jwt_blacklist WHERE jti = $1;
