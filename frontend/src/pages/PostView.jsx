import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { postsAPI } from '../services/api'
import PostCard from '../components/PostCard'
import './Feed.css'

export default function PostView() {
    const { id } = useParams()
    const navigate = useNavigate()
    const [post, setPost] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        loadPost()
    }, [id])

    const loadPost = async () => {
        setLoading(true)
        setError(null)
        try {
            const res = await postsAPI.getPost(id)
            setPost(res.data)
        } catch (err) {
            console.error('Failed to load post:', err)
            setError('Post not found or you do not have permission to view it')
        } finally {
            setLoading(false)
        }
    }

    const handleDelete = () => {
        navigate('/')
    }

    if (loading) {
        return (
            <div className="feed-page">
                <div className="loading-spinner">Loading post...</div>
            </div>
        )
    }

    if (error || !post) {
        return (
            <div className="feed-page">
                <div className="error-message">
                    <h2>{error || 'Post not found'}</h2>
                    <button className="btn btn-primary" onClick={() => navigate('/')}>
                        Go to Feed
                    </button>
                </div>
            </div>
        )
    }

    return (
        <div className="feed-page">
            <h1>Post</h1>
            <PostCard
                post={post}
                onUpdate={setPost}
                onDelete={handleDelete}
            />
        </div>
    )
}
