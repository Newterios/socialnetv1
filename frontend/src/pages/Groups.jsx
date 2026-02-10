import { useState, useEffect } from 'react'
import { groupsAPI, uploadAPI } from '../services/api'
import { useAuth } from '../context/AuthContext'

export default function Groups() {
    const { user } = useAuth()
    const [groups, setGroups] = useState([])
    const [selectedGroup, setSelectedGroup] = useState(null)
    const [groupPosts, setGroupPosts] = useState([])
    const [groupMembers, setGroupMembers] = useState([])
    const [showCreateModal, setShowCreateModal] = useState(false)
    const [showMembersModal, setShowMembersModal] = useState(false)
    const [showSettingsModal, setShowSettingsModal] = useState(false)
    const [newGroup, setNewGroup] = useState({ title: '', description: '' })
    const [newPost, setNewPost] = useState({ content: '', media_url: '' })
    const [groupSettings, setGroupSettings] = useState({ title: '', description: '' })
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState('')
    const [uploadingFile, setUploadingFile] = useState(false)

    useEffect(() => {
        loadGroups()
    }, [])

    const loadGroups = async () => {
        try {
            const res = await groupsAPI.getMyGroups()
            setGroups(res.data || [])
        } catch (err) {
            setMessage('Failed to load groups')
        }
    }

    const loadGroupDetails = async (groupId) => {
        try {
            const [groupRes, postsRes] = await Promise.all([
                groupsAPI.getGroup(groupId),
                groupsAPI.getGroupPosts(groupId),
            ])
            setSelectedGroup(groupRes.data)
            setGroupPosts(postsRes.data || [])
        } catch (err) {
            setMessage('Failed to load group details')
        }
    }

    const loadMembers = async (groupId) => {
        try {
            const res = await groupsAPI.getGroupMembers(groupId)
            setGroupMembers(res.data || [])
            setShowMembersModal(true)
        } catch (err) {
            setMessage('Failed to load members')
        }
    }

    const handleCreateGroup = async (e) => {
        e.preventDefault()
        if (!newGroup.title.trim()) return
        setLoading(true)
        try {
            const res = await groupsAPI.createGroup(newGroup)
            setGroups([res.data, ...groups])
            setNewGroup({ title: '', description: '' })
            setShowCreateModal(false)
            setMessage('Group created!')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to create group')
        }
        setLoading(false)
    }

    const handleJoinGroup = async (groupId) => {
        try {
            await groupsAPI.joinGroup(groupId)
            loadGroups()
            setMessage('Joined group!')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to join group')
        }
    }

    const handleLeaveGroup = async (groupId) => {
        try {
            await groupsAPI.leaveGroup(groupId)
            loadGroups()
            setSelectedGroup(null)
            setMessage('Left group')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to leave group')
        }
    }

    const handleFileUpload = async (e) => {
        const file = e.target.files[0]
        if (!file) return
        setUploadingFile(true)
        try {
            const res = await uploadAPI.uploadFile(file)
            setNewPost({ ...newPost, media_url: res.data.url })
            setMessage('File uploaded!')
        } catch (err) {
            setMessage('Failed to upload file')
        }
        setUploadingFile(false)
    }

    const handlePostToGroup = async (e) => {
        e.preventDefault()
        if (!newPost.content.trim() && !newPost.media_url) return
        setLoading(true)
        try {
            await groupsAPI.postToGroup(selectedGroup.id, newPost)
            setNewPost({ content: '', media_url: '' })
            loadGroupDetails(selectedGroup.id)
            setMessage('Post created!')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to post')
        }
        setLoading(false)
    }

    const handleUpdateSettings = async (e) => {
        e.preventDefault()
        setLoading(true)
        try {
            await groupsAPI.updateGroupSettings(selectedGroup.id, groupSettings)
            setSelectedGroup({ ...selectedGroup, ...groupSettings })
            setShowSettingsModal(false)
            loadGroups()
            setMessage('Settings updated!')
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to update settings')
        }
        setLoading(false)
    }

    const openSettings = () => {
        setGroupSettings({
            title: selectedGroup.title,
            description: selectedGroup.description || '',
        })
        setShowSettingsModal(true)
    }

    const isOwner = selectedGroup?.owner_id === user?.id

    return (
        <div className="groups-page">
            <div className="groups-header">
                <h1>Groups</h1>
                <button onClick={() => setShowCreateModal(true)}>Create Group</button>
            </div>

            {message && <div className="message">{message}</div>}

            <div className="groups-layout">
                <div className="groups-sidebar">
                    <h3>My Groups</h3>
                    {groups.length === 0 ? (
                        <p>You haven't joined any groups yet</p>
                    ) : (
                        <ul className="groups-list">
                            {groups.map((group) => (
                                <li
                                    key={group.id}
                                    className={selectedGroup?.id === group.id ? 'active' : ''}
                                    onClick={() => loadGroupDetails(group.id)}
                                >
                                    <span className="group-title">{group.title}</span>
                                    <span className="group-members">{group.member_count} members</span>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>

                <div className="group-content">
                    {selectedGroup ? (
                        <>
                            <div className="group-header">
                                <div className="group-info">
                                    <h2>{selectedGroup.title}</h2>
                                    <p>{selectedGroup.description}</p>
                                    <div className="group-meta">
                                        <span>Created by: {selectedGroup.owner?.username || 'Unknown'}</span>
                                        <span>Created: {new Date(selectedGroup.created_at).toLocaleDateString()}</span>
                                        <span>{selectedGroup.member_count} members</span>
                                    </div>
                                </div>
                                <div className="group-actions">
                                    <button onClick={() => loadMembers(selectedGroup.id)}>View Members</button>
                                    {isOwner && <button onClick={openSettings}>Settings</button>}
                                    {!isOwner && (
                                        <button className="danger-btn" onClick={() => handleLeaveGroup(selectedGroup.id)}>
                                            Leave Group
                                        </button>
                                    )}
                                </div>
                            </div>

                            <div className="post-form">
                                <form onSubmit={handlePostToGroup}>
                                    <textarea
                                        value={newPost.content}
                                        onChange={(e) => setNewPost({ ...newPost, content: e.target.value })}
                                        placeholder="Share something with the group..."
                                        rows={3}
                                    />
                                    <div className="post-form-actions">
                                        <input type="file" accept="image/*" onChange={handleFileUpload} />
                                        {newPost.media_url && (
                                            <img src={newPost.media_url} alt="Preview" className="media-preview" />
                                        )}
                                        <button type="submit" disabled={loading || uploadingFile}>
                                            {loading ? 'Posting...' : 'Post'}
                                        </button>
                                    </div>
                                </form>
                            </div>

                            <div className="group-posts">
                                {groupPosts.length === 0 ? (
                                    <p>No posts yet. Be the first to share!</p>
                                ) : (
                                    groupPosts.map((post) => (
                                        <div key={post.id} className="post-card">
                                            <div className="post-author">
                                                <span className="emoji">{post.author?.emoji_avatar || '👤'}</span>
                                                <strong>{post.author?.username || 'Unknown'}</strong>
                                                <span className="post-date">{new Date(post.created_at).toLocaleString()}</span>
                                            </div>
                                            <p className="post-content">{post.content}</p>
                                            {post.media_url && (
                                                <img src={post.media_url} alt="Post media" className="post-media" />
                                            )}
                                        </div>
                                    ))
                                )}
                            </div>
                        </>
                    ) : (
                        <div className="no-group-selected">
                            <p>Select a group to view its content</p>
                        </div>
                    )}
                </div>
            </div>

            {showCreateModal && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h3>Create New Group</h3>
                        <form onSubmit={handleCreateGroup}>
                            <div className="form-group">
                                <label>Group Name</label>
                                <input
                                    type="text"
                                    value={newGroup.title}
                                    onChange={(e) => setNewGroup({ ...newGroup, title: e.target.value })}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Description</label>
                                <textarea
                                    value={newGroup.description}
                                    onChange={(e) => setNewGroup({ ...newGroup, description: e.target.value })}
                                    rows={3}
                                />
                            </div>
                            <div className="modal-actions">
                                <button type="button" onClick={() => setShowCreateModal(false)}>Cancel</button>
                                <button type="submit" disabled={loading}>Create</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {showMembersModal && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h3>Group Members ({groupMembers.length})</h3>
                        <ul className="members-list">
                            {groupMembers.map((member) => (
                                <li key={member.id}>
                                    <span className="emoji">{member.user?.emoji_avatar || '👤'}</span>
                                    <a href={`/users/${member.user_id}`}>{member.user?.username}</a>
                                    {member.role === 'owner' && <span className="owner-badge">Owner</span>}
                                    <span className="join-date">Joined: {new Date(member.joined_at).toLocaleDateString()}</span>
                                </li>
                            ))}
                        </ul>
                        <button onClick={() => setShowMembersModal(false)}>Close</button>
                    </div>
                </div>
            )}

            {showSettingsModal && (
                <div className="modal-overlay">
                    <div className="modal">
                        <h3>Group Settings</h3>
                        <form onSubmit={handleUpdateSettings}>
                            <div className="form-group">
                                <label>Group Name</label>
                                <input
                                    type="text"
                                    value={groupSettings.title}
                                    onChange={(e) => setGroupSettings({ ...groupSettings, title: e.target.value })}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Description</label>
                                <textarea
                                    value={groupSettings.description}
                                    onChange={(e) => setGroupSettings({ ...groupSettings, description: e.target.value })}
                                    rows={3}
                                />
                            </div>
                            <div className="modal-actions">
                                <button type="button" onClick={() => setShowSettingsModal(false)}>Cancel</button>
                                <button type="submit" disabled={loading}>Save</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    )
}
