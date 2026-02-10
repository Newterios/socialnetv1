import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { postsAPI } from '../services/api'
import PostCard from '../components/PostCard'
import CreatePost from '../components/CreatePost'
import './Feed.css'

export default function Feed() {
    const [posts, setPosts] = useState([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        loadFeed()
    }, [])

    const loadFeed = async () => {
        try {
            const res = await postsAPI.getFeed()
            setPosts(res.data || [])
        } catch (err) {
            console.error('Failed to load feed')
        } finally {
            setLoading(false)
        }
    }

    const handlePostCreated = (newPost) => {
        setPosts([newPost, ...posts])
    }

    const handlePostUpdate = (postId, updates) => {
        setPosts(posts.map(p => p.id === postId ? { ...p, ...updates } : p))
    }

    const handlePostDelete = (postId) => {
        setPosts(posts.filter(p => p.id !== postId))
    }

    return (
        <div className="page-container">
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.4 }}
            >
                <h1 className="page-title">Feed</h1>

                <CreatePost onPostCreated={handlePostCreated} />

                <div className="feed-posts">
                    {loading ? (
                        <div className="feed-loading">
                            {[1, 2, 3].map(i => (
                                <div key={i} className="post-skeleton">
                                    <div className="skeleton-header">
                                        <div className="skeleton skeleton-avatar"></div>
                                        <div className="skeleton-info">
                                            <div className="skeleton skeleton-name"></div>
                                            <div className="skeleton skeleton-date"></div>
                                        </div>
                                    </div>
                                    <div className="skeleton skeleton-content"></div>
                                    <div className="skeleton skeleton-actions"></div>
                                </div>
                            ))}
                        </div>
                    ) : posts.length === 0 ? (
                        <motion.div
                            className="feed-empty"
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                        >
                            <div className="feed-empty-icon">
                                <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                                    <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
                                </svg>
                            </div>
                            <h3>No posts yet</h3>
                            <p>Be the first one to share something!</p>
                        </motion.div>
                    ) : (
                        <AnimatePresence mode="popLayout">
                            {posts.map((post, index) => (
                                <motion.div
                                    key={post.id}
                                    initial={{ opacity: 0, y: 20 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    exit={{ opacity: 0, scale: 0.95 }}
                                    transition={{ delay: index * 0.05 }}
                                >
                                    <PostCard
                                        post={post}
                                        onUpdate={handlePostUpdate}
                                        onDelete={handlePostDelete}
                                    />
                                </motion.div>
                            ))}
                        </AnimatePresence>
                    )}
                </div>
            </motion.div>
        </div>
    )
}
