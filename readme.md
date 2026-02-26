# API Quest — Desent Solutions Backend Challenge
## Overview
Build and deploy a REST API that passes 8 levels of increasingly difficult backend challenges. The test runner at `https://www.desent.io/coding-test-backend` will hit your deployed API endpoints in real time.
**Tech Stack:** Golang (Echo framework) + MongoDB  
**Time Budget:** 2–4 hours
---
## The 8 Levels
### Level 1 — Ping
| Method | Endpoint   | Description                        |
|--------|------------|------------------------------------|
| GET    | `/ping`    | Return a simple health check response |
**Expected:** Return `200 OK` with a JSON body like `{"message": "pong"}` (or similar — the test likely checks for a 200 status).
---
### Level 2 — Echo
| Method | Endpoint   | Description                         |
|--------|------------|-------------------------------------|
| POST   | `/echo`    | Echo back the request body as JSON  |
**Expected:** Accept a JSON body and return it as-is with `200 OK`.
---
### Level 3 — CRUD: Create & Read
| Method | Endpoint        | Description               |
|--------|-----------------|---------------------------|
| POST   | `/books`        | Create a new book         |
| GET    | `/books`        | List all books            |
| GET    | `/books/:id`    | Get a single book by ID   |
**MongoDB Collection:** `books`
**Book Model (suggested):**
```go
type Book struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title     string             `json:"title" bson:"title"`
    Author    string             `json:"author" bson:"author"`
    ISBN      string             `json:"isbn" bson:"isbn"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
