-- name: CreateProfile :one
INSERT INTO profiles (user_id, first_name, last_name, age, gender, city, interests)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProfileByID :one
SELECT * FROM profiles
WHERE id = $1;

-- name: GetProfileByUserID :one
SELECT * FROM profiles
WHERE user_id = $1;

-- name: UpdateProfile :one
UPDATE profiles 
SET first_name = $2, last_name = $3, age = $4, gender = $5, city = $6, interests = $7
WHERE user_id = $1
RETURNING *;

-- name: SearchProfiles :many
SELECT * FROM profiles
WHERE 
    ($1::text IS NULL OR gender = $1) AND
    ($2::text IS NULL OR city = $2) AND
    ($3::text[] IS NULL OR interests && $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: GetProfilesCount :one
SELECT COUNT(*) FROM profiles
WHERE 
    ($1::text IS NULL OR gender = $1) AND
    ($2::text IS NULL OR city = $2) AND
    ($3::text[] IS NULL OR interests && $3); 