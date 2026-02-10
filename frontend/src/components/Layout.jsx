import { Outlet } from 'react-router-dom'
import Navbar from './Navbar'
import Sidebar from './Sidebar'
import './Layout.css'

export default function Layout() {
    return (
        <div className="app">
            <Navbar />
            <Sidebar />
            <main className="main-content">
                <Outlet />
            </main>
        </div>
    )
}
