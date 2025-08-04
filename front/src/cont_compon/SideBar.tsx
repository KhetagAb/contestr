import { Home, User, Settings, LogOut } from 'lucide-react'

export const Sidebar = () => {
    return (
        <aside className="sidebar-hover">
            <div className="icon-wrap">
                <div className="icon-box">
                    <img src="/logo.svg" alt="Логотип" className="logo-icon" />
                    <span className="tooltip">Главная</span>
                </div>
            </div>

            <div className="icon-wrap">
                <div className="icon-box">
                    <Home size={24} />
                    <span className="tooltip">Главная</span>
                </div>
                <div className="icon-box">
                    <User size={24} />
                    <span className="tooltip">Профиль</span>
                </div>
                <div className="icon-box">
                    <Settings size={24} />
                    <span className="tooltip">Настройки</span>
                </div>
                <div className="icon-box">
                    <LogOut size={24} />
                    <span className="tooltip">Выход</span>
                </div>
            </div>
        </aside>
    )
}
