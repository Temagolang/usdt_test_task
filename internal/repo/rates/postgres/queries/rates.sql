-- name: InsertRate :exec
INSERT INTO rates (ask, bid, fetched_at)
VALUES ($1, $2, $3);