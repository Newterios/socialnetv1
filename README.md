# SocialNet - Social Network Application

A full-featured social network application built with Go backend and React frontend.

## Features

### Core Features
- **User Authentication**: Register, login, JWT-based sessions
- **Google OAuth**: Sign in with Google via Firebase
- **User Profiles**: Bio, avatar, emoji avatars
- **Posts & Feed**: Create, edit, delete posts with media support
- **Social Interactions**: Like, comment, share posts
- **Friends System**: Send/accept friend requests, block users
- **Private Messaging**: Direct messages between users
- **Groups**: Create/join groups, post to groups
- **Notifications**: Real-time activity notifications
- **Admin Panel**: Reports, statistics, user management

### New Features (v2.0)
- **Emoji Avatars**: Choose from 10 predefined emojis, find users with same emoji
- **Privacy Settings**: Control who can message you and see your last seen
- **Online Status**: Real-time online/offline tracking
- **Account Deletion**: Permanently delete your account
- **Admin Statistics**: Site-wide stats dashboard
- **Admin Broadcast**: Send notifications to all users
- **Grant Admin Rights**: Admins can promote other users
- **Group Members**: View member list with join dates
- **Group Settings**: Owners can edit group name/description
- **File Upload**: Image attachments for posts

## Tech Stack

### Backend
- Go 1.21+
- SQLite (primary data)
- MongoDB (analytics)
- Firebase Admin SDK (OAuth)
- JWT Authentication

### Frontend
- React 18 + Vite
- Axios for API calls
- React Router

## Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- MongoDB (optional)
- Firebase project (optional)

### Environment Variables

```bash
# Server
PORT=8080
JWT_SECRET=your-secret-key

# Database
DATABASE_PATH=./data/socialnet.db

# MongoDB (optional)
MONGODB_URI=mongodb://localhost:27017

# Firebase (optional - for Google OAuth)
FIREBASE_KEY_PATH=./firebase-key.json

# Admin settings
INITIAL_ADMINS=admin@example.com,admin2@example.com

# File uploads
UPLOAD_DIR=./uploads
```

### Backend Setup

```bash
cd backend
go mod download
go build -o socialnet .
./socialnet
```

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register new user |
| POST | `/login` | Login with email/password |
| POST | `/auth/google` | Login with Google ID token |
| GET | `/auth/password-requirements` | Get password requirements |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users/{id}` | Get user profile |
| GET | `/users/search?q=` | Search users |
| PUT | `/profile` | Update profile |
| DELETE | `/profile` | Delete account |
| PUT | `/profile/privacy` | Update privacy settings |
| PUT | `/profile/emoji` | Set emoji avatar |
| PUT | `/profile/status` | Update online status |

### Emojis
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/emojis` | Get all available emojis |
| GET | `/emojis/{id}/users` | Get users with specific emoji |

### Posts
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/posts` | Get feed |
| POST | `/posts` | Create post |
| GET | `/posts/{id}` | Get single post |
| PUT | `/posts/{id}` | Update post |
| DELETE | `/posts/{id}` | Delete post |

### Social
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/friends` | Get friends list |
| GET | `/friends/requests` | Get pending requests |
| POST | `/friends/` | Send friend request |
| PUT | `/friends/{id}` | Accept/reject request |
| POST | `/likes/{postId}` | Like a post |
| DELETE | `/likes/{postId}` | Unlike a post |
| GET | `/comments/{postId}` | Get comments |
| POST | `/comments/{postId}` | Add comment |

### Messages
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/conversations` | Get all conversations |
| POST | `/conversations` | Start new conversation |
| GET | `/conversations/{id}/messages` | Get messages |
| POST | `/conversations/{id}/messages` | Send message |

### Groups
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/groups` | Get my groups |
| POST | `/groups` | Create group |
| GET | `/groups/{id}` | Get group details |
| POST | `/groups/{id}/join` | Join group |
| POST | `/groups/{id}/leave` | Leave group |
| GET | `/groups/{id}/posts` | Get group posts |
| POST | `/groups/{id}/posts` | Post to group |
| GET | `/groups/{id}/members` | Get group members |
| PUT | `/groups/{id}/settings` | Update group settings |

### Admin
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/reports` | Get reports (admin) |
| POST | `/reports` | Create report |
| PUT | `/reports/{id}` | Review report |
| DELETE | `/admin/delete/{type}/{id}` | Delete content |
| GET | `/admin/stats` | Get site statistics |
| POST | `/admin/grant` | Grant admin rights |
| POST | `/admin/broadcast` | Send broadcast |
| GET | `/admin/broadcast` | Get broadcasts |

### Upload
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/upload` | Upload file |
| GET | `/uploads/{file}` | Serve uploaded file |

## Project Structure

```
socialnet/
├── main.go                    # Application entry point
├── internal/
│   ├── config/               # Configuration
│   ├── database/             # SQLite & MongoDB
│   ├── http/
│   │   ├── handler/          # HTTP handlers
│   │   ├── middleware/       # Auth, rate limiting
│   │   └── router.go         # Route definitions
│   ├── model/                # Data models
│   ├── repository/           # Data access
│   ├── security/             # JWT, Firebase
│   ├── service/              # Business logic
│   └── worker/               # Background jobs
└── frontend/
    └── src/
        ├── context/          # React contexts
        ├── pages/            # Page components
        └── services/         # API services
```

## Security

- JWT-based authentication with configurable expiration
- Password hashing with bcrypt
- Rate limiting per IP
- Input validation and sanitization
- CORS configuration
- Admin-only endpoints protected

## License

MIT
