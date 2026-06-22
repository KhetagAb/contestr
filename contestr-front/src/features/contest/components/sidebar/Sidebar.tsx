import { useState, type ReactNode } from "react";
import { FiMenu } from "react-icons/fi";
import { CONTEST_HOME_QUERY, useAppPath } from "@/app/AppPath";
import { useFollowedParticipant } from "@/features/contest/follow/FollowedParticipantContext";
import { useContestTimetableStarts } from "@/shared/hooks/useContestTimetableStarts";
import { useContests } from "@/shared/hooks/useContests";
import { formatContestStartLabel } from "@/shared/utils/contestStartTime";
import logo from "@/assets/icons/logo.svg";
import adminUserIcon from "@/assets/images/admin-user-icon.png";
import { AdminSidebarMenu } from "./AdminSidebarMenu";

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
    const { navigate } = useAppPath();
    const [hoveredMenu, setHoveredMenu] = useState<string | null>(null);
    const { openParticipantPicker } = useFollowedParticipant();
    const { contests, isLoading } = useContests();
    const contestsMenuOpen = hoveredMenu === "contests";
    const contestIds = contests.map((item) => item.contest_id);
    const { startTimeByContestId } = useContestTimetableStarts(
        contestIds,
        contestsMenuOpen && !isLoading && contestIds.length > 0,
    );

    const isAdminMenuOpen = hoveredMenu === "user" && adminSession != null;

    return (
        <aside className="sidebar-hover">
            <div className="icon-wrap">
                <a
                    href={`/?${CONTEST_HOME_QUERY}`}
                    className="icon-box sidebar-home-link"
                    aria-label="На главную"
                    onClick={(event) => {
                        event.preventDefault();
                        navigate(`/?${CONTEST_HOME_QUERY}`);
                    }}
                >
                    <img src={logo} alt="" className="logo-icon" />
                </a>
            </div>

            <div className="icon-wrap">
                <div
                    className="menu-item-container"
                    onMouseEnter={() => setHoveredMenu("contests")}
                    onMouseLeave={() => setHoveredMenu(null)}
                >
                    <div className="icon-box" aria-label="Список контестов">
                        <FiMenu className="sidebar-contests-menu-icon" aria-hidden />
                    </div>
                    <Submenu isOpen={contestsMenuOpen}>
                        {isLoading && (
                            <span className="submenu-item submenu-item-text">Загрузка…</span>
                        )}
                        {!isLoading && contests.length === 0 && (
                            <span className="submenu-item submenu-item-text">
                                Контесты не настроены
                            </span>
                        )}
                        {contests.map((item) => {
                            const startLabel = formatContestStartLabel(
                                startTimeByContestId.get(item.contest_id),
                            );
                            return (
                                <a
                                    key={item.contest_id}
                                    className="submenu-item submenu-item--contest"
                                    href={`/?contestId=${item.contest_id}`}
                                >
                                    <span className="submenu-item__body">
                                        <span className="submenu-item-text submenu-item__title">
                                            {item.name}
                                        </span>
                                        {startLabel && (
                                            <span className="submenu-item__meta">{startLabel}</span>
                                        )}
                                    </span>
                                </a>
                            );
                        })}
                    </Submenu>
                </div>
            </div>

            <div
                className={`menu-item-container sidebar-admin-user${isAdminMenuOpen ? " sidebar-admin-user--open" : ""}`}
                onMouseEnter={() => setHoveredMenu("user")}
                onMouseLeave={() => setHoveredMenu(null)}
            >
                <button
                    type="button"
                    className="icon-box sidebar-admin-user-trigger"
                    aria-label="Найти себя в таблице"
                    onClick={() => openParticipantPicker()}
                >
                    <img src={adminUserIcon} alt="" className="sidebar-admin-user-icon" />
                </button>
                {adminSession ? (
                    <Submenu isOpen={isAdminMenuOpen}>
                        <AdminSidebarMenu session={adminSession} />
                    </Submenu>
                ) : null}
            </div>
        </aside>
    );
}
