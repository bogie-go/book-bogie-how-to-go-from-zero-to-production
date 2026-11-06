-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: IsStaff :one
SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND role = 'staff');
