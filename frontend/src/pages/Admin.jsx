import { useState, useEffect } from 'react'
import { adminAPI, usersAPI, emojiAPI, groupsAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'
import './Admin.css'

const NAV_ITEMS = [
    { id: 'stats', icon: '📊', label: 'Statistics' },
    { id: 'reports', icon: '📋', label: 'Reports' },
    { id: 'emoji', icon: '😊', label: 'Emoji' },
    { id: 'broadcast', icon: '📢', label: 'Broadcast' },
    { id: 'admins', icon: '🛡️', label: 'Admins' },
    { id: 'groups', icon: '👥', label: 'Groups' },
]

export default function Admin() {
    const { user } = useAuth()
    const [activeTab, setActiveTab] = useState('stats')
    const [reports, setReports] = useState([])
    const [filterStatus, setFilterStatus] = useState('pending')
    const [stats, setStats] = useState(null)
    const [broadcasts, setBroadcasts] = useState([])
    const [broadcastMessage, setBroadcastMessage] = useState('')
    const [searchQuery, setSearchQuery] = useState('')
    const [searchResults, setSearchResults] = useState([])
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState('')
    const [emojiStats, setEmojiStats] = useState([])
    const [showEmojiBroadcast, setShowEmojiBroadcast] = useState(false)
    const [selectedEmoji, setSelectedEmoji] = useState(null)
    const [emojiBroadcastMessage, setEmojiBroadcastMessage] = useState('')
    const [expandedEmojis, setExpandedEmojis] = useState({})
    const [allGroups, setAllGroups] = useState([])
    const [newGroup, setNewGroup] = useState({ title: '', description: '' })
    const [creatingGroup, setCreatingGroup] = useState(false)

    useEffect(() => {
        if (!user?.is_admin) return
        loadData()
    }, [user, filterStatus, activeTab])

    useEffect(() => {
        if (message) {
            const t = setTimeout(() => setMessage(''), 3500)
            return () => clearTimeout(t)
        }
    }, [message])

    const loadData = async () => {
        setLoading(true)
        try {
            switch (activeTab) {
                case 'stats':
                    const statsRes = await adminAPI.getStats()
                    setStats(statsRes.data)
                    break
                case 'reports':
                    const reportsRes = await adminAPI.getReports(filterStatus)
                    setReports(reportsRes.data || [])
                    break
                case 'broadcast':
                    const broadcastRes = await adminAPI.getBroadcasts()
                    setBroadcasts(broadcastRes.data || [])
                    break
                case 'emoji':
                    await loadEmojiStats()
                    break
                case 'groups':
                    await loadAllGroups()
                    break
            }
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to load data')
        }
        setLoading(false)
    }

    const loadEmojiStats = async () => {
        const EMOJIS = [
            { id: 'smile', char: '😊', name: 'Smile' },
            { id: 'cool', char: '😎', name: 'Cool' },
            { id: 'fire', char: '🔥', name: 'Fire' },
            { id: 'lightning', char: '⚡', name: 'Lightning' },
            { id: 'gaming', char: '🎮', name: 'Gaming' },
            { id: 'art', char: '🎨', name: 'Art' },
            { id: 'music', char: '🎵', name: 'Music' },
            { id: 'book', char: '📚', name: 'Book' },
            { id: 'star', char: '🌟', name: 'Star' },
            { id: 'strong', char: '💪', name: 'Strong' },
        ]
        const results = await Promise.all(
            EMOJIS.map(async (emoji) => {
                try {
                    const res = await emojiAPI.getUsersWithEmoji(emoji.id)
                    return { ...emoji, users: res.data || [], count: (res.data || []).length }
                } catch {
                    return { ...emoji, users: [], count: 0 }
                }
            })
        )
        setEmojiStats(results.sort((a, b) => b.count - a.count))
    }

    const loadAllGroups = async () => {
        try {
            const res = await groupsAPI.getMyGroups()
            setAllGroups(res.data || [])
        } catch {
            setAllGroups([])
        }
    }

    const handleCreateGroup = async (e) => {
        e.preventDefault()
        if (!newGroup.title.trim()) return
        setCreatingGroup(true)
        try {
            await groupsAPI.createGroup(newGroup)
            setNewGroup({ title: '', description: '' })
            setMessage('Group created!')
            await loadAllGroups()
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to create group')
        }
        setCreatingGroup(false)
    }

    const handleEmojiBroadcast = async (e) => {
        e.preventDefault()
        if (!emojiBroadcastMessage.trim() || !selectedEmoji) return
        setLoading(true)
        try {
            await adminAPI.broadcastToEmoji(selectedEmoji.id, emojiBroadcastMessage)
            setMessage(`Broadcast sent to ${selectedEmoji.count} users with ${selectedEmoji.char}!`)
            setShowEmojiBroadcast(false)
            setEmojiBroadcastMessage('')
            setSelectedEmoji(null)
        } catch (err) {
            setMessage(err.response?.data?.error || 'Broadcast failed')
        }
        setLoading(false)
    }

    const handleReview = async (reportId, status) => {
        try {
            await adminAPI.reviewReport(reportId, status)
            setReports(reports.filter((r) => r.id !== reportId))
            setMessage('Report updated')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Update failed')
        }
    }

    const handleDelete = async (report) => {
        try {
            await adminAPI.deleteContent(report.target_type, report.target_id)
            await adminAPI.reviewReport(report.id, 'resolved')
            setReports(reports.filter((r) => r.id !== report.id))
            setMessage('Content deleted')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Delete failed')
        }
    }

    const handleBroadcast = async (e) => {
        e.preventDefault()
        if (!broadcastMessage.trim()) return
        setLoading(true)
        try {
            await adminAPI.broadcast(broadcastMessage)
            setBroadcastMessage('')
            setMessage('Broadcast sent to all users!')
            const res = await adminAPI.getBroadcasts()
            setBroadcasts(res.data || [])
        } catch (err) {
            setMessage(err.response?.data?.error || 'Broadcast failed')
        }
        setLoading(false)
    }

    const handleSearch = async (e) => {
        e.preventDefault()
        if (!searchQuery.trim()) return
        try {
            const res = await usersAPI.searchUsers(searchQuery)
            setSearchResults(res.data || [])
        } catch {
            setMessage('Search failed')
        }
    }

    const handleGrantAdmin = async (userId) => {
        try {
            await adminAPI.grantAdmin(userId)
            setSearchResults(searchResults.map(u =>
                u.id === userId ? { ...u, is_admin: true } : u
            ))
            setMessage('Admin rights granted!')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Error')
        }
    }

    if (!user?.is_admin) {
        return (
            <div className="admin-page">
                <div className="admin-denied">
                    <h1>🚫 Access Denied</h1>
                    <p>You do not have permission to access this page.</p>
                </div>
            </div>
        )
    }

    const tabMeta = NAV_ITEMS.find(n => n.id === activeTab)

    return (
        <div className="admin-page">
            <aside className="admin-sidebar">
                <div className="admin-sidebar-title">Admin Panel</div>
                {NAV_ITEMS.map(item => (
                    <button
                        key={item.id}
                        className={`admin-nav-item ${activeTab === item.id ? 'active' : ''}`}
                        onClick={() => setActiveTab(item.id)}
                    >
                        <span className="admin-nav-icon">{item.icon}</span>
                        {item.label}
                    </button>
                ))}
            </aside>

            <main className="admin-content">
                {message && <div className="admin-toast">{message}</div>}

                <div className="admin-content-header">
                    <h2>{tabMeta?.icon} {tabMeta?.label}</h2>
                </div>

                {activeTab === 'stats' && stats && (
                    <div className="admin-stats-grid">
                        <div className="admin-stat-card">
                            <h3>👤 Users</h3>
                            <span className="stat-value">{stats.total_users}</span>
                        </div>
                        <div className="admin-stat-card">
                            <h3>🟢 Online</h3>
                            <span className="stat-value">{stats.online_users}</span>
                        </div>
                        <div className="admin-stat-card">
                            <h3>📝 Posts</h3>
                            <span className="stat-value">{stats.total_posts}</span>
                        </div>
                        <div className="admin-stat-card">
                            <h3>👥 Groups</h3>
                            <span className="stat-value">{stats.total_groups}</span>
                        </div>
                        <div className="admin-stat-card">
                            <h3>💬 Messages</h3>
                            <span className="stat-value">{stats.total_messages}</span>
                        </div>
                        <div className="admin-stat-card warning">
                            <h3>⚠️ Reports</h3>
                            <span className="stat-value">{stats.pending_reports}</span>
                        </div>
                    </div>
                )}

                {activeTab === 'reports' && (
                    <div className="admin-reports-section">
                        <div className="admin-filter-bar">
                            <select value={filterStatus} onChange={(e) => setFilterStatus(e.target.value)}>
                                <option value="pending">Pending</option>
                                <option value="reviewed">Reviewed</option>
                                <option value="resolved">Resolved</option>
                            </select>
                        </div>

                        {reports.length === 0 ? (
                            <p className="admin-empty-state">No reports with status "{filterStatus}"</p>
                        ) : (
                            reports.map((report) => (
                                <div key={report.id} className="admin-report-card">
                                    <div className="report-header">
                                        <span className="report-type">{report.target_type}</span>
                                        <span className="report-id">#{report.target_id}</span>
                                    </div>
                                    <p className="report-reason">{report.reason}</p>
                                    <p className="report-meta">
                                        Reported by: {report.reporter?.username || 'Unknown'} · {new Date(report.created_at).toLocaleString()}
                                    </p>
                                    {filterStatus === 'pending' && (
                                        <div className="report-actions">
                                            <button onClick={() => handleReview(report.id, 'reviewed')}>Mark Reviewed</button>
                                            <button onClick={() => handleReview(report.id, 'resolved')}>Dismiss</button>
                                            <button className="danger-btn" onClick={() => handleDelete(report)}>Delete Content</button>
                                        </div>
                                    )}
                                </div>
                            ))
                        )}
                    </div>
                )}

                {activeTab === 'emoji' && (
                    <div className="admin-emoji-section">
                        {loading ? (
                            <p className="admin-empty-state">Loading...</p>
                        ) : emojiStats.length === 0 ? (
                            <p className="admin-empty-state">No data</p>
                        ) : (
                            <div className="admin-emoji-grid">
                                {emojiStats.map((emoji) => (
                                    <div key={emoji.id} className="admin-emoji-card-wrapper">
                                        <div
                                            className="admin-emoji-card"
                                            onClick={() => emoji.count > 0 && setExpandedEmojis(prev => ({
                                                ...prev,
                                                [emoji.id]: !prev[emoji.id]
                                            }))}
                                            style={{ cursor: emoji.count > 0 ? 'pointer' : 'default' }}
                                        >
                                            <div className="admin-emoji-icon">{emoji.char}</div>
                                            <div className="admin-emoji-info">
                                                <h4>{emoji.name}</h4>
                                                <p>{emoji.count} {emoji.count === 1 ? 'user' : 'users'}</p>
                                            </div>
                                            <div className="admin-emoji-actions">
                                                {emoji.count > 0 && (
                                                    <>
                                                        <button
                                                            onClick={(e) => {
                                                                e.stopPropagation()
                                                                setExpandedEmojis(prev => ({
                                                                    ...prev,
                                                                    [emoji.id]: !prev[emoji.id]
                                                                }))
                                                            }}
                                                        >
                                                            {expandedEmojis[emoji.id] ? 'Hide' : 'Users'}
                                                        </button>
                                                        <button
                                                            className="btn-primary"
                                                            onClick={(e) => {
                                                                e.stopPropagation()
                                                                setSelectedEmoji(emoji)
                                                                setShowEmojiBroadcast(true)
                                                            }}
                                                        >
                                                            Broadcast
                                                        </button>
                                                    </>
                                                )}
                                            </div>
                                        </div>
                                        {expandedEmojis[emoji.id] && emoji.users?.length > 0 && (
                                            <div className="admin-emoji-users">
                                                <h5>Users with {emoji.char} {emoji.name}:</h5>
                                                <ul>
                                                    {emoji.users.map((u) => (
                                                        <li key={u.id || u.ID}>
                                                            <span className="emoji-user-name">{u.username || u.Username}</span>
                                                            {(u.email || u.Email) && (
                                                                <span className="emoji-user-email"> — {u.email || u.Email}</span>
                                                            )}
                                                        </li>
                                                    ))}
                                                </ul>
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}

                {activeTab === 'broadcast' && (
                    <div className="admin-broadcast-section">
                        <div className="admin-form-card">
                            <form onSubmit={handleBroadcast}>
                                <div className="form-group">
                                    <label>Send notification to all users:</label>
                                    <textarea
                                        value={broadcastMessage}
                                        onChange={(e) => setBroadcastMessage(e.target.value)}
                                        placeholder="Enter broadcast message..."
                                        rows={3}
                                    />
                                </div>
                                <button type="submit" className="admin-form-btn" disabled={loading || !broadcastMessage.trim()}>
                                    {loading ? 'Sending...' : '📤 Send Broadcast'}
                                </button>
                            </form>
                        </div>

                        <div className="broadcast-history">
                            <h3>Broadcast History</h3>
                            {broadcasts.length === 0 ? (
                                <p className="admin-empty-state">No broadcasts yet</p>
                            ) : (
                                <div className="broadcast-list">
                                    {broadcasts.map((b) => (
                                        <div key={b.id} className="broadcast-item">
                                            <p>{b.message}</p>
                                            <span className="broadcast-date">{new Date(b.created_at).toLocaleString()}</span>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {activeTab === 'admins' && (
                    <div className="admin-search-section">
                        <form className="admin-search-form" onSubmit={handleSearch}>
                            <input
                                type="text"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                placeholder="Search by username..."
                            />
                            <button type="submit">🔍 Search</button>
                        </form>

                        {searchResults.length > 0 && (
                            <div className="admin-user-results">
                                {searchResults.map((u) => (
                                    <div key={u.id} className="admin-user-row">
                                        <div className="admin-user-info">
                                            {u.username}
                                            <span className="user-email">({u.email})</span>
                                        </div>
                                        {u.is_admin ? (
                                            <span className="admin-badge">🛡️ Admin</span>
                                        ) : (
                                            <button className="admin-grant-btn" onClick={() => handleGrantAdmin(u.id)}>
                                                Grant Admin
                                            </button>
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}

                {activeTab === 'groups' && (
                    <div className="admin-groups-section">
                        <div className="admin-form-card">
                            <form onSubmit={handleCreateGroup}>
                                <div className="form-group">
                                    <label>Group Name</label>
                                    <input
                                        type="text"
                                        value={newGroup.title}
                                        onChange={(e) => setNewGroup({ ...newGroup, title: e.target.value })}
                                        placeholder="Enter group name..."
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label>Description</label>
                                    <textarea
                                        value={newGroup.description}
                                        onChange={(e) => setNewGroup({ ...newGroup, description: e.target.value })}
                                        placeholder="Group description (optional)..."
                                        rows={3}
                                    />
                                </div>
                                <button type="submit" className="admin-form-btn" disabled={creatingGroup || !newGroup.title.trim()}>
                                    {creatingGroup ? 'Creating...' : '➕ Create Group'}
                                </button>
                            </form>
                        </div>

                        <div className="admin-groups-list">
                            <h3>Existing Groups ({allGroups.length})</h3>
                            {allGroups.length === 0 ? (
                                <p className="admin-empty-state">No groups yet</p>
                            ) : (
                                allGroups.map((group) => (
                                    <div key={group.id} className="admin-group-item">
                                        <div>
                                            <div className="admin-group-name">{group.title}</div>
                                            {group.description && <div className="admin-group-desc">{group.description}</div>}
                                        </div>
                                        <div className="admin-group-meta">
                                            <div className="admin-group-members">{group.member_count || 0} members</div>
                                            <div className="admin-group-date">
                                                {group.created_at ? new Date(group.created_at).toLocaleDateString() : ''}
                                            </div>
                                        </div>
                                    </div>
                                ))
                            )}
                        </div>
                    </div>
                )}
            </main>

            {showEmojiBroadcast && selectedEmoji && (
                <div className="modal-overlay" onClick={() => setShowEmojiBroadcast(false)}>
                    <div className="modal" onClick={(e) => e.stopPropagation()}>
                        <h3>Broadcast to {selectedEmoji.char} Users</h3>
                        <p>Send notification to {selectedEmoji.count} users with emoji {selectedEmoji.name}</p>
                        <form onSubmit={handleEmojiBroadcast}>
                            <textarea
                                value={emojiBroadcastMessage}
                                onChange={(e) => setEmojiBroadcastMessage(e.target.value)}
                                placeholder="Enter message..."
                                rows={4}
                                required
                            />
                            <div className="modal-actions">
                                <button type="button" onClick={() => setShowEmojiBroadcast(false)}>Cancel</button>
                                <button type="submit" className="btn-primary" disabled={loading}>
                                    {loading ? 'Sending...' : 'Send'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    )
}
