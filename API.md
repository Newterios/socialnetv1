# SocialNet API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication

All endpoints except `/register` and `/login` require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Response Format

Success responses return JSON with relevant data.
Error responses return JSON with error message and appropriate HTTP status code.

## Endpoints

### Authentication

#### Register User
```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "Password123!",
  "full_name": "Full Name"
}

Response: 201 Created
{
  "id": 1,
  "email": "user@example.com",
  "username": "username",
  "full_name": "Full Name",
  "bio": "",
  "avatar_url": "",
  "is_admin": false,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password123!"
}

Response: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "full_name": "Full Name",
    ...
  }
}
```

### Users

#### Get User Profile
```http
GET /users/:id
Authorization: Bearer <token>

Response: 200 OK
{
  "id": 1,
  "email": "user@example.com",
  "username": "username",
  "full_name": "Full Name",
  "bio": "User bio",
  "avatar_url": "https://...",
  "is_admin": false,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Update Profile
```http
PUT /users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "full_name": "Updated Name",
  "bio": "New bio",
  "avatar_url": "https://..."
}

Response: 200 OK
{"message": "profile updated"}
```

#### Search Users
```http
GET /users/search?q=john
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 2,
    "username": "john_doe",
    "full_name": "John Doe",
    ...
  }
]
```

### Posts

#### Create Post
```http
POST /posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "This is my post content",
  "media_url": "https://..."
}

Response: 201 Created
{
  "id": 1,
  "user_id": 1,
  "content": "This is my post content",
  "media_url": "https://...",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "author": {...},
  "like_count": 0,
  "liked": false
}
```

#### Get Post
```http
GET /posts/:id
Authorization: Bearer <token>

Response: 200 OK
{
  "id": 1,
  "user_id": 1,
  "content": "Post content",
  "author": {...},
  "like_count": 5,
  "liked": true,
  ...
}
```

#### Update Post
```http
PUT /posts/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Updated content",
  "media_url": "https://..."
}

Response: 200 OK
{"message": "post updated"}
```

#### Delete Post
```http
DELETE /posts/:id
Authorization: Bearer <token>

Response: 200 OK
{"message": "post deleted"}
```

#### Get Feed
```http
GET /feed
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "content": "Post content",
    "author": {...},
    "like_count": 5,
    ...
  }
]
```

### Social Features

#### Like Post
```http
POST /posts/:id/like
Authorization: Bearer <token>

Response: 201 Created
{"message": "post liked"}
```

#### Unlike Post
```http
DELETE /posts/:id/like
Authorization: Bearer <token>

Response: 200 OK
{"message": "post unliked"}
```

#### Comment on Post
```http
POST /posts/:id/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Great post!"
}

Response: 201 Created
{
  "id": 1,
  "post_id": 1,
  "user_id": 2,
  "content": "Great post!",
  "created_at": "2024-01-01T00:00:00Z",
  "author": {...}
}
```

#### Get Comments
```http
GET /posts/:id/comments
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "post_id": 1,
    "content": "Great post!",
    "author": {...},
    ...
  }
]
```

### Friends

#### Send Friend Request
```http
POST /friends/request
Authorization: Bearer <token>
Content-Type: application/json

{
  "addressee_id": 2
}

Response: 201 Created
{"message": "friend request sent"}
```

#### Get Pending Requests
```http
GET /friends/pending
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "requester_id": 2,
    "status": "pending",
    "requester": {...},
    ...
  }
]
```

#### Accept Friend Request
```http
PUT /friends/:id/accept
Authorization: Bearer <token>

Response: 200 OK
{"message": "friend request accepted"}
```

#### Block User
```http
PUT /friends/:id/block
Authorization: Bearer <token>

Response: 200 OK
{"message": "user blocked"}
```

#### Get Friends List
```http
GET /friends
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 2,
    "username": "friend1",
    "full_name": "Friend One",
    ...
  }
]
```

### Messaging

#### Start Conversation
```http
POST /conversations
Authorization: Bearer <token>
Content-Type: application/json

{
  "participant_id": 2
}

Response: 201 Created
{
  "id": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Send Message
```http
POST /conversations/:id/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "body": "Hello there!"
}

Response: 201 Created
{
  "id": 1,
  "conversation_id": 1,
  "user_id": 1,
  "body": "Hello there!",
  "created_at": "2024-01-01T00:00:00Z",
  ...
}
```

#### Get Messages
```http
GET /conversations/:id/messages
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "body": "Hello there!",
    "author": {...},
    ...
  }
]
```

#### Get Conversations
```http
GET /conversations
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    ...
  }
]
```

### Groups

#### Create Group
```http
POST /groups
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My Group",
  "description": "Group description"
}

Response: 201 Created
{
  "id": 1,
  "owner_id": 1,
  "title": "My Group",
  "description": "Group description",
  "created_at": "2024-01-01T00:00:00Z",
  "owner": {...},
  "member_count": 1,
  "is_member": true
}
```

#### Join Group
```http
POST /groups/:id/join
Authorization: Bearer <token>

Response: 200 OK
{"message": "joined group"}
```

#### Leave Group
```http
DELETE /groups/:id/leave
Authorization: Bearer <token>

Response: 200 OK
{"message": "left group"}
```

#### Post to Group
```http
POST /groups/:id/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Group post content"
}

Response: 201 Created
{
  "id": 1,
  "group_id": 1,
  "user_id": 1,
  "content": "Group post content",
  "created_at": "2024-01-01T00:00:00Z",
  ...
}
```

### Notifications

#### Get Notifications
```http
GET /notifications
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "user_id": 1,
    "type": "like",
    "message": "user123 liked your post",
    "read": false,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Mark as Read
```http
PUT /notifications/:id/read
Authorization: Bearer <token>

Response: 200 OK
{"message": "notification marked as read"}
```

### Admin

#### Create Report
```http
POST /reports
Authorization: Bearer <token>
Content-Type: application/json

{
  "target_type": "post",
  "target_id": 1,
  "reason": "Spam content"
}

Response: 201 Created
{"message": "report created"}
```

#### Get Reports (Admin Only)
```http
GET /admin/reports?status=pending
Authorization: Bearer <admin_token>

Response: 200 OK
[
  {
    "id": 1,
    "reporter_id": 2,
    "target_type": "post",
    "target_id": 1,
    "reason": "Spam content",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z",
    "reporter": {...}
  }
]
```

#### Delete Content (Admin Only)
```http
DELETE /admin/content/post/1
Authorization: Bearer <admin_token>

Response: 200 OK
{"message": "content deleted"}
```

## Error Codes

- `400 Bad Request` - Invalid request format or validation error
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
