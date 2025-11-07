import { useState, useRef, useEffect } from "react";
import { createPortal } from "react-dom";
import { FiFileText } from "react-icons/fi";
// import { LuTrophy } from "react-icons/lu";
import { GiBalloonDog, GiLightBulb } from "react-icons/gi";
import { CONTEST_IDS } from "../consts";
import Clock from "./Time_cont.tsx";
import type { GetRegattaContestStandingsResponses } from "../client";
// import badminton from "./public/badminton.svg";
// import football from "./public/football.svg";
// import tableTennis from "./public/table_tennis.svg";
// import tennis from "./public/tennis.svg";
// import volleyball from "./public/volleyball.svg";
// import basketball from "./public/basketball.svg";
// import shooting from "./public/shooting.svg";
import logo from "./public/logo.svg";

const mockData: GetRegattaContestStandingsResponses["200"] = {
    contest_name: "Regatta ITMO 2025",
    contest_id: 42,
    current_tour_start_time: Math.floor(Date.now() / 1000) - 600, // старт 10 мин назад
    current_tour_duration: 120, // 2 часа
    rows: [], // таблица нам не нужна
};

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
                setHeight(ref.current.scrollHeight + "px");
            } else {
                setHeight("0px");
            }
        }
    }, [isOpen, children]);

    return (
        <div className="submenu-wrapper" style={{ maxHeight: height }}>
            <div ref={ref} className="submenu">
                {children}
            </div>
        </div>
    );
};

// Универсальный тултип через портал
const SubTooltip = ({
                        text,
                        target,
                    }: {
    text: string;
    target: HTMLElement | null;
}) => {
    if (!target) return null;

    const rect = target.getBoundingClientRect();

    const style: React.CSSProperties = {
        position: "fixed",
        top: rect.top + window.scrollY + rect.height / 2,
        left: rect.left + window.scrollX + rect.width + 10,
        transform: "translateY(-50%)",
        background: "#f0eee1",
        fontFamily: "'Montserrat Alternates', sans-serif", // <- шрифт
        fontSize: "var(--fs-tooltip)",
        padding: "0.5em 1em",
        borderRadius: "1em",
        whiteSpace: "nowrap",
        zIndex: 9999,
        pointerEvents: "none",
    };

    return createPortal(<div style={style}>{text}</div>, document.body);
};

export const Sidebar = () => {
    const [openMenu, setOpenMenu] = useState<string | null>(null);
    const [hoveredItem, setHoveredItem] = useState<HTMLElement | null>(null);
    const [tooltipText, setTooltipText] = useState("");

    const toggleMenu = (menu: string) =>
        setOpenMenu(openMenu === menu ? null : menu);

    return (
        <aside className="sidebar-hover">
            <div className="icon-wrap">
                <div className="icon-box">
                    <img src={logo} alt="Логотип" className="logo-icon" />
                </div>
            </div>

            <div className="icon-wrap">
                {/* Контесты */}
                <div className="menu-item-container">
                    <div className="icon-box" onClick={() => toggleMenu("contests")}>
                        <FiFileText size={35} />
                        <span className="tooltip">Контесты</span>
                    </div>
                    <Submenu isOpen={openMenu === "contests"}>
                        {[
                            { icon: <GiBalloonDog className="subMenu-icon" />, alt: "Младший дивизион" },
                            { icon: <GiLightBulb className="subMenu-icon" />, alt: "Старший дивизион" },
                        ].map((item, idx) => (
                            <div
                                key={idx}
                                className="submenu-item"
                                onMouseEnter={(e) => {
                                    setHoveredItem(e.currentTarget);
                                    setTooltipText(item.alt);
                                }}
                                onMouseLeave={() => setHoveredItem(null)}
                                onClick={() =>
                                    (window.location.href = "/?contestId=" + CONTEST_IDS[idx])
                                }
                            >
                                {item.icon}
                            </div>
                        ))}
                    </Submenu>
                </div>

                {/*<div className="menu-item-container">*/}
                {/*    <div className="icon-box" onClick={() => toggleMenu("sport")}>*/}
                {/*        <LuTrophy size={35} />*/}
                {/*        <span className="tooltip">Спорт</span>*/}
                {/*    </div>*/}
                {/*    <Submenu isOpen={openMenu === "sport"}>*/}
                {/*        {[*/}
                {/*            { src: badminton, alt: "Бадминтон" },*/}
                {/*            { src: basketball, alt: "Баскетбол" },*/}
                {/*            { src: football, alt: "Футбол" },*/}
                {/*            { src: tableTennis, alt: "Настольный теннис" },*/}
                {/*            { src: tennis, alt: "Теннис" },*/}
                {/*            { src: volleyball, alt: "Волейбол" },*/}
                {/*            { src: shooting, alt: "Стрельба" },*/}
                {/*        ].map((sport) => (*/}
                {/*            <div*/}
                {/*                key={sport.alt}*/}
                {/*                className="submenu-item"*/}
                {/*                onMouseEnter={(e) => {*/}
                {/*                    setHoveredItem(e.currentTarget);*/}
                {/*                    setTooltipText(sport.alt);*/}
                {/*                }}*/}
                {/*                onMouseLeave={() => setHoveredItem(null)}*/}
                {/*            >*/}
                {/*                <img src={sport.src} alt={sport.alt} className="sport-icon" />*/}
                {/*            </div>*/}
                {/*        ))}*/}
                {/*    </Submenu>*/}
                {/*</div>*/}
            </div>

            <div className="time">
                <Clock data={mockData} />

            </div>

            <div className="time-text"><p><span>Время с </span> <span>начала тура </span></p></div>

            <SubTooltip text={tooltipText} target={hoveredItem} />
        </aside>
    );
};