```
**Notes:**
- POST `/books` should return `201 Created` with the created book (including its generated ID).
- GET `/books` returns an array of all books.
- GET `/books/:id` returns a single book or `404` if not found.
---
### Level 4 — CRUD: Update & Delete
| Method | Endpoint         | Description           |
|--------|------------------|-----------------------|
| PUT    | `/books/:id`     | Update a book by ID   |
| DELETE | `/books/:id`     | Delete a book by ID   |
**Notes:**
- PUT should return `200 OK` with the updated book.
- DELETE should return `200 OK` or `204 No Content`.
- Both should return `404` if the book ID doesn't exist.
---
### Level 5 — Auth Guard
| Method | Endpoint          | Description                                  |
|--------|-------------------|----------------------------------------------|
| POST   | `/auth/token`     | Generate an auth token                       |
| GET    | `/books`          | Protected — requires valid Bearer token      |
**Plan:**
- POST `/auth/token` accepts credentials (e.g., `{"username": "...", "password": "..."}`) and returns a token (JWT or simple token).
- GET `/books` now requires `Authorization: Bearer <token>` header.
- Return `401 Unauthorized` if the token is missing or invalid.
**Implementation:**
- Use a simple JWT or a static/generated token stored in MongoDB or in-memory.
- Add Echo middleware for token validation on the `/books` routes.
---
### Level 6 — Search & Paginate
| Method | Endpoint                        | Description                    |
|--------|---------------------------------|--------------------------------|
| GET    | `/books?author=X`               | Filter books by author         |
| GET    | `/books?page=1&limit=2`         | Paginate the book list         |
**Plan:**
- Support query parameter `author` to filter books.
- Support `page` and `limit` query parameters for pagination.
- Default values: `page=1`, `limit=10` (or whatever sensible default).
- Return paginated results — likely as `{"data": [...], "page": 1, "limit": 2, "total": N}`.
**MongoDB:** Use `.Find()` with filter + `.Skip()` and `.Limit()` for pagination.
---
### Level 7 — Error Handling
| Method | Endpoint              | Description                                 |
|--------|-----------------------|---------------------------------------------|
| POST   | `/books` (invalid)    | Handle invalid/malformed request body       |
| GET    | `/books/:id` (not found) | Handle non-existent book ID              |
**Plan:**
- POST `/books` with missing/invalid fields should return `400 Bad Request` with a descriptive error message.
- GET `/books/:id` with a non-existent ID should return `404 Not Found`.
- Use proper JSON error response format, e.g.: `{"error": "Book not found"}` or `{"error": "Invalid request body", "details": "..."}`.
---
### Level 8 — Boss: Speed Run
| Method | Endpoint             | Description                        |
|--------|----------------------|------------------------------------|
| ALL    | All previous endpoints | Re-test everything — must be fast |
**Plan:**
- This level re-runs all previous endpoint tests.
- Focus on performance: keep response times low.
- Ensure MongoDB connection pooling is configured.
- Make sure all previous levels still work correctly together.
---
## Project Structure (Golang Echo + MongoDB)
```
api-quest/
├── main.go                  # Entry point, server setup
├── go.mod
├── go.sum
├── config/
│   └── config.go            # Env vars, MongoDB URI, port, JWT secret
├── database/
│   └── mongo.go             # MongoDB connection & client setup
├── models/
│   └── book.go              # Book struct & validation
├── handlers/
│   ├── ping.go              # Level 1: GET /ping
│   ├── echo.go              # Level 2: POST /echo
│   ├── book.go              # Levels 3,4,6,7: CRUD + search + pagination
│   └── auth.go              # Level 5: POST /auth/token
├── middleware/
│   └── auth.go              # JWT/token validation middleware
├── routes/
│   └── routes.go            # Route registration
├── Dockerfile               # For deployment
└── README.md
```
---
## Implementation Plan (Step-by-Step)
### Step 1: Project Setup
- Initialize Go module: `go mod init api-quest`
- Install dependencies:
  - `github.com/labstack/echo/v4` (Echo framework)
  - `go.mongodb.org/mongo-driver/mongo` (MongoDB driver)
  - `github.com/golang-jwt/jwt/v5` (JWT, for auth level)
- Set up MongoDB connection with connection pooling.
- Configure via environment variables: `PORT`, `MONGODB_URI`, `JWT_SECRET`.
### Step 2: Level 1 — Ping
- Register `GET /ping` route.
- Return `{"message": "pong"}` with status `200`.
### Step 3: Level 2 — Echo
- Register `POST /echo` route.
- Bind the JSON body and return it as-is.
### Step 4: Level 3 — CRUD Create & Read
- Create `books` collection in MongoDB.
- Implement `POST /books` — validate, insert, return `201`.
- Implement `GET /books` — find all, return array.
- Implement `GET /books/:id` — find by ID, return book or `404`.
### Step 5: Level 4 — CRUD Update & Delete
- Implement `PUT /books/:id` — find and update, return updated book or `404`.
- Implement `DELETE /books/:id` — find and delete, return success or `404`.
### Step 6: Level 5 — Auth Guard
- Implement `POST /auth/token` — accept credentials, return JWT token.
- Add auth middleware to protect `GET /books` (and possibly other book routes).
- Validate `Authorization: Bearer <token>` header.
### Step 7: Level 6 — Search & Paginate
- Add query param parsing to `GET /books`: `author`, `page`, `limit`.
- Build MongoDB filter dynamically.
- Implement pagination with skip/limit.
- Return paginated response with metadata.
### Step 8: Level 7 — Error Handling
- Add request body validation on `POST /books` (required fields: title, author, etc.).
- Return `400` with descriptive error for invalid payloads.
- Ensure `GET /books/:id` returns `404` for non-existent IDs with proper error JSON.
- Add global error handler in Echo for consistent error responses.
### Step 9: Level 8 — Boss Speed Run
- Optimize MongoDB queries (ensure indexes on `_id`, `author`).
- Verify connection pooling is working.
- Test all endpoints end-to-end for correctness and speed.
- Profile and reduce any unnecessary overhead.
### Step 10: Deploy
- Create a `Dockerfile` for containerized deployment.
- Deploy to Render, Railway, or Fly.io.
- Set environment variables on the platform.
- Paste the public URL into the Desent test page and run.
---
## Rules of Engagement
- Use any language, framework, or tools you want. AI is encouraged.
- Your API must be publicly accessible — they can't test localhost.
- In-memory storage is fine — no database required (but we're using MongoDB).
- Levels must be completed in order. Each one unlocks the next.
- You can retry any failed level as many times as you want.
---
## Key Dependencies
| Package                                | Purpose               |
|----------------------------------------|-----------------------|
| `github.com/labstack/echo/v4`         | HTTP framework        |
| `go.mongodb.org/mongo-driver/mongo`   | MongoDB driver        |
| `github.com/golang-jwt/jwt/v5`        | JWT token generation  |
| `github.com/joho/godotenv`            | .env file loading     |
---
## Environment Variables
```env
PORT=8080
MONGODB_URI=mongodb+srv://<user>:<pass>@cluster.mongodb.net/apiquest?retryWrites=true&w=majority
JWT_SECRET=your-super-secret-key
```
---
## Deployment Options
- **Render** — Free tier, easy Docker deploy
- **Railway** — Fast deploys from GitHub
- **Fly.io** — Good free tier, Dockerfile support
- **Vercel** — Possible with Go serverless functions (less ideal for stateful API)
---
## Notes
- The test runner calls your API in real time; ensure CORS is enabled if needed.
- MongoDB Atlas free tier (M0) works well for this challenge.
- For the speed run (Level 8), response time matters — keep handlers lean and DB queries efficient.
- Add an index on the `author` field in the `books` collection for search performance.