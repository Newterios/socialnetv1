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
