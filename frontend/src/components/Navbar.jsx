import { useState, useEffect, useRef } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { useAuth } from '../context/AuthContext'
import { notificationsAPI, usersAPI } from '../services/api'
import './Navbar.css'

export default function Navbar() {
    const { user, logout } = useAuth()
    const navigate = useNavigate()
    const [showNotifications, setShowNotifications] = useState(false)
    const [showProfile, setShowProfile] = useState(false)
    const [showSearch, setShowSearch] = useState(false)
    const [notifications, setNotifications] = useState([])
    const [searchQuery, setSearchQuery] = useState('')
    const [searchResults, setSearchResults] = useState([])
    const [clearingNotifications, setClearingNotifications] = useState(false)
    const searchRef = useRef(null)

    useEffect(() => {
        loadNotifications()
    }, [])

    useEffect(() => {
        const handleClickOutside = (e) => {
            if (searchRef.current && !searchRef.current.contains(e.target)) {
                setShowSearch(false)
            }
        }
        document.addEventListener('mousedown', handleClickOutside)
        return () => document.removeEventListener('mousedown', handleClickOutside)
    }, [])

    const loadNotifications = async () => {
        try {
            const res = await notificationsAPI.getList()
            setNotifications(res.data || [])
        } catch (err) {
            console.error('Failed to load notifications')
        }
    }

    const handleSearch = async (e) => {
        const query = e.target.value
        setSearchQuery(query)

        if (query.length >= 2) {
            try {
                const res = await usersAPI.search(query)
                setSearchResults(res.data || [])
            } catch (err) {
                setSearchResults([])
            }
        } else {
            setSearchResults([])
        }
    }

    const handleLogout = () => {
        logout()
        navigate('/login')
    }

    const handleClearNotifications = async () => {
        setClearingNotifications(true)
        try {
            await notificationsAPI.clearAll()
            setNotifications([])
        } catch (err) {
            console.error('Failed to clear notifications')
        } finally {
            setClearingNotifications(false)
        }
    }

    const unreadCount = notifications.filter(n => !n.read).length

    return (
        <nav className="navbar">
            <div className="navbar-content">
                <Link to="/" className="navbar-logo">
                    <div className="navbar-logo-icon">
                        <svg width="32" height="32" viewBox="0 0 48 48" fill="none">
                            <circle cx="24" cy="24" r="24" fill="url(#navLogoGradient)" />
                            <path d="M16 24C16 19.5817 19.5817 16 24 16C28.4183 16 32 19.5817 32 24" stroke="white" strokeWidth="2.5" strokeLinecap="round" />
                            <circle cx="24" cy="30" r="4" fill="white" />
                            <defs>
                                <linearGradient id="navLogoGradient" x1="0" y1="0" x2="48" y2="48">
                                    <stop stopColor="#8b5cf6" />
                                    <stop offset="1" stopColor="#06b6d4" />
                                </linearGradient>
                            </defs>
                        </svg>
                    </div>
                    <span className="navbar-logo-text">SocialNet</span>
                </Link>

                <div className="navbar-search" ref={searchRef}>
                    <div className="search-input-wrapper">
                        <svg className="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <circle cx="11" cy="11" r="8" />
                            <path d="M21 21l-4.35-4.35" />
                        </svg>
                        <input
                            type="text"
                            className="search-input"
                            placeholder="Search users..."
                            value={searchQuery}
                            onChange={handleSearch}
                            onFocus={() => setShowSearch(true)}
                        />
                    </div>

                    <AnimatePresence>
                        {showSearch && searchResults.length > 0 && (
                            <motion.div
                                className="search-dropdown"
                                initial={{ opacity: 0, y: -10 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -10 }}
                            >
                                {searchResults.map(u => (
                                    <Link
                                        key={u.id}
                                        to={`/profile/${u.id}`}
                                        className="search-result"
                                        onClick={() => {
                                            setShowSearch(false)
                                            setSearchQuery('')
                                        }}
                                    >
                                        <div className="avatar avatar-sm">
                                            {u.avatar_url ? (
                                                <img src={u.avatar_url} alt="" />
                                            ) : (
                                                u.username?.charAt(0).toUpperCase()
                                            )}
                                        </div>
                                        <div className="search-result-info">
                                            <span className="search-result-name">{u.full_name}</span>
                                            <span className="search-result-username">@{u.username}</span>
                                        </div>
                                    </Link>
                                ))}
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>

                <div className="navbar-actions">
                    <div className="navbar-notification">
                        <button
                            className="btn btn-icon btn-ghost"
                            onClick={() => setShowNotifications(!showNotifications)}
                        >
                            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                                <path d="M13.73 21a2 2 0 0 1-3.46 0" />
                            </svg>
                            {unreadCount > 0 && (
                                <span className="notification-badge">{unreadCount}</span>
                            )}
                        </button>

                        <AnimatePresence>
                            {showNotifications && (
                                <motion.div
                                    className="notifications-dropdown"
                                    initial={{ opacity: 0, y: -10, scale: 0.95 }}
                                    animate={{ opacity: 1, y: 0, scale: 1 }}
                                    exit={{ opacity: 0, y: -10, scale: 0.95 }}
                                >
                                    <div className="notifications-header">
                                        <h3>Notifications</h3>
                                        <div className="notifications-header-actions">
                                            {unreadCount > 0 && <span className="badge badge-primary">{unreadCount} new</span>}
                                            {notifications.length > 0 && (
                                                <button
                                                    className="notifications-clear"
                                                    onClick={handleClearNotifications}
                                                    disabled={clearingNotifications}
                                                >
                                                    {clearingNotifications ? 'Clearing...' : 'Clear'}
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                    <div className="notifications-list">
                                        {notifications.length === 0 ? (
                                            <div className="notifications-empty">No notifications yet</div>
                                        ) : (
                                            notifications.slice(0, 5).map(n => (
                                                <div key={n.id} className={`notification-item ${!n.read ? 'unread' : ''}`}>
                                                    <p>{n.message}</p>
                                                    <span className="notification-time">
                                                        {new Date(n.created_at).toLocaleDateString()}
                                                    </span>
                                                </div>
                                            ))
                                        )}
                                    </div>
                                </motion.div>
                            )}
                        </AnimatePresence>
                    </div>

                    <div className="navbar-profile">
                        <button
                            className="profile-trigger"
                            onClick={() => setShowProfile(!showProfile)}
                        >
                            <div className="avatar avatar-sm">
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
                                        return <span style={{ fontSize: '20px' }}>{emoji?.char || user.emoji_avatar}</span>
                                    })()
                                ) : user?.avatar_url ? (
                                    <img src={user.avatar_url} alt="" />
                                ) : (
                                    user?.username?.charAt(0).toUpperCase()
                                )}
                            </div>
                        </button>

                        <AnimatePresence>
                            {showProfile && (
                                <motion.div
                                    className="profile-dropdown"
                                    initial={{ opacity: 0, y: -10, scale: 0.95 }}
                                    animate={{ opacity: 1, y: 0, scale: 1 }}
                                    exit={{ opacity: 0, y: -10, scale: 0.95 }}
                                >
                                    <div className="profile-dropdown-header">
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
                                                    return <span style={{ fontSize: '32px' }}>{emoji?.char || user.emoji_avatar}</span>
                                                })()
                                            ) : user?.avatar_url ? (
                                                <img src={user.avatar_url} alt="" />
                                            ) : (
                                                user?.username?.charAt(0).toUpperCase()
                                            )}
                                        </div>
                                        <div className="profile-dropdown-info">
                                            <span className="profile-dropdown-name">{user?.full_name}</span>
                                            <span className="profile-dropdown-email">{user?.email}</span>
                                        </div>
                                    </div>
                                    <div className="divider"></div>
                                    <Link to={`/profile/${user?.id}`} className="profile-dropdown-item" onClick={() => setShowProfile(false)}>
                                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                                            <circle cx="12" cy="7" r="4" />
                                        </svg>
                                        View Profile
                                    </Link>
                                    <Link to="/settings" className="profile-dropdown-item" onClick={() => setShowProfile(false)}>
                                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                            <circle cx="12" cy="12" r="3" />
                                            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
                                        </svg>
                                        Settings
                                    </Link>
                                    <div className="divider"></div>
                                    <button className="profile-dropdown-item logout" onClick={handleLogout}>
                                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                                            <polyline points="16 17 21 12 16 7" />
                                            <line x1="21" y1="12" x2="9" y2="12" />
                                        </svg>
                                        Sign Out
                                    </button>
                                </motion.div>
                            )}
                        </AnimatePresence>
                    </div>
                </div>
            </div>
        </nav>
    )
}
