import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { motion } from 'framer-motion'
import { usersAPI, postsAPI, friendsAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'
import PostCard from '../components/PostCard'
import './Profile.css'

export default function Profile() {
    const { id } = useParams()
    const { user: currentUser } = useAuth()
    const [profile, setProfile] = useState(null)
    const [posts, setPosts] = useState([])
    const [loading, setLoading] = useState(true)
    const [friendStatus, setFriendStatus] = useState(null)

    const isOwn = currentUser?.id === parseInt(id)

    useEffect(() => {
        loadProfile()
        if (!isOwn) {
            loadFriendStatus()
        }
    }, [id])

    const loadProfile = async () => {
        setLoading(true)
        try {
            const [profileRes, feedRes] = await Promise.all([
                usersAPI.getProfile(id),
                postsAPI.getFeed()
            ])
            setProfile(profileRes.data)
            const userPosts = (feedRes.data || []).filter(p =>
                p.user_id === parseInt(id) || p.author?.id === parseInt(id)
            )
            setPosts(userPosts)
        } catch (err) {
            console.error('Failed to load profile')
        } finally {
            setLoading(false)
        }
    }

    const loadFriendStatus = async () => {
        try {
            const res = await friendsAPI.getList()
            const friends = res.data || []
            const isFriend = friends.some(f => f.id === parseInt(id))
            if (isFriend) {
                setFriendStatus('accepted')
            }
        } catch (err) {
            console.error('Failed to load friend status')
        }
    }

    const handleSendFriendRequest = async () => {
        try {
            await friendsAPI.sendRequest(parseInt(id))
            setFriendStatus('pending')
        } catch (err) {
            console.error('Failed to send friend request')
        }
    }

    if (loading) {
        return (
            <div className="page-container">
                <div className="profile-skeleton">
                    <div className="skeleton profile-cover-skeleton"></div>
                    <div className="profile-header-skeleton">
                        <div className="skeleton profile-avatar-skeleton"></div>
                        <div className="skeleton profile-name-skeleton"></div>
                    </div>
                </div>
            </div>
        )
    }

    if (!profile) {
        return (
            <div className="page-container">
                <div className="profile-not-found">
                    <h2>User not found</h2>
                </div>
            </div>
        )
    }

    return (
        <div className="profile-page">
            <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.4 }}
            >
                <div className="profile-cover">
                    <div className="profile-cover-gradient"></div>
                </div>

                <div className="profile-container">
                    <div className="profile-header">
                        <motion.div
                            className="profile-avatar-wrapper"
                            initial={{ scale: 0.8, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            transition={{ delay: 0.1 }}
                        >
                            <div className="avatar avatar-lg profile-avatar">
                                {profile.emoji_avatar ? (
                                    (() => {
                                        const EMOJIS = [
                                            { id: 'smile', char: '😊' },
                                            { id: 'cool', char: '😎' },
                                            { id: 'fire', char: '🔥' },
                                            { id: 'lightning', char: '⚡' },
                                            { id: 'gaming', char: '🎮' },
                                            { id: 'art', char: '🎨' },
                                            { id: 'music', char: '🎵' },
                                            { id: 'book', char: '📚' },
                                            { id: 'star', char: '🌟' },
                                            { id: 'strong', char: '💪' },
                                        ]
                                        const emoji = EMOJIS.find(e => e.id === profile.emoji_avatar)
                                        return <span style={{ fontSize: '64px' }}>{emoji?.char || profile.emoji_avatar}</span>
                                    })()
                                ) : profile.avatar_url ? (
                                    <img src={profile.avatar_url} alt="" />
                                ) : (
                                    profile.username?.charAt(0).toUpperCase()
                                )}
                            </div>
                        </motion.div>

                        <div className="profile-info">
                            <h1 className="profile-name">{profile.full_name}</h1>
                            <p className="profile-username">@{profile.username}</p>
                            {profile.bio && <p className="profile-bio">{profile.bio}</p>}

                            <div className="profile-stats">
                                <div className="profile-stat">
                                    <span className="profile-stat-value">{posts.length}</span>
                                    <span className="profile-stat-label">Posts</span>
                                </div>
                                <div className="profile-stat">
                                    <span className="profile-stat-value">0</span>
                                    <span className="profile-stat-label">Friends</span>
                                </div>
                            </div>
                        </div>

                        <div className="profile-actions">
                            {!isOwn && friendStatus !== 'accepted' && (
                                <motion.button
                                    className="btn btn-primary"
                                    onClick={handleSendFriendRequest}
                                    disabled={friendStatus === 'pending'}
                                    whileHover={{ scale: 1.02 }}
                                    whileTap={{ scale: 0.98 }}
                                >
                                    {friendStatus === 'pending' ? 'Request Sent' : 'Add Friend'}
                                </motion.button>
                            )}
                        </div>
                    </div>

                    <div className="profile-content">
                        <div className="profile-section">
                            <h2 className="profile-section-title">Posts</h2>

                            {posts.length === 0 ? (
                                <div className="profile-empty">
                                    <p>No posts yet</p>
                                </div>
                            ) : (
                                <div className="profile-posts">
                                    {posts.map((post, index) => (
                                        <motion.div
                                            key={post.id}
                                            initial={{ opacity: 0, y: 20 }}
                                            animate={{ opacity: 1, y: 0 }}
                                            transition={{ delay: index * 0.05 }}
                                        >
                                            <PostCard post={post} />
                                        </motion.div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </motion.div>
        </div>
    )
}
