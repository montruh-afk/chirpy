# Chirpy
Chirpy is a backend REST API for a microblogging platform, built using Go and PostgreSQL. It features user authentication with JWTs and refresh tokens, password hashing, a base for content moderation filters as well as that for an admin panel (used for metrics tracking), webhook handling for premium tier upgrades and health checks.

## Getting Started
Chirpy requires both Go (net/http) and PostgreSQL (via SQLc and Goose for migrations) to be installed to work correctly.

### Installing postgres
**Using macOS with brew:**
type `brew install postgresql@15` in the terminal

**Using linux / WSL (Debian):**
`sudo apt update && sudo apt install postgresql postgresql-contrib`

###### Ensure the installation worked:
`psql --version`

### Installing Go
Docs available at [Installing Go](https://go.dev/doc/install).


## Installation & Running

1. Clone the repository:
```bash
git clone https://github.com/montruh-afk/chirpy.git
cd chirpy
```
2. Set up your .env file:
A few variables are required to start up chirpy, copy the example environment file and fill in your values:

```bash
cp .env.example .env
```
3. Database Setup
Before running Chirpy, a user needs a database created and the schema applied with Goose:
**Create a database in PostgreSQL:**
```bash
createdb chirpy
```

**Run migrations (requires Goose):**
```bash
goose -dir sql/schema postgres "postgres://localhost:5432/chirpy?sslmode=disable" up
```

3. Run the server:
```bash 
go run .
```
By default, the server will be available at **http://localhost:8080**

## Key Endpoints

- `POST /api/users` - Create user
- `POST /api/login` - Authenticate user & get JWT tokens
- `POST /api/refresh` - Refresh access token
- `POST /api/chirps` - Create a chirp
- `GET /api/chirps` - List chirps (with optional author/sort queries)
example usage:
`GET http://localhost:8080/api/chirps?sort=asc`
`GET http://localhost:8080/api/chirps?sort=desc`
`GET http://localhost:8080/api/chirps?author=<author_id>`

- `POST /api/polka/webhooks` - Webhook for user upgrades
- `GET /admin/metrics` - View server hit counts

All handlers are registered in the file *initializer.go*