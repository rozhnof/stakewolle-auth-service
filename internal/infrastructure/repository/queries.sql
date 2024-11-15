-- queries.sql

-- name: GetUserByID :one
SELECT 
    u.id AS user_id, 
    u.username,
    u.referrer_id,
    u.hash_password,
    s.id AS session_id,
    s.refresh_token,
    s.expired_at AS session_expired_at,
    r.id AS referral_code_id,
    r.code AS referral_code,
    r.expired_at AS referral_code_expired_at
FROM
    users u 
LEFT JOIN session s ON u.id = s.user_id
LEFT JOIN referral_code r ON u.id = r.user_id 
WHERE 
    u.id = $1 
    AND s.deleted_at IS NULL 
    AND r.deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT 
    u.id AS user_id, 
    u.username,
    u.referrer_id,
    u.hash_password,
    s.id AS session_id,
    s.refresh_token,
    s.expired_at AS session_expired_at,
    r.id AS referral_code_id,
    r.code AS referral_code,
    r.expired_at AS referral_code_expired_at
FROM
    users u 
LEFT JOIN session s ON u.id = s.user_id
LEFT JOIN referral_code r ON u.id = r.user_id 
WHERE 
    u.username = $1 
    AND s.deleted_at IS NULL 
    AND r.deleted_at IS NULL;


-- name: GetUserIDByReferralCode :one
SELECT 
    u.id AS user_id
FROM
    users u 
LEFT JOIN session s ON u.id = s.user_id
LEFT JOIN referral_code r ON u.id = r.user_id 
WHERE 
    r.code = $1 
    AND s.deleted_at IS NULL 
    AND r.deleted_at IS NULL;


-- name: CreateUser :one
INSERT INTO users (
    username,
    referrer_id,
    hash_password
) VALUES (
    $1, $2, $3
)
RETURNING 
    id,
    username,
    referrer_id,
    hash_password;


-- name: UpdateUser :one
UPDATE 
    users
SET  
    username = $2,
    hash_password = $3
WHERE 
    id = $1
RETURNING 
    id,
    username,
    referrer_id,
    hash_password;


-- name: DeleteUser :one
UPDATE 
    users
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    id = $1
RETURNING 
    deleted_at;


-- name: CreateSession :one
INSERT INTO session (
    user_id, 
    refresh_token, 
    expired_at
)
VALUES ($1, $2, $3)
RETURNING 
    id, 
    user_id,
    refresh_token, 
    expired_at;


-- name: DeleteSessionByUserID :one
UPDATE 
    session
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    user_id = $1
RETURNING 
    deleted_at;


-- name: CreateReferralCode :one
INSERT INTO referral_code (
    user_id, 
    code,
    expired_at
)
VALUES ($1, $2, $3)
RETURNING 
    id, 
    user_id, 
    code,
    expired_at;


-- name: DeleteReferralCodeByUserID :one
UPDATE 
    referral_code
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    user_id = $1
RETURNING 
    deleted_at;


-- name: List :many
SELECT 
    u.id AS user_id, 
    u.username,
    u.referrer_id,
    u.hash_password,
    s.id AS session_id,
    s.refresh_token,
    s.expired_at AS session_expired_at,
    r.id AS referral_code_id,
    r.code AS referral_code,
    r.expired_at AS referral_code_expired_at
FROM
    users u 
LEFT JOIN session s ON u.id = s.user_id
LEFT JOIN referral_code r ON u.id = r.user_id 
WHERE 
    u.id = ANY(@user_ids::UUID[])
    AND u.referrer_id = ANY(@referrer_ids::UUID[])
    AND s.deleted_at IS NULL 
    AND r.deleted_at IS NULL
LIMIT $1
OFFSET $2;