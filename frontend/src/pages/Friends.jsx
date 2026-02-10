import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { friendsAPI, messagesAPI } from '../services/api'
import './Friends.css'

export default function Friends() {
    const navigate = useNavigate()
    const [friends, setFriends] = useState([])
    const [pending, setPending] = useState([])
    const [activeTab, setActiveTab] = useState('friends')
    const [loading, setLoading] = useState(true)
    const [startingConversationWith, setStartingConversationWith] = useState(null)

    useEffect(() => {
        loadData()
    }, [])

    const loadData = async () => {
        try {
            const [friendsRes, pendingRes] = await Promise.all([
                friendsAPI.getList(),
                friendsAPI.getPending()
            ])
            setFriends(friendsRes.data || [])
            setPending(pendingRes.data || [])
        } catch (err) {
            console.error('Failed to load friends data')
        } finally {
            setLoading(false)
        }
    }

    const handleAccept = async (requestId) => {
        try {
            await friendsAPI.accept(requestId)
            const acceptedRequest = pending.find(p => p.id === requestId)
            setPending(pending.filter(p => p.id !== requestId))
            if (acceptedRequest?.requester) {
                setFriends([...friends, acceptedRequest.requester])
            }
        } catch (err) {
            console.error('Failed to accept request')
        }
    }

    const handleBlock = async (requestId) => {
        try {
            await friendsAPI.block(requestId)
            setPending(pending.filter(p => p.id !== requestId))
        } catch (err) {
            console.error('Failed to block user')
        }
    }

    const handleMessage = async (friend) => {
        setStartingConversationWith(friend.id)
        try {
            const res = await messagesAPI.createConversation(friend.id)
            navigate('/messages', {
                state: {
                    conversationId: res.data?.id,
                    participant: friend
                }
            })
        } catch (err) {
            console.error('Failed to start conversation')
        } finally {
            setStartingConversationWith(null)
        }
    }

    return (
        <div className="page-container">
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
            >
                <h1 className="page-title">Friends</h1>

                <div className="friends-tabs">
                    <button
                        className={`tab ${activeTab === 'friends' ? 'active' : ''}`}
                        onClick={() => setActiveTab('friends')}
                    >
                        Friends
                        {friends.length > 0 && <span className="tab-count">{friends.length}</span>}
                    </button>
                    <button
                        className={`tab ${activeTab === 'pending' ? 'active' : ''}`}
                        onClick={() => setActiveTab('pending')}
                    >
                        Requests
                        {pending.length > 0 && <span className="tab-count pending">{pending.length}</span>}
                    </button>
                </div>

                {loading ? (
                    <div className="friends-loading">
                        {[1, 2, 3].map(i => (
                            <div key={i} className="friend-skeleton card">
                                <div className="skeleton" style={{ width: 48, height: 48, borderRadius: '50%' }}></div>
                                <div className="skeleton" style={{ height: 20, width: 150 }}></div>
                            </div>
                        ))}
                    </div>
                ) : (
                    <AnimatePresence mode="wait">
                        {activeTab === 'friends' ? (
                            <motion.div
                                key="friends"
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                exit={{ opacity: 0, x: 20 }}
                            >
                                {friends.length === 0 ? (
                                    <div className="friends-empty card">
                                        <div className="friends-empty-icon">
                                            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                                                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                                                <circle cx="9" cy="7" r="4" />
                                                <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                                                <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                                            </svg>
                                        </div>
                                        <h3>No friends yet</h3>
                                        <p>Search for users and send friend requests</p>
                                    </div>
                                ) : (
                                    <div className="friends-list">
                                        {friends.map((friend, index) => (
                                            <motion.div
                                                key={friend.id}
                                                className="friend-card card"
                                                initial={{ opacity: 0, y: 20 }}
                                                animate={{ opacity: 1, y: 0 }}
                                                transition={{ delay: index * 0.05 }}
                                            >
                                                <Link to={`/profile/${friend.id}`} className="friend-info">
                                                    <div className="avatar">
                                                        {friend.avatar_url ? (
                                                            <img src={friend.avatar_url} alt="" />
                                                        ) : (
                                                            friend.username?.charAt(0).toUpperCase()
                                                        )}
                                                    </div>
                                                    <div className="friend-details">
                                                        <span className="friend-name">{friend.full_name}</span>
                                                        <span className="friend-username">@{friend.username}</span>
                                                    </div>
                                                </Link>
                                                <button
                                                    className="btn btn-secondary btn-sm"
                                                    onClick={() => handleMessage(friend)}
                                                    disabled={startingConversationWith === friend.id}
                                                >
                                                    {startingConversationWith === friend.id ? 'Opening...' : 'Message'}
                                                </button>
                                            </motion.div>
                                        ))}
                                    </div>
                                )}
                            </motion.div>
                        ) : (
                            <motion.div
                                key="pending"
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                exit={{ opacity: 0, x: 20 }}
                            >
                                {pending.length === 0 ? (
                                    <div className="friends-empty card">
                                        <h3>No pending requests</h3>
                                        <p>Friend requests will appear here</p>
                                    </div>
                                ) : (
                                    <div className="friends-list">
                                        {pending.map((request, index) => (
                                            <motion.div
                                                key={request.id}
                                                className="friend-card card"
                                                initial={{ opacity: 0, y: 20 }}
                                                animate={{ opacity: 1, y: 0 }}
                                                transition={{ delay: index * 0.05 }}
                                            >
                                                <div className="friend-info">
                                                    <div className="avatar">
                                                        {request.requester?.avatar_url ? (
                                                            <img src={request.requester.avatar_url} alt="" />
                                                        ) : (
                                                            request.requester?.username?.charAt(0).toUpperCase() || 'U'
                                                        )}
                                                    </div>
                                                    <div className="friend-details">
                                                        <span className="friend-name">{request.requester?.full_name || 'User'}</span>
                                                        <span className="friend-username">@{request.requester?.username}</span>
                                                    </div>
                                                </div>
                                                <div className="request-actions">
                                                    <motion.button
                                                        className="btn btn-primary btn-sm"
                                                        onClick={() => handleAccept(request.id)}
                                                        whileHover={{ scale: 1.02 }}
                                                        whileTap={{ scale: 0.98 }}
                                                    >
                                                        Accept
                                                    </motion.button>
                                                    <button
                                                        className="btn btn-ghost btn-sm"
                                                        onClick={() => handleBlock(request.id)}
                                                    >
                                                        Decline
                                                    </button>
                                                </div>
                                            </motion.div>
                                        ))}
                                    </div>
                                )}
                            </motion.div>
                        )}
                    </AnimatePresence>
                )}
            </motion.div>
        </div>
    )
}
