import { useState } from 'react'
import { Link } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { postsAPI, reportsAPI, socialAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'
import './PostCard.css'

export default function PostCard({ post, onUpdate, onDelete }) {
    const { user } = useAuth()
    const [liked, setLiked] = useState(post.liked)
    const [likeCount, setLikeCount] = useState(post.like_count || 0)
    const [showComments, setShowComments] = useState(false)
    const [comments, setComments] = useState([])
    const [newComment, setNewComment] = useState('')
    const [loadingComments, setLoadingComments] = useState(false)
    const [showMenu, setShowMenu] = useState(false)

    const isOwner = user?.id === post.user_id || user?.id === post.author?.id

    const handleLike = async () => {
        try {
            if (liked) {
                await socialAPI.unlikePost(post.id)
                setLiked(false)
                setLikeCount(c => c - 1)
            } else {
                await socialAPI.likePost(post.id)
                setLiked(true)
                setLikeCount(c => c + 1)
            }
        } catch (err) {
            console.error('Like failed')
        }
    }

    const loadComments = async () => {
        if (comments.length > 0) {
            setShowComments(!showComments)
            return
        }

        setLoadingComments(true)
        try {
            const res = await socialAPI.getComments(post.id)
            setComments(res.data || [])
            setShowComments(true)
        } catch (err) {
            console.error('Failed to load comments')
        } finally {
            setLoadingComments(false)
        }
    }

    const handleComment = async (e) => {
        e.preventDefault()
        if (!newComment.trim()) return

        try {
            const res = await socialAPI.addComment(post.id, { content: newComment })
            const comment = res.data || {}
            if (!comment.author && user) {
                comment.author = user
            }
            setComments([...comments, comment])
            setNewComment('')
        } catch (err) {
            console.error('Failed to add comment')
        }
    }

    const handleDelete = async () => {
        if (!confirm('Are you sure you want to delete this post?')) return

        try {
            await postsAPI.deletePost(post.id)
            onDelete?.(post.id)
        } catch (err) {
            console.error('Failed to delete post')
        }
    }

    const handleReport = async (targetType, targetId) => {
        const reason = window.prompt('Describe the issue')
        if (!reason || !reason.trim()) return

        try {
            await reportsAPI.create(targetType, targetId, reason.trim())
            setShowMenu(false)
        } catch (err) {
            console.error('Failed to report content')
        }
    }

    const formatDate = (dateStr) => {
        const date = new Date(dateStr)
        const now = new Date()
        const diff = now - date
        const mins = Math.floor(diff / 60000)
        const hours = Math.floor(diff / 3600000)
        const days = Math.floor(diff / 86400000)

        if (mins < 1) return 'Just now'
        if (mins < 60) return `${mins}m ago`
        if (hours < 24) return `${hours}h ago`
        if (days < 7) return `${days}d ago`
        return date.toLocaleDateString()
    }

    const author = post.author || { username: 'user', full_name: 'User' }

    return (
        <motion.article
            className="post-card card"
            layout
        >
            <div className="post-header">
                <Link to={`/profile/${author.id || post.user_id}`} className="post-author">
                    <div className="avatar">
                        {author.emoji_avatar ? (
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
                                const emoji = EMOJIS.find(e => e.id === author.emoji_avatar)
                                return <span style={{ fontSize: '24px' }}>{emoji?.char || author.emoji_avatar}</span>
                            })()
                        ) : author.avatar_url ? (
                            <img src={author.avatar_url} alt="" />
                        ) : (
                            author.username?.charAt(0).toUpperCase()
                        )}
                    </div>
                    <div className="post-author-info">
                        <span className="post-author-name">{author.full_name}</span>
                        <span className="post-author-meta">
                            @{author.username} · {formatDate(post.created_at)}
                        </span>
                    </div>
                </Link>

                <div className="post-menu">
                    <button
                        className="btn btn-icon btn-ghost"
                        onClick={() => setShowMenu(!showMenu)}
                    >
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                            <circle cx="12" cy="5" r="2" />
                            <circle cx="12" cy="12" r="2" />
                            <circle cx="12" cy="19" r="2" />
                        </svg>
                    </button>

                    <AnimatePresence>
                        {showMenu && (
                            <motion.div
                                className="post-menu-dropdown"
                                initial={{ opacity: 0, scale: 0.95, y: -10 }}
                                animate={{ opacity: 1, scale: 1, y: 0 }}
                                exit={{ opacity: 0, scale: 0.95, y: -10 }}
                            >
                                {!isOwner && (
                                    <button className="post-menu-item warning" onClick={() => handleReport('post', post.id)}>
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
                                            <line x1="12" y1="9" x2="12" y2="13" />
                                            <line x1="12" y1="17" x2="12.01" y2="17" />
                                        </svg>
                                        Report post
                                    </button>
                                )}
                                {isOwner && (
                                    <button className="post-menu-item danger" onClick={handleDelete}>
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                            <polyline points="3 6 5 6 21 6" />
                                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                                        </svg>
                                        Delete
                                    </button>
                                )}
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            </div>

            <div className="post-content">
                <p>{post.content}</p>
                {post.media_url && (
                    <div className="post-media">
                        <img src={post.media_url} alt="" />
                    </div>
                )}
            </div>

            <div className="post-actions">
                <motion.button
                    className={`post-action ${liked ? 'liked' : ''}`}
                    onClick={handleLike}
                    whileTap={{ scale: 0.9 }}
                >
                    <motion.svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill={liked ? 'currentColor' : 'none'}
                        stroke="currentColor"
                        strokeWidth="2"
                        animate={liked ? { scale: [1, 1.3, 1] } : {}}
                        transition={{ duration: 0.3 }}
                    >
                        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
                    </motion.svg>
                    <span>{likeCount}</span>
                </motion.button>

                <button
                    className={`post-action ${showComments ? 'active' : ''}`}
                    onClick={loadComments}
                >
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                    </svg>
                    <span>{comments.length || ''}</span>
                </button>

                <button className="post-action" onClick={async () => {
                    const url = `https://socialnet-v1.onrender.com/post/${post.id}`
                    try {
                        await navigator.clipboard.writeText(url)
                        alert('Post link copied to clipboard!')
                    } catch (err) {
                        console.error('Failed to copy link')
                    }
                }}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <circle cx="18" cy="5" r="3" />
                        <circle cx="6" cy="12" r="3" />
                        <circle cx="18" cy="19" r="3" />
                        <line x1="8.59" y1="13.51" x2="15.42" y2="17.49" />
                        <line x1="15.41" y1="6.51" x2="8.59" y2="10.49" />
                    </svg>
                </button>
            </div>

            <AnimatePresence>
                {showComments && (
                    <motion.div
                        className="post-comments"
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                    >
                        {loadingComments ? (
                            <div className="comments-loading">Loading comments...</div>
                        ) : (
                            <>
                                <div className="comments-list">
                                    {comments.map(comment => (
                                        <div key={comment.id} className="comment">
                                            <div className="avatar avatar-sm">
                                                {comment.author?.avatar_url ? (
                                                    <img src={comment.author.avatar_url} alt="" />
                                                ) : (
                                                    comment.author?.username?.charAt(0).toUpperCase() || 'U'
                                                )}
                                            </div>
                                            <div className="comment-body">
                                                <div className="comment-header">
                                                    <span className="comment-author">{comment.author?.full_name || comment.author?.username || 'User'}</span>
                                                    <div className="comment-meta">
                                                        <span className="comment-time">{formatDate(comment.created_at)}</span>
                                                        {comment.user_id !== user?.id && (
                                                            <button
                                                                className="comment-report"
                                                                onClick={() => handleReport('comment', comment.id)}
                                                            >
                                                                Report
                                                            </button>
                                                        )}
                                                    </div>
                                                </div>
                                                <p className="comment-content">{comment.content}</p>
                                            </div>
                                        </div>
                                    ))}
                                </div>

                                <form className="comment-form" onSubmit={handleComment}>
                                    <input
                                        type="text"
                                        className="input-field"
                                        placeholder="Write a comment..."
                                        value={newComment}
                                        onChange={e => setNewComment(e.target.value)}
                                    />
                                    <button type="submit" className="btn btn-primary btn-sm" disabled={!newComment.trim()}>
                                        Post
                                    </button>
                                </form>
                            </>
                        )}
                    </motion.div>
                )}
            </AnimatePresence>
        </motion.article>
    )
}
