-- name: InsertUser :one
INSERT INTO users (email, username, role)
VALUES ($1, $2, $3) RETURNING id;