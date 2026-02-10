import { useState, useRef } from 'react'
import { motion } from 'framer-motion'
import { postsAPI, uploadAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'
import { useToast } from './Toast'
import './CreatePost.css'

const EMOJIS = ['😊', '😂', '❤️', '🔥', '👍', '🎉', '💯', '✨', '🙌', '😍']

export default function CreatePost({ onPostCreated }) {
    const { user } = useAuth()
    const toast = useToast()
    const [content, setContent] = useState('')
    const [mediaUrl, setMediaUrl] = useState('')
    const [loading, setLoading] = useState(false)
    const [focused, setFocused] = useState(false)
    const [uploading, setUploading] = useState(false)
    const [showEmojis, setShowEmojis] = useState(false)
    const fileInputRef = useRef(null)

    const handleSubmit = async (e) => {
        e.preventDefault()
        if ((!content.trim() && !mediaUrl) || loading) return

        setLoading(true)
        try {
            const res = await postsAPI.create({ content, media_url: mediaUrl })
            onPostCreated?.(res.data)
            setContent('')
            setMediaUrl('')
            setFocused(false)
            toast.success('Post created!')
        } catch (err) {
            toast.error('Failed to create post')
        } finally {
            setLoading(false)
        }
    }

    const handleImageClick = () => {
        fileInputRef.current?.click()
    }

    const handleFileChange = async (e) => {
        const file = e.target.files?.[0]
        if (!file) return

        if (!file.type.startsWith('image/')) {
            toast.error('Please select an image file')
            return
        }

        if (file.size > 5 * 1024 * 1024) {
            toast.error('Image must be less than 5MB')
            return
        }

        setUploading(true)
        try {
            const res = await uploadAPI.uploadFile(file)
            setMediaUrl(res.data.url)
            toast.success('Image uploaded!')
        } catch (err) {
            toast.error('Failed to upload image')
        } finally {
            setUploading(false)
        }
    }

    const addEmoji = (emoji) => {
        setContent(content + emoji)
        setShowEmojis(false)
    }

    const removeImage = () => {
        setMediaUrl('')
    }

    return (
        <motion.div
            className={`create-post card ${focused ? 'focused' : ''}`}
            layout
        >
            <div className="create-post-header">
                <div className="avatar">
                    {user?.emoji_avatar ? (
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
                            const emoji = EMOJIS.find(e => e.id === user.emoji_avatar)
                            return <span style={{ fontSize: '24px' }}>{emoji?.char || user.emoji_avatar}</span>
                        })()
                    ) : user?.avatar_url ? (
                        <img src={user.avatar_url} alt="" />
                    ) : (
                        user?.username?.charAt(0).toUpperCase()
                    )}
                </div>
                <form className="create-post-form" onSubmit={handleSubmit}>
                    <textarea
                        className="create-post-input"
                        placeholder="What's on your mind?"
                        value={content}
                        onChange={e => setContent(e.target.value)}
                        onFocus={() => setFocused(true)}
                        rows={focused ? 3 : 1}
                    />

                    {mediaUrl && (
                        <div className="create-post-preview">
                            <img src={mediaUrl} alt="Preview" />
                            <button type="button" className="remove-image" onClick={removeImage}>
                                ✕
                            </button>
                        </div>
                    )}

                    <motion.div
                        className="create-post-actions"
                        initial={false}
                        animate={{
                            height: focused ? 'auto' : 0,
                            opacity: focused ? 1 : 0
                        }}
                    >
                        <div className="create-post-tools">
                            <input
                                type="file"
                                ref={fileInputRef}
                                onChange={handleFileChange}
                                accept="image/*"
                                style={{ display: 'none' }}
                            />
                            <button
                                type="button"
                                className="btn btn-icon btn-ghost"
                                title="Add image"
                                onClick={handleImageClick}
                                disabled={uploading}
                            >
                                {uploading ? (
                                    <span className="btn-loader-sm"></span>
                                ) : (
                                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                                        <circle cx="8.5" cy="8.5" r="1.5" />
                                        <polyline points="21 15 16 10 5 21" />
                                    </svg>
                                )}
                            </button>
                            <div className="emoji-picker-wrapper">
                                <button
                                    type="button"
                                    className="btn btn-icon btn-ghost"
                                    title="Add emoji"
                                    onClick={() => setShowEmojis(!showEmojis)}
                                >
                                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                        <circle cx="12" cy="12" r="10" />
                                        <path d="M8 14s1.5 2 4 2 4-2 4-2" />
                                        <line x1="9" y1="9" x2="9.01" y2="9" />
                                        <line x1="15" y1="9" x2="15.01" y2="9" />
                                    </svg>
                                </button>
                                {showEmojis && (
                                    <div className="emoji-picker">
                                        {EMOJIS.map(emoji => (
                                            <button
                                                key={emoji}
                                                type="button"
                                                className="emoji-btn"
                                                onClick={() => addEmoji(emoji)}
                                            >
                                                {emoji}
                                            </button>
                                        ))}
                                    </div>
                                )}
                            </div>
                        </div>
                        <motion.button
                            type="submit"
                            className="btn btn-primary"
                            disabled={(!content.trim() && !mediaUrl) || loading || uploading}
                            whileHover={{ scale: 1.02 }}
                            whileTap={{ scale: 0.98 }}
                        >
                            {loading ? (
                                <span className="btn-loader"></span>
                            ) : (
                                'Post'
                            )}
                        </motion.button>
                    </motion.div>
                </form>
            </div>
        </motion.div>
    )
}
