import { useState, useEffect } from 'react'
import { useAuth } from '../context/AuthContext'
import { usersAPI, emojiAPI } from '../services/api'
import { useNavigate } from 'react-router-dom'
import { useToast } from '../components/Toast'
import './Settings.css'

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

export default function Settings() {
    const { user, updateUser, logout } = useAuth()
    const navigate = useNavigate()
    const toast = useToast()
    const [formData, setFormData] = useState({
        full_name: '',
        bio: '',
        avatar_url: '',
    })
    const [privacySettings, setPrivacySettings] = useState({
        show_last_seen: 'all',
        allow_messages_from: 'all',
    })
    const [selectedEmoji, setSelectedEmoji] = useState('')
    const [loading, setLoading] = useState(false)
    const [showDeleteModal, setShowDeleteModal] = useState(false)
    const [deleteConfirm, setDeleteConfirm] = useState('')
    const [usersWithEmoji, setUsersWithEmoji] = useState([])
    const [showEmojiUsers, setShowEmojiUsers] = useState(false)

    useEffect(() => {
        if (user) {
            setFormData({
                full_name: user.full_name || '',
                bio: user.bio || '',
                avatar_url: user.avatar_url || '',
            })
            setSelectedEmoji(user.emoji_avatar || '')

            const showLastSeen = user.show_last_seen || 'all'
            const allowMessagesFrom = user.allow_messages_from || 'all'

            setPrivacySettings({
                show_last_seen: showLastSeen,
                allow_messages_from: allowMessagesFrom,
            })
        }
    }, [user])

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value })
    }

    const handlePrivacyChange = (e) => {
        setPrivacySettings({ ...privacySettings, [e.target.name]: e.target.value })
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setLoading(true)
        try {
            await usersAPI.updateProfile(user.id, formData)
            updateUser({ ...user, ...formData })
            toast.success('Profile updated successfully!')
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to update profile')
        }
        setLoading(false)
    }

    const handlePrivacySave = async () => {
        setLoading(true)
        try {
            await usersAPI.updatePrivacySettings(privacySettings)
            updateUser({ ...user, ...privacySettings })
            toast.success('Privacy settings saved!')
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to save privacy settings')
        }
        setLoading(false)
    }

    const handleEmojiSelect = async (emojiId) => {
        setLoading(true)
        try {
            await usersAPI.setEmojiAvatar(emojiId)
            setSelectedEmoji(emojiId)
            updateUser({ ...user, emoji_avatar: emojiId })
            toast.success('Emoji avatar updated!')
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to update emoji')
        }
        setLoading(false)
    }

    const handleFindUsersWithEmoji = async (emojiId) => {
        try {
            const res = await emojiAPI.getUsersWithEmoji(emojiId)
            setUsersWithEmoji(res.data || [])
            setShowEmojiUsers(true)
        } catch (err) {
            toast.error('Failed to fetch users')
        }
    }

    const handleDeleteAccount = async () => {
        if (deleteConfirm !== 'DELETE') {
            toast.error('Please type DELETE to confirm')
            return
        }
        setLoading(true)
        try {
            await usersAPI.deleteAccount()
            logout()
            navigate('/login')
        } catch (err) {
            toast.error(err.response?.data?.error || 'Failed to delete account')
        }
        setLoading(false)
    }

    if (!user) return null

    return (
        <div className="settings-page">
            <h1>Settings</h1>

            <section className="settings-section">
                <h2>Profile Information</h2>
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Full Name</label>
                        <input
                            type="text"
                            name="full_name"
                            value={formData.full_name}
                            onChange={handleChange}
                        />
                    </div>
                    <div className="form-group">
                        <label>Bio</label>
                        <textarea
                            name="bio"
                            value={formData.bio}
                            onChange={handleChange}
                            rows={4}
                        />
                    </div>
                    <div className="form-group">
                        <label>Avatar URL</label>
                        <input
                            type="text"
                            name="avatar_url"
                            value={formData.avatar_url}
                            onChange={handleChange}
                        />
                    </div>
                    <button type="submit" disabled={loading}>
                        {loading ? 'Saving...' : 'Save Profile'}
                    </button>
                </form>
            </section>

            <section className="settings-section">
                <h2>Emoji Avatar</h2>
                <p>Select an emoji to represent you across the app:</p>
                <div className="emoji-grid">
                    {EMOJIS.map((emoji) => (
                        <button
                            key={emoji.id}
                            className={`emoji-btn ${selectedEmoji === emoji.id ? 'selected' : ''}`}
                            onClick={() => handleEmojiSelect(emoji.id)}
                            title={emoji.name}
                        >
                            <span className="emoji">{emoji.char}</span>
                        </button>
                    ))}
                </div>
                {selectedEmoji && (
                    <button
                        className="link-btn"
                        onClick={() => handleFindUsersWithEmoji(selectedEmoji)}
                    >
                        Find users with same emoji
                    </button>
                )}
                {showEmojiUsers && usersWithEmoji.length > 0 && (
                    <div className="emoji-users">
                        <h4>Users with same emoji:</h4>
                        <ul>
                            {usersWithEmoji.map((u) => (
                                <li key={u.user_id}>
                                    <a href={`/users/${u.user_id}`}>{u.username}</a>
                                </li>
                            ))}
                        </ul>
                    </div>
                )}
            </section>

            <section className="settings-section">
                <h2>Privacy Settings</h2>
                <div className="form-group">
                    <label>Show Last Seen To</label>
                    <select
                        name="show_last_seen"
                        value={privacySettings.show_last_seen}
                        onChange={handlePrivacyChange}
                    >
                        <option value="all">Everyone</option>
                        <option value="friends">Friends Only</option>
                        <option value="nobody">Nobody</option>
                    </select>
                </div>
                <div className="form-group">
                    <label>Allow Messages From</label>
                    <select
                        name="allow_messages_from"
                        value={privacySettings.allow_messages_from}
                        onChange={handlePrivacyChange}
                    >
                        <option value="all">Everyone</option>
                        <option value="friends">Friends Only</option>
                    </select>
                </div>
                <button onClick={handlePrivacySave} disabled={loading}>
                    Save Privacy Settings
                </button>
            </section>

            <section className="settings-section account-info">
                <h2>Account</h2>
                <p>Email: {user.email}</p>
                <p>Username: @{user.username}</p>
                <p>Member since: {new Date(user.created_at).toLocaleDateString()}</p>
            </section>

            <section className="settings-section danger-zone">
                <h2>Danger Zone</h2>
                <p>Deleting your account is permanent and cannot be undone.</p>
                <button
                    className="danger-btn"
                    onClick={() => setShowDeleteModal(true)}
                >
                    Delete Account
                </button>
            </section>

            {showDeleteModal && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h3>Delete Account</h3>
                        <p>This action cannot be undone. All your data will be permanently deleted.</p>
                        <p>Type <strong>DELETE</strong> to confirm:</p>
                        <input
                            type="text"
                            value={deleteConfirm}
                            onChange={(e) => setDeleteConfirm(e.target.value)}
                            placeholder="Type DELETE"
                        />
                        <div className="modal-actions">
                            <button onClick={() => setShowDeleteModal(false)}>Cancel</button>
                            <button className="danger-btn" onClick={handleDeleteAccount} disabled={loading}>
                                {loading ? 'Deleting...' : 'Delete Account'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}
