import { createContext, useContext, useState, useEffect } from 'react'
import { authAPI, usersAPI } from '../services/api'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const storedUser = localStorage.getItem('user')
        const token = localStorage.getItem('token')

        if (storedUser && token) {
            setUser(JSON.parse(storedUser))
            usersAPI.updateOnlineStatus(true).catch(() => { })
        }
        setLoading(false)
    }, [])

    useEffect(() => {
        const handleBeforeUnload = () => {
            if (user) {
                usersAPI.updateOnlineStatus(false).catch(() => { })
            }
        }
        window.addEventListener('beforeunload', handleBeforeUnload)
        return () => window.removeEventListener('beforeunload', handleBeforeUnload)
    }, [user])

    const login = async (email, password) => {
        const response = await authAPI.login({ email, password })
        const { token, user: userData } = response.data

        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(userData))
        setUser(userData)
        usersAPI.updateOnlineStatus(true).catch(() => { })

        return userData
    }

    const googleLogin = async (idToken) => {
        const response = await authAPI.googleLogin(idToken)
        const { token, user: userData } = response.data

        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(userData))
        setUser(userData)
        usersAPI.updateOnlineStatus(true).catch(() => { })

        return userData
    }

    const register = async (data) => {
        const response = await authAPI.register(data)
        return response.data
    }

    const logout = () => {
        usersAPI.updateOnlineStatus(false).catch(() => { })
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        setUser(null)
    }

    const deleteAccount = async () => {
        await usersAPI.deleteAccount()
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        setUser(null)
    }

    const updateUser = (updates) => {
        const updatedUser = { ...user, ...updates }
        localStorage.setItem('user', JSON.stringify(updatedUser))
        setUser(updatedUser)
    }

    return (
        <AuthContext.Provider value={{ user, loading, login, googleLogin, register, logout, deleteAccount, updateUser }}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext)
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
