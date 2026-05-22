import { useState, type ReactNode } from "react";
import { FiMenu } from "react-icons/fi";
import { useFollowedParticipant } from "@/features/contest/follow/FollowedParticipantContext";
import { useContests } from "@/shared/hooks/useContests";
import logo from "@/assets/icons/logo.svg";
import adminUserIcon from "@/assets/images/admin-user-icon.png";
import { ParticipantFollowMenu } from "./ParticipantFollowMenu";
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
    const { isParticipantPickerOpen, closeParticipantPicker } = useFollowedParticipant();
    const { contests, isLoading } = useContests();

    const isUserMenuOpen = hoveredMenu === "user" || isParticipantPickerOpen;

    const closeUserMenu = () => {
        setHoveredMenu(null);
        closeParticipantPicker();
    };

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
                    <div
                        className="icon-box"
                        aria-label="Список контестов"
                    >
                        <FiMenu className="sidebar-contests-menu-icon" aria-hidden />
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

            <div
                className={`menu-item-container sidebar-admin-user${isUserMenuOpen ? " sidebar-admin-user--open" : ""}`}
                onMouseEnter={() => setHoveredMenu("user")}
                onMouseLeave={closeUserMenu}
            >
                <div
                    className="icon-box sidebar-admin-user-trigger"
                    aria-label="Выбор участника"
                >
                    <img src={adminUserIcon} alt="" className="sidebar-admin-user-icon" />
                </div>
                <Submenu isOpen={isUserMenuOpen}>
                    <ParticipantFollowMenu adminSession={adminSession} />
                </Submenu>
            </div>
        </aside>
    );
}
