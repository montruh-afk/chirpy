-- name: CreateUser :one

INSERT INTO users (id, created_at, updated_at, email, hashed_password) 

VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

--

-- name: DeleteUser :exec

DELETE FROM users WHERE id = $1;

--

-- name: ResetUsers :exec

DELETE FROM users;

--

-- name: GetuserByEmail :one

SELECT id, created_at, updated_at, email, hashed_password FROM users WHERE email = $1;

--

-- name: GetUserFromRefreshToken :one

SELECT users.id, users.created_at, users.updated_at, users.email, refresh_tokens.token FROM refresh_tokens 
INNER JOIN users ON refresh_tokens.user_id = users.id 
WHERE refresh_tokens.token = $1;

--

-- name: GetUserByID :one 
SELECT id, email FROM users WHERE id = $1;

--