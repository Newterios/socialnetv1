# 2. Architecture & Design (Monolith)

This milestone focuses on planning and architecture, not full implementation.

2.1 High-level system
Client
- Web browser (HTML pages) or REST client (Postman)

Server
- Go web server with RESTful endpoints
- Authentication and authorization
- Business logic services
- Database access layer

Database
- Relational DB (PostgreSQL for production-like setup or SQLite for local development)

2.2 Monolith structure
We keep everything in one repository and one deployable service. Inside the service, we still separate responsibilities by packages.

Suggested Go project structure (design pattern friendly)
- cmd/api
  - application entry point and server startup
- internal/config
  - configuration loading
- internal/http
  - router, middlewares, handlers (controllers)
- internal/service
  - business logic, use cases
- internal/repository
  - database queries and persistence
- internal/model
  - domain models and DTOs
- internal/security
  - password hashing, JWT or session logic
- internal/storage
  - file storage logic (optional)
- pkg
  - reusable helpers that can be used outside internal (optional)

2.3 Data flow
Example: Create Post
- Client sends request to POST /posts
- Handler validates input and user session
- Service checks permissions and creates Post model
- Repository writes Post to database
- Service triggers notifications (if needed)
- Handler returns created post as JSON or redirects to post page

2.4 Main modules and responsibilities
Auth module
- Registration, login, logout
- Session or JWT management
- Password hashing

User module
- Profile view and edit
- User search
- User settings (basic)

Feed and posts module
- Feed generation (simple: latest posts from friends and groups)
- Post CRUD
- Comments and likes

Friends module
- Friend request workflow
- Friends list

Messaging module
- Conversations and messages
- Simple polling based updates (MVP)

Groups module
- Group CRUD
- Group membership

Admin module
- Reports
- Content moderation actions

2.5 Security and validation (planned)
- Input validation for all endpoints
- Password hashing using strong algorithm
- Protection against simple attacks: CSRF for forms, rate limiting for login, secure cookies for sessions
- Authorization checks for edit and delete operations

2.6 Testing strategy (planned)
- Unit tests for services
- Integration tests for repositories (using test database)
- Basic HTTP handler tests with net/http/httptest
