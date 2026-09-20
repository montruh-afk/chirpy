-- name: CreateUser :one

INSERT INTO users (id, created_at, updated_at, email, hashed_password) 

VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, email, created_at, updated_at, is_chirpy_red;

--

-- name: DeleteUser :exec

DELETE FROM users WHERE id = $1;

--

-- name: ResetUsers :exec

DELETE FROM users;

--

-- name: GetuserByEmail :one

SELECT * FROM users WHERE email = $1;

--

-- name: GetUserFromRefreshToken :one

SELECT users.id, users.email, users.created_at, users.updated_at, users.is_chirpy_red, refresh_tokens.token FROM refresh_tokens 
INNER JOIN users ON refresh_tokens.user_id = users.id 
WHERE refresh_tokens.token = $1;

--

-- name: GetUserByID :one 
SELECT * FROM users WHERE id = $1;

--

-- name: UpdateUser :one 
UPDATE users
SET updated_at = NOW(), email = $1, hashed_password = $2
WHERE id = $3

RETURNING id, email, created_at, updated_at, is_chirpy_red;

--

-- name: SetChirpyRed :exec

UPDATE users
SET updated_at = NOW(), is_chirpy_red = true WHERE id = $1;

--