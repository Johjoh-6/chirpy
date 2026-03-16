-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  user_id,
  token,
  expires_at,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, NOW(), NOW()
) RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();
