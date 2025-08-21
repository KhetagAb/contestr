import { useState } from "react";
import { FiFileText } from "react-icons/fi";
import { LuTrophy } from "react-icons/lu";
import { CONTEST_IDS } from "../consts";
import Time_cont from "./Time_cont.tsx";
// import badminton from './badminton.svg';
// import basketball from 'front/SportComponents/public/basketball.svg'
// import football from 'front/SportComponents/public/football.svg'
// import {table_tennis} from './table_tennis.svg'
// import {tennis} from './tennis.svg'
// import {volleyball} from './volleyball.svg'
// import telegramLogo from "../../SportComponents/public/telegramLogo.svg";

export const Sidebar = () => {
    const [openMenu, setOpenMenu] = useState<string | null>(null);

    const toggleMenu = (menu: string) => {
        setOpenMenu(openMenu === menu ? null : menu);
    };

    return (
        <aside className="sidebar-hover">
            {/* Логотип */}
            <div className="icon-wrap">
                <div className="icon-box">
                    <img src="/logo.svg" alt="Логотип" className="logo-icon" />
                </div>
            </div>

            <div className="icon-wrap">
                <div className="icon-box" onClick={() => toggleMenu("junior")}>
                    <FiFileText size={35} />
                    <span className="tooltip">Контесты</span>
                </div>
                <div className={`submenu-wrapper ${openMenu === "junior" ? "open" : ""}`}>
                    <div className="submenu">
                        <div
                            className="submenu-item"
                            onClick={() => (window.location.href = "/?contestId=" + CONTEST_IDS[0])}
                        >
                            Дети
                        </div>
                        <div className="submenu-item">Взрослые</div>
                    </div>
                </div>

                {/* Старшие */}
                <div className="icon-box" onClick={() => toggleMenu("senior")}>
                    <LuTrophy size={35} />
                    <span className="tooltip">Спорт</span>
                </div>
                <div className={`submenu-wrapper ${openMenu === "senior" ? "open" : ""}`}>
                    <div className="submenu">
                        <div
                            className="submenu-item"
                            onClick={() => (window.location.href = "/?contestId=" + CONTEST_IDS[1])}
                        >
                            <p></p>
                        </div>
                        <div className="submenu-item">
                            <p></p>

                        </div>

                </div>
            </div>
            </div>
            <div className="time">
                <Time_cont />
            </div>
        </aside>
    );
};
