# Assignment 3 - Project Design & Setup Milestone Report

## Project Title
SocialNet - Social Network (VK/Facebook-like)

## Team
- Aitbek Nugmanov
- Syrym Shadiyarbek
- Rakhat Balgabekov

## Repository
GitHub: https://github.com/newterios/socialnet


---


## 1. Project Proposal

# 1. Project Proposal

Project title
SocialNet - Social Network (VK and Facebook style)

1.1 Relevance
Social networks are used daily for communication, content sharing, and community building. Our project is a learning-focused social network that demonstrates real backend patterns in Go and prepares us for building RESTful web applications.

1.2 Problem
Many social platforms are overloaded with features and have complex settings. For our educational project, we focus on core functions that most users need: accounts, profiles, posts, feed, friends, messaging, groups, and notifications.

1.3 Competitor analysis
VK
- Strong communities, posts, and messaging
- Many features, can be heavy for new users

Facebook
- Powerful feed and groups
- Complex privacy and settings

Instagram
- Focus on media posts and stories
- Less focus on long text discussions

Telegram channels and chats
- Fast messaging and communities
- Not a classic profile and feed social network

Our direction
- Classic profile-based social network with posts, comments, friends, and messages
- Simple scope with clear data model and clean monolith architecture

1.4 Target users
- University students and clubs
- Small communities that want a simple internal social network
- People who prefer text posts and community groups

1.5 Planned features (MVP for final defense)
Accounts and profiles
- Registration, login, logout
- Profile page with avatar and bio
- Edit profile

Posts and feed
- Create, edit, delete post
- Upload media for post (optional)
- Feed: recent posts from friends and groups
- Like, comment

Friends
- Send friend request
- Accept or decline
- List of friends

Messaging
- Private chat between two users
- Send and read messages
- Message history

Groups
- Create group
- Join or leave group
- Group feed

Notifications
- Friend request notification
- Like or comment notification
- New message notification

Admin and moderation (basic)
- Report post or comment
- Admin can delete reported content

1.6 Non-goals for this milestone
- Full UI design
- Advanced recommendation algorithms
- Real-time websockets in MVP (can be a bonus later)

1.7 Assumptions
- Monolith backend is used for now
- Web interface can be simple HTML templates or minimal frontend
- Database is relational (PostgreSQL or SQLite) because ERD fits well


## 2. Architecture and Design

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


## 3. Diagrams

All diagrams are exported as PNG into `docs/diagrams/exported/`.

### 3.1 Use Case Diagram
File: `docs/diagrams/exported/use_case_socialnet.png`

### 3.2 ER Diagram
File: `docs/diagrams/exported/erd_socialnet.png`

### 3.3 UML Class Diagram
File: `docs/diagrams/exported/uml_class_socialnet.png`

### 3.4 Sequence Diagram (Create Post)
File: `docs/diagrams/exported/sequence_create_post.png`



## 4. Project Plan (Gantt - Weeks 7-10)

# 4. Project Plan (Gantt) - Weeks 7-10

Rules we follow
- Plan is for weeks 7-10 only
- Tasks are split evenly between team members
- Final polish tasks are not listed yet, as required

Team members
- Aitbek Nugmanov
- Syrym Shadiyarbek
- Rakhat Balgabekov

4.1 Timeline overview
Week 7
- Finalize requirements and data model
- Create repository structure and skeleton
- Prepare diagrams and API draft

Week 8
- Implement authentication and user profiles
- Implement posts and comments base endpoints
- First defense preparation and demo

Week 9
- Implement friends and messaging
- Implement groups
- Add basic security and tests

Week 10
- Stabilize core flows
- Add admin reporting and moderation
- Deployment draft (optional)

4.2 Gantt table (simple)
| Week | Aitbek Nugmanov | Syrym Shadiyarbek | Rakhat Balgabekov |
| - | - | - | - |
| 7 | Repository setup, project structure, CI draft | ERD and data model, DB schema draft | Use Case and UML diagrams, API draft |
| 8 | Auth module + session, profile edit | Posts CRUD + feed draft | Comments and likes + validation |
| 9 | Friends workflow + endpoints | Messaging module + endpoints | Groups module + endpoints |
| 10 | Admin reports + moderation | Security improvements + testing | Deployment draft + integration checks |

4.3 Deliverables for Assignment 3 (this milestone)
- Project proposal markdown
- Architecture description and modules
- Diagrams (Use Case, ERD, UML)
- Gantt plan for weeks 7-10
- Git repository with branches and commits per member


## 5. Team Contribution Summary

# 6. Team Work Distribution (Equal Contribution)

We split the work so that each member has clear ownership, and we can defend any part.

Aitbek Nugmanov
- Repository creation and structure
- Architecture description and module boundaries
- Git setup and workflow guide
- Draft of authentication and profile flows

Syrym Shadiyarbek
- ERD design and database entities
- Data model rules (constraints, relations)
- Draft of repositories and persistence layer approach
- Week 9 messaging module plan

Rakhat Balgabekov
- Use Case diagram and actors
- UML class diagram and sequence diagram
- API draft for posts, comments, likes
- Week 10 deployment draft ideas

Joint tasks
- Review documents for consistent English B2 level
- Prepare defense answers and demo scenario


## 6. Current Status

At this milestone, the project is initialized with Go module `socialnet`, basic server entry point, and the documentation package required for Assignment 3.

The next milestone will implement core features: authentication, profiles, posts, and feeds.
