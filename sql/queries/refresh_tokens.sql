-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = Now()
WHERE token = $1
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id, expires_at, revoked_at
FROM refresh_tokens
WHERE token = $1;
