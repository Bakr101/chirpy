-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at Desc;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirp :one
DELETE FROM chirps
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetChirpsByID :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpsByIDDesc :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC;