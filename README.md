# chirpy
lesson from Boot.dev build with Go, Postgres, SQLC. We don't use any frameworks or ORM libraries.

## Goals
The Goals of this project are:
- Understand what web servers are and how they power real-world web applications
- Build a production-style HTTP server in Go, without the use of a framework
- Use JSON, headers, and status codes to communicate with clients via a RESTful API
- Learn what makes Go a great language for building fast web servers
- Use type safe SQL to store and retrieve data from a Postgres database
- Implement a secure authentication/authorization system with well-tested cryptography libraries
- Build and understand webhooks and API keys
- Document the REST API with markdown

# Usage 
You need to have a `.env` file with the following variables:
```env
DB_URL="postgres://youname:@localhost:5432/chirpy?sslmode=disable"
PLATFORM=dev // or other
JWT_SECRET=
POLKA_KEY=
```
This are loaded by the `dotenv` package.

## API

Base URL: `http://localhost:8080`

### Conventions

- All JSON responses use `Content-Type: application/json`
- Error responses are JSON shaped like:
  - `{ "error": "message" }`
- Auth headers:
  - JWT endpoints expect: `Authorization: Bearer <jwt>`
  - Polka webhook expects: `Authorization: ApiKey <polka_key>`

---

## Health

### `GET /api/healthz`
Returns plain text `"OK"`.

**Response**
- `200 OK` (`text/plain`): `OK`

---

## Admin

### `GET /admin/metrics`
Returns an HTML page with the file-server hit count.

**Response**
- `200 OK` (`text/html`)

### `POST /admin/reset`
Deletes all users and resets metrics **only if** `PLATFORM=dev`.

**Responses**
- `200 OK`: `Reset` (plain text)
- `403 Forbidden`: `{ "error": "Platform is not dev" }`

---

## Users

### `POST /api/users`
Create a new user.

**Request body**
- `email` (string, required)
- `password` (string, required)

**Responses**
- `201 Created`:
  - `{ "id", "created_at", "updated_at", "email", "is_chirpy_red" }`
- `400 Bad Request`:
  - `{ "error": "Email is required" }`
  - `{ "error": "Password is required" }`
  - or JSON decode error string

### `PUT /api/users`
Update the authenticated user’s email/password.

**Auth**
- `Authorization: Bearer <jwt>`

**Request body**
- `email` (string)
- `password` (string)

**Responses**
- `200 OK`:
  - `{ "id", "created_at", "updated_at", "email", "is_chirpy_red" }`
- `400 Bad Request`: JSON decode error string
- `401 Unauthorized`:
  - bearer/jwt validation error string
  - `{ "error": "User not found" }` (if user id from JWT doesn’t exist)

---

## Auth

### `POST /api/login`
Log in with email/password. Returns a short-lived access JWT and a long-lived refresh token.

**Request body**
- `email` (string, required)
- `password` (string, required)

**Responses**
- `200 OK`:
  - `{ "id", "created_at", "updated_at", "email", "is_chirpy_red", "token", "refresh_token" }`
- `400 Bad Request`:
  - `{ "error": "Email is required" }`
  - `{ "error": "Password is required" }`
- `401 Unauthorized`:
  - `{ "error": "Incorrect email or password" }`

### `POST /api/refresh`
Exchange a valid (not revoked, not expired) refresh token for a new access JWT.

**Auth**
- `Authorization: Bearer <refresh_token>`

**Responses**
- `200 OK`: `{ "token": "<jwt>" }`
- `401 Unauthorized`: error string (missing/invalid bearer, invalid/expired/revoked refresh token)

### `POST /api/revoke`
Revoke a refresh token.

**Auth**
- `Authorization: Bearer <refresh_token>`

**Responses**
- `204 No Content`
- `401 Unauthorized`: error string (missing/invalid bearer)
- `500 Internal Server Error`: database error string

---

## Chirps

### `POST /api/chirps`
Create a chirp for the authenticated user.

**Auth**
- `Authorization: Bearer <jwt>`

**Request body**
- `body` (string, required, 1..140 chars)
- Profanity filter replaces (case-insensitive exact-word match): `kerfuffle`, `sharbert`, `fornax` with `****`.

**Responses**
- `201 Created`:
  - `{ "id", "created_at", "updated_at", "body", "user_id" }`
- `400 Bad Request`:
  - `{ "error": "Chirp is required" }`
  - `{ "error": "Chirp is too long" }`
  - or JSON decode error string
- `401 Unauthorized`: bearer/jwt validation error string

### `GET /api/chirps`
List chirps.

**Query params**
- `author_id` (uuid, optional): filter to a specific user
- `sort` (string, optional): `asc` (default) or `desc`
  - any other value => `400`

**Responses**
- `200 OK`: array of chirps:
  - `[ { "id", "created_at", "updated_at", "body", "user_id" }, ... ]`
- `400 Bad Request`:
  - `{ "error": "invalid sort" }`
  - `{ "error": "invalid author_id" }`

### `GET /api/chirps/{chirpID}`
Get a single chirp by id.

**Path params**
- `chirpID` (uuid)

**Responses**
- `200 OK`: `{ "id", "created_at", "updated_at", "body", "user_id" }`
- `400 Bad Request`: `{ "error": "chirpID is required" }` (uuid parse failed)
- `404 Not Found`: `{ "error": "chirp not found" }`

### `DELETE /api/chirps/{chirpID}`
Delete a chirp. Only the chirp owner can delete.

**Auth**
- `Authorization: Bearer <jwt>`

**Path params**
- `chirpID` (uuid)

**Responses**
- `204 No Content`
- `400 Bad Request`: `{ "error": "chirpID is required" }`
- `401 Unauthorized`: `{ "error": "Unauthorized" }` (missing/invalid token)
- `403 Forbidden`: `{ "error": "Forbidden" }` (not the owner)
- `404 Not Found`: `{ "error": "chirp not found" }`

---

## Webhooks (Polka)

### `POST /api/polka/webhooks`
Handles Polka events. Only `user.upgraded` is acted upon.

**Auth**
- `Authorization: ApiKey <POLKA_KEY>`

**Request body**
- `event` (string)
- `data.user_id` (string uuid)

**Behavior**
- If `event` is not `user.upgraded`: returns `204` and does nothing.
- If `event` is `user.upgraded`: sets `is_chirpy_red=true` for that user.

**Responses**
- `204 No Content` (success or ignored event)
- `401 Unauthorized`:
  - invalid/missing API key
- `400 Bad Request`: invalid `user_id` uuid
- `404 Not Found`: `{ "error": "user not found" }`
