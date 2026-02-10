import axios from 'axios'

const API_URL = '/api'

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const isAuthEndpoint = error.config?.url?.includes('/login') ||
      error.config?.url?.includes('/register') ||
      error.config?.url?.includes('/auth/')
    if (error.response?.status === 401 && !isAuthEndpoint) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authAPI = {
  register: (data) => api.post('/register', data),
  login: (data) => api.post('/login', data),
  googleLogin: (idToken) => api.post('/auth/google', { id_token: idToken }),
  getPasswordRequirements: () => api.get('/auth/password-requirements'),
}

export const usersAPI = {
  getProfile: (id) => api.get(`/users/${id}`),
  updateProfile: (id, data) => api.put('/profile', data),
  searchUsers: (query) => api.get(`/users/search?q=${query}`),
  search: (query) => api.get(`/users/search?q=${query}`),
  deleteAccount: () => api.delete('/profile'),
  updatePrivacySettings: (data) => api.put('/profile/privacy', data),
  setEmojiAvatar: (emojiId) => api.put('/profile/emoji', { emoji_id: emojiId }),
  updateOnlineStatus: (isOnline) => api.put('/profile/status', { is_online: isOnline }),
}

export const emojiAPI = {
  getEmojis: () => api.get('/emojis'),
  getUsersWithEmoji: (emojiId) => api.get(`/emojis/${emojiId}/users`),
}

export const postsAPI = {
  getFeed: () => api.get('/posts'),
  getPost: (id) => api.get(`/posts/${id}`),
  create: (data) => api.post('/posts', data),
  createPost: (data) => api.post('/posts', data),
  updatePost: (id, data) => api.put(`/posts/${id}`, data),
  deletePost: (id) => api.delete(`/posts/${id}`),
  delete: (id) => api.delete(`/posts/${id}`),
}

export const socialAPI = {
  getFriends: () => api.get('/friends'),
  getPendingRequests: () => api.get('/friends/requests'),
  sendFriendRequest: (addresseeId) => api.post('/friends/', { addressee_id: addresseeId }),
  acceptFriendRequest: (requestId) => api.put(`/friends/${requestId}`, {}),
  likePost: (postId) => api.post(`/likes/${postId}`, {}),
  unlikePost: (postId) => api.delete(`/likes/${postId}`),
  getComments: (postId) => api.get(`/comments/${postId}`),
  addComment: (postId, data) => api.post(`/comments/${postId}`, data),
}

export const messagesAPI = {
  getConversations: () => api.get('/conversations'),
  startConversation: (participantId) => api.post('/conversations', { participant_id: participantId }),
  createConversation: (participantId) => api.post('/conversations', { participant_id: participantId }),
  getMessages: (conversationId) => api.get(`/conversations/${conversationId}/messages`),
  sendMessage: (conversationId, body) => api.post(`/conversations/${conversationId}/messages`, { body }),
}

export const groupsAPI = {
  getMyGroups: () => api.get('/groups'),
  getGroup: (id) => api.get(`/groups/${id}`),
  createGroup: (data) => api.post('/groups', data),
  joinGroup: (id) => api.post(`/groups/${id}/join`, {}),
  leaveGroup: (id) => api.post(`/groups/${id}/leave`, {}),
  getGroupPosts: (id) => api.get(`/groups/${id}/posts`),
  postToGroup: (id, data) => api.post(`/groups/${id}/posts`, data),
  getGroupMembers: (id) => api.get(`/groups/${id}/members`),
  updateGroupSettings: (id, data) => api.put(`/groups/${id}/settings`, data),
}

export const notificationsAPI = {
  getNotifications: () => api.get('/notifications'),
  getList: () => api.get('/notifications'),
  markAsRead: (id) => api.put(`/notifications/${id}/read`, {}),
  clearAll: () => api.delete('/notifications/clear'),
}

export const reportsAPI = {
  createReport: (data) => api.post('/reports', data),
}

export const adminAPI = {
  getReports: (status = 'pending') => api.get(`/reports?status=${status}`),
  reviewReport: (id, status) => api.put(`/reports/${id}`, { status }),
  deleteContent: (targetType, targetId) => api.delete(`/admin/delete/${targetType}/${targetId}`),
  getStats: () => api.get('/admin/stats'),
  grantAdmin: (userId) => api.post('/admin/grant', { user_id: userId }),
  broadcast: (message) => api.post('/admin/broadcast', { message }),
  getBroadcasts: () => api.get('/admin/broadcast'),
  broadcastToEmoji: (emojiId, message) => api.post(`/admin/broadcast/emoji/${emojiId}`, { message }),
}

export const uploadAPI = {
  uploadFile: (file) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export const friendsAPI = {
  sendRequest: (addresseeId) => api.post('/friends/', { addressee_id: addresseeId }),
  getFriends: () => api.get('/friends'),
  getList: () => api.get('/friends'),
  getPendingRequests: () => api.get('/friends/requests'),
  getPending: () => api.get('/friends/requests'),
  acceptRequest: (requestId) => api.put(`/friends/${requestId}`, {}),
  accept: (requestId) => api.put(`/friends/${requestId}`, {}),
  block: (requestId) => api.delete(`/friends/${requestId}`),
}

export default api
