import { useState, useEffect, useRef } from 'react'
import { useLocation } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { messagesAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'
import './Messages.css'

export default function Messages() {
    const { user } = useAuth()
    const location = useLocation()
    const [conversations, setConversations] = useState([])
    const [activeConversation, setActiveConversation] = useState(null)
    const [messages, setMessages] = useState([])
    const [newMessage, setNewMessage] = useState('')
    const [loading, setLoading] = useState(true)
    const messagesEndRef = useRef(null)

    useEffect(() => {
        loadConversations()
    }, [])

    useEffect(() => {
        if (activeConversation) {
            loadMessages(activeConversation.id)
        }
    }, [activeConversation])

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }, [messages])

    useEffect(() => {
        if (loading) return
        const conversationId = location.state?.conversationId
        if (!conversationId) return
        if (activeConversation?.id === conversationId) return

        const existing = conversations.find(c => c.id === conversationId)
        if (existing) {
            setActiveConversation(existing)
        } else {
            setActiveConversation({
                id: conversationId,
                participant: location.state?.participant || null
            })
        }
    }, [loading, conversations, location.state, activeConversation?.id])

    const loadConversations = async () => {
        try {
            const res = await messagesAPI.getConversations()
            setConversations(res.data || [])
        } catch (err) {
            console.error('Failed to load conversations')
        } finally {
            setLoading(false)
        }
    }

    const loadMessages = async (convId) => {
        try {
            const res = await messagesAPI.getMessages(convId)
            setMessages(res.data || [])
        } catch (err) {
            console.error('Failed to load messages')
        }
    }

    const handleSend = async (e) => {
        e.preventDefault()
        if (!newMessage.trim() || !activeConversation) return

        try {
            const res = await messagesAPI.sendMessage(activeConversation.id, newMessage)
            setMessages([...messages, res.data])
            setNewMessage('')
        } catch (err) {
            console.error('Failed to send message')
        }
    }

    const formatTime = (dateStr) => {
        const date = new Date(dateStr)
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }

    return (
        <div className="messages-page">
            <motion.div
                className="messages-container"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
            >
                <div className="conversations-sidebar">
                    <div className="conversations-header">
                        <h2>Messages</h2>
                    </div>

                    <div className="conversations-list">
                        {loading ? (
                            <div className="conversations-loading">Loading...</div>
                        ) : conversations.length === 0 ? (
                            <div className="conversations-empty">
                                <p>No conversations yet</p>
                                <span>Start a conversation with a friend</span>
                            </div>
                        ) : (
                            conversations.map(conv => (
                                <motion.div
                                    key={conv.id}
                                    className={`conversation-item ${activeConversation?.id === conv.id ? 'active' : ''}`}
                                    onClick={() => setActiveConversation(conv)}
                                    whileHover={{ x: 4 }}
                                >
                                    <div className="avatar">
                                        {conv.participant?.emoji_avatar ? (
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
                                                const emoji = EMOJIS.find(e => e.id === conv.participant.emoji_avatar)
                                                return <span style={{ fontSize: '28px' }}>{emoji?.char || conv.participant.emoji_avatar}</span>
                                            })()
                                        ) : conv.participant?.avatar_url ? (
                                            <img src={conv.participant.avatar_url} alt="" />
                                        ) : (
                                            conv.participant?.username?.charAt(0).toUpperCase() || 'U'
                                        )}
                                    </div>
                                    <div className="conversation-info">
                                        <span className="conversation-name">
                                            {conv.participant?.full_name || conv.participant?.username || 'User'}
                                        </span>
                                        <span className="conversation-preview">
                                            {conv.last_message?.body || 'No messages'}
                                        </span>
                                    </div>
                                </motion.div>
                            ))
                        )}
                    </div>
                </div>

                <div className="chat-area">
                    {!activeConversation ? (
                        <div className="chat-empty">
                            <div className="chat-empty-icon">
                                <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                                </svg>
                            </div>
                            <h3>Select a conversation</h3>
                            <p>Choose a conversation from the sidebar to start messaging</p>
                        </div>
                    ) : (
                        <>
                            <div className="chat-header">
                                <div className="avatar">
                                    {activeConversation.participant?.emoji_avatar ? (
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
                                            const emoji = EMOJIS.find(e => e.id === activeConversation.participant.emoji_avatar)
                                            return <span style={{ fontSize: '28px' }}>{emoji?.char || activeConversation.participant.emoji_avatar}</span>
                                        })()
                                    ) : activeConversation.participant?.avatar_url ? (
                                        <img src={activeConversation.participant.avatar_url} alt="" />
                                    ) : (
                                        activeConversation.participant?.username?.charAt(0).toUpperCase() || 'U'
                                    )}
                                </div>
                                <div className="chat-header-info">
                                    <h3>{activeConversation.participant?.full_name || activeConversation.participant?.username || 'User'}</h3>
                                    <span>Online</span>
                                </div>
                            </div>

                            <div className="chat-messages">
                                <AnimatePresence>
                                    {messages.map((msg, idx) => (
                                        <motion.div
                                            key={msg.id}
                                            className={`message ${msg.user_id === user?.id ? 'own' : ''}`}
                                            initial={{ opacity: 0, y: 10 }}
                                            animate={{ opacity: 1, y: 0 }}
                                            transition={{ delay: idx * 0.02 }}
                                        >
                                            <div className="message-bubble">
                                                <p>{msg.body}</p>
                                                <span className="message-time">{formatTime(msg.created_at)}</span>
                                            </div>
                                        </motion.div>
                                    ))}
                                </AnimatePresence>
                                <div ref={messagesEndRef} />
                            </div>

                            <form className="chat-input" onSubmit={handleSend}>
                                <input
                                    type="text"
                                    className="input-field"
                                    placeholder="Type a message..."
                                    value={newMessage}
                                    onChange={e => setNewMessage(e.target.value)}
                                />
                                <motion.button
                                    type="submit"
                                    className="btn btn-primary"
                                    disabled={!newMessage.trim()}
                                    whileHover={{ scale: 1.02 }}
                                    whileTap={{ scale: 0.98 }}
                                >
                                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                        <line x1="22" y1="2" x2="11" y2="13" />
                                        <polygon points="22 2 15 22 11 13 2 9 22 2" />
                                    </svg>
                                </motion.button>
                            </form>
                        </>
                    )}
                </div>
            </motion.div>
        </div>
    )
}
