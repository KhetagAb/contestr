import { useState, useRef, useEffect } from "react";
import { FiFileText } from "react-icons/fi";
import { LuTrophy } from "react-icons/lu";
import { CONTEST_IDS } from "../consts";
import Time_cont from "./Time_cont.tsx";
import badminton from "../../SportComponents/public/badminton.svg";
import basketball from "../../SportComponents/public/basketball.svg";


// Компонент подменю
const Submenu = ({
                     isOpen,
                     children,
                 }: {
    isOpen: boolean;
    children: React.ReactNode;
}) => {
    const ref = useRef<HTMLDivElement>(null);
    const [height, setHeight] = useState("0px");

    useEffect(() => {
        if (ref.current) {
            if (isOpen) {
                setHeight(ref.current.scrollHeight + "px"); // раскрываем по фактической высоте
            } else {
                setHeight("0px"); // закрываем
            }
        }
    }, [isOpen, children]);

    return (
        <div
            className="submenu-wrapper"
            style={{ maxHeight: height }}
        >
            <div ref={ref} className="submenu">
                {children}
            </div>
        </div>
    );
};

export const Sidebar = () => {
    const [openMenu, setOpenMenu] = useState<string | null>(null);

    const toggleMenu = (menu: string) => {
        setOpenMenu(openMenu === menu ? null : menu);
    };

    return (
        <aside className="sidebar-hover">
            <div className="icon-wrap">
                <div className="icon-box">
                    <img src="/logo.svg" alt="Логотип" className="logo-icon" />
                </div>
            </div>

            <div className="icon-wrap">
                <div className="menu-item-container">
                    <div className="icon-box" onClick={() => toggleMenu("contests")}>
                        <FiFileText size={35} />
                        <span className="tooltip">Контесты</span>
                    </div>
                    <Submenu isOpen={openMenu === "contests"}>
                        <div
                            className="submenu-item"
                            onClick={() =>
                                (window.location.href = "/?contestId=" + CONTEST_IDS[0])
                            }
                        >
                            Дети
                        </div>
                        <div className="submenu-item">Взрослые</div>
                    </Submenu>
                </div>

                <div className="menu-item-container">
                    <div className="icon-box" onClick={() => toggleMenu("sport")}>
                        <LuTrophy size={35} />
                        <span className="tooltip">Спорт</span>
                    </div>
                    <Submenu isOpen={openMenu === "sport"}>
                        <div
                            className="submenu-item"
                            onClick={() =>
                                (window.location.href = "/?contestId=" + CONTEST_IDS[1])
                            }
                        >
                            <img src={badminton} alt="Бадминтон" className="sport-icon" />
                            <span className="tooltip">Бадминтон</span>
                        </div>
                        <div className="submenu-item">
                            <img src={basketball} alt="Бадминтон" className="sport-icon" />
                            <span className="tooltip">basketball</span>
                        </div>
                        <div className="submenu-item">
                            <img src={basketball} alt="Бадминтон" className="sport-icon" />
                            <span className="tooltip">basketball</span>
                        </div>
                        <div className="submenu-item">
                            <img src={basketball} alt="Бадминтон" className="sport-icon" />
                            <span className="tooltip">basketball</span>
                        </div>
                        <div className="submenu-item">
                            <img src={basketball} alt="Бадминтон" className="sport-icon" />
                            <span className="tooltip">basketball</span>
                        </div>
                        <div className="submenu-item">
                            <img src={basketball} alt="Бадминтон" className="sport-icon" />
                            <span className="tooltip">basketball</span>
                        </div>
                    </Submenu>
                </div>
            </div>

            <div className="time">
                <Time_cont />
            </div>
        </aside>
    );
};
