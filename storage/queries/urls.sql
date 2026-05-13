-- name: AddURL :one
INSERT INTO
    urls (short_code, original_url)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: GetLongURL :one
SELECT
    original_url
FROM
    urls
WHERE
    short_code = $1;

-- name: UpdateURL :one
UPDATE urls
SET
    original_url = $2,
    updated_at = NOW()
WHERE
    short_code = $1
RETURNING
    *;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE
    short_code = $1;

-- name: IncrementClicks :exec
UPDATE urls
SET
    clicks = clicks + 1
WHERE
    short_code = $1;

-- name: GetURLByLongURL :one
SELECT
    *
FROM
    urls
WHERE
    original_url = $1
LIMIT
    1;
