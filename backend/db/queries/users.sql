-- name: GetUserByID :one
SELECT id, linkedin_sub, email, role, status, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByLinkedInSub :one
SELECT id, linkedin_sub, email, role, status, created_at, updated_at
FROM users
WHERE linkedin_sub = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, linkedin_sub, email, role, status, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (linkedin_sub, email, role, status)
VALUES ($1, $2, $3, $4)
RETURNING id, linkedin_sub, email, role, status, created_at, updated_at;

-- name: UpdateUserStatus :one
UPDATE users
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING id, linkedin_sub, email, role, status, created_at, updated_at;
