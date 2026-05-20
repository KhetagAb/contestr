import { useState, type ReactNode } from "react";
import { FiFileText } from "react-icons/fi";
import { useContests } from "@/shared/hooks/useContests";
import logo from "@/assets/icons/logo.svg";
import adminUserIcon from "@/assets/images/admin-user-icon.png";
import { ContestClock } from "./ContestClock";

export type AdminSidebarSession = {
    username: string;
    onLogout: () => void;
};

const Submenu = ({
    isOpen,
    children,
}: {
    isOpen: boolean;
    children: ReactNode;
}) => (
    <div className={`submenu-wrapper${isOpen ? " submenu-wrapper--open" : ""}`}>
        {children}
    </div>
);

type SidebarProps = {
    adminSession?: AdminSidebarSession | null;
};

export function Sidebar({ adminSession }: SidebarProps) {
    const [hoveredMenu, setHoveredMenu] = useState<string | null>(null);
    const { contests, isLoading } = useContests();

    return (
        <aside className="sidebar-hover">
            <div className="icon-wrap">
                <div className="icon-box">
                    <img src={logo} alt="Логотип" className="logo-icon" />
                </div>
            </div>

            <div className="icon-wrap">
                <div
                    className="menu-item-container"
                    onMouseEnter={() => setHoveredMenu("contests")}
                    onMouseLeave={() => setHoveredMenu(null)}
                >
                    <div className="icon-box">
                        <FiFileText size={35} />
                    </div>
                    <Submenu isOpen={hoveredMenu === "contests"}>
                        {isLoading && (
                            <span className="submenu-item submenu-item-text">Загрузка…</span>
                        )}
                        {!isLoading && contests.length === 0 && (
                            <span className="submenu-item submenu-item-text">Контесты не настроены</span>
                        )}
                        {contests.map((item) => (
                            <a
                                key={item.contest_id}
                                className="submenu-item"
                                href={`/?contestId=${item.contest_id}`}
                            >
                                <span className="submenu-item-text">{item.name}</span>
                            </a>
                        ))}
                    </Submenu>
                </div>
            </div>

            <div className="time">
                <ContestClock />
            </div>

            <div className="time-text">
                <p>
                    <span>Время до </span>
                    <span>конца тура</span>
                </p>
            </div>

            {adminSession && (
                <div
                    className="menu-item-container sidebar-admin-user"
                    onMouseEnter={() => setHoveredMenu("admin-user")}
                    onMouseLeave={() => setHoveredMenu(null)}
                >
                    <div className="icon-box sidebar-admin-user-trigger" aria-label="Аккаунт администратора">
                        <img src={adminUserIcon} alt="" className="sidebar-admin-user-icon" />
                    </div>
                    <Submenu isOpen={hoveredMenu === "admin-user"}>
                        <div className="admin-user-popover">
                            <p className="admin-user-popover-greeting">
                                Вы вошли как <span className="h1_pink">{adminSession.username}</span>
                            </p>
                            <a href="/admin" className="admin-user-popover-link">
                                Админ-панель
                            </a>
                            <button
                                type="button"
                                className="admin-user-popover-btn"
                                onClick={adminSession.onLogout}
                            >
                                Выйти
                            </button>
                        </div>
                    </Submenu>
                </div>
            )}
        </aside>
    );
}
