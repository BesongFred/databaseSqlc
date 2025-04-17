-- name: CreateMessage :one
INSERT INTO message (thread_id,content)
VALUES ($1, $2)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetMessagesByThread :many
SELECT * FROM message
WHERE thread_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteMessage :exec
DELETE FROM message WHERE id = $1;

-- name: UpdateMessage :exec
UPDATE message 
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAll :exec
DELETE FROM message;


-- name: CreateThread :one
INSERT INTO thread (title) 
VALUES ($1)
RETURNING *; 

-- name: GetThreadByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetThread :one
SELECT * FROM thread
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders(amount,number)
VALUES($1,$2)
RETURNING *;