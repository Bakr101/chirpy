-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_passwords)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: ChangePassEmail :one
UPDATE users
SET email = $1, hashed_passwords = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: SetChirpyRed :one
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
RETURNING *;