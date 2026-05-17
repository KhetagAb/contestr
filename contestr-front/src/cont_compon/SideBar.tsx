import { useState } from "react";
// import { createPortal } from "react-dom";
import { FiFileText } from "react-icons/fi";
// import { LuTrophy } from "react-icons/lu";
// import { GiBalloonDog, GiLightBulb } from "react-icons/gi";
import { CONTESTS } from "../consts";
import Clock from "./Time_cont.tsx";
// import badminton from "./public/badminton.svg";
// import football from "./public/football.svg";
// import tableTennis from "./public/table_tennis.svg";
// import tennis from "./public/tennis.svg";
// import volleyball from "./public/volleyball.svg";
// import basketball from "./public/basketball.svg";
// import shooting from "./public/shooting.svg";
import logo from "./public/logo.svg";

const Submenu = ({
                     isOpen,
                     children,
                 }: {
    isOpen: boolean;
    children: React.ReactNode;
}) => (
    <div className={`submenu-wrapper${isOpen ? " submenu-wrapper--open" : ""}`}>
        {children}
    </div>
);

// // Универсальный тултип через портал
// const SubTooltip = ({
//                         text,
//                         target,
//                     }: {
//     text: string;
//     target: HTMLElement | null;
// }) => {
//     if (!target) return null;

//     const rect = target.getBoundingClientRect();

//     const style: React.CSSProperties = {
//         position: "fixed",
//         top: rect.top + window.scrollY + rect.height / 2,
//         left: rect.left + window.scrollX + rect.width + 10,
//         transform: "translateY(-50%)",
//         background: "#f0eee1",
//         fontFamily: "'Montserrat Alternates', sans-serif", // <- шрифт
//         fontSize: "var(--fs-tooltip)",
//         padding: "0.5em 1em",
//         borderRadius: "1em",
//         whiteSpace: "nowrap",
//         zIndex: 9999,
//         pointerEvents: "none",
//     };

//     return createPortal(<div style={style}>{text}</div>, document.body);
// };

export const Sidebar = () => {
    const [hoveredMenu, setHoveredMenu] = useState<string | null>(null);

    return (
        <aside className="sidebar-hover">
            <div className="icon-wrap">
                <div className="icon-box">
                    <img src={logo} alt="Логотип" className="logo-icon" />
                </div>
            </div>

            <div className="icon-wrap">
                {/* Контесты */}
                <div
                    className="menu-item-container"
                    onMouseEnter={() => setHoveredMenu("contests")}
                    onMouseLeave={() => setHoveredMenu(null)}
                >
                    <div className="icon-box">
                        <FiFileText size={35} />
                    </div>
                    <Submenu isOpen={hoveredMenu === "contests"}>
                        {
                            CONTESTS.map((item) => (
                                <a key={item.id} className="submenu-item"
                                href={`/?contestId=${item.id}`}>
                                    <item.IconComponent className="subMenu-icon"></item.IconComponent>
                                    <span className="submenu-item-text">{item.name}</span>
                                </a>
                            ))
                        }
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
                <Clock/>

            </div>

            <div className="time-text"><p><span>Время с </span> <span>начала тура </span></p></div>

            {/* <SubTooltip text={tooltipText} target={hoveredItem} /> */}
        </aside>
    );
};
