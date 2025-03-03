-- name: CreateUserWithEmail :one
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
RETURNING id, email, name, created_at, updated_at;


-- name: CreateUserWithGoogle :one
INSERT INTO users (google_id, email, name)
VALUES ($1, $2, $3)
RETURNING id, email, google_id, name;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, google_id, name, created_at, updated_at
FROM users
WHERE email = $1;


-- name: GetUserByGoogleID :one
SELECT id, email, google_id, name, created_at, updated_at
FROM users
WHERE google_id = $1;


-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE email = $1;