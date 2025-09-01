-- name: InsertCollection :one
INSERT INTO collections (user_id, name, native_language, learning_language)
VALUES ($1, $2, $3, $4) RETURNING id;

-- name: SelectCollections :many
SELECT id, user_id, name, native_language, learning_language FROM collections
WHERE user_id = $1;

-- name: SelectCollection :one
SELECT c.id, c.user_id, c.name, c.native_language, c.learning_language FROM collections c
WHERE c.id = $1;