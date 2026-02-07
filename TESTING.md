# SocialNet - Testing Guide

## Automated Testing

**Quickest way to test:**

1. Make sure the server is running (`go run main.go`)
2. Open a new terminal
3. Run the test script:

**macOS / Linux (zsh/bash):**
```bash
cd /path/to/AP1_Assignment3_SocialNet
chmod +x test_api.sh
./test_api.sh
```

**Windows (PowerShell):**
```powershell
cd C:\path\to\AP1_Assignment3_SocialNet
.\test_api.ps1
```

The script automatically tests:
- User registration
- Authentication (Login)
- Create post
- Get post
- Feed
- Likes
- Comments
- Profile update
- Friends (requests)
- Groups
- Notifications
- User search

---

## Manual Testing with curl

### 1. Registration
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","username":"testuser","password":"Test123!","full_name":"Test User"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"Test123!"}'
```
Save the token from the response.

### 3. Create Post (replace YOUR_TOKEN)
```bash
curl -X POST http://localhost:8080/posts \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"My first post!"}'
```

### 4. Get Feed
```bash
curl -X GET http://localhost:8080/feed \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Testing with Postman

1. Import the base URL: `http://localhost:8080`
2. Create a Collection with requests from `API.md`
3. Add Environment variable `token` after login
4. Test all endpoints

---

## Database Inspection

After testing, you can view the data:

```bash
sqlite3 socialnet.db
.tables
SELECT * FROM users;
SELECT * FROM posts;
.exit
```

---

## Testing Checklist

### Basic Functionality
- [x] Registration and login
- [x] Create/edit/delete posts
- [x] Likes and comments
- [x] Feed (news feed)

### Social Features
- [x] Send friend requests
- [x] Accept/block requests
- [x] Friends list
- [x] User search

### Messaging
- [x] Create conversation
- [x] Send messages
- [x] Get message history

### Groups
- [x] Create group
- [x] Join/leave group
- [x] Group posts

### Notifications
- [x] Get notification list
- [x] Mark as read
- [x] Unread counter

### Admin Functions
- [x] Create report
- [x] View reports (admin only)
- [x] Delete content (admin only)

---

## Expected Results

- **Server start**: Server running on port 8080
- **Database initialized**: DB created with tables
- **Workers started**: Background processes running
- **API responses**: Correct JSON responses
- **Authentication**: JWT tokens working
- **Authorization**: Permission checks working
- **Validation**: Input data validation
- **Errors**: Clear error messages

---

## Quick Check

Simply open your browser: http://localhost:8080

You should see:
```json
{"message":"SocialNet API","version":"1.0"}
```

Done!
