import { PiBalloonThin } from "react-icons/pi";
import { TfiCup } from "react-icons/tfi";
import { CONTEST_IDS } from "../consts";
// import Time_cont from './Time_cont.tsx'


export const Sidebar = () => {
    return (
        <aside className="sidebar-hover">
            <div className="icon-wrap">
                <div className="icon-box">
                    <img src="/logo.svg" alt="Логотип" className="logo-icon" />
                    {/* <span className="tooltip">Главная</span> */}
                </div>
            </div>

            <div className="icon-wrap">
                <div className="icon-box" onClick={() => {
                    window.location.href = "/?contestId=" + CONTEST_IDS[0]
                }}>
                    <PiBalloonThin size={45} />
                    <span className="tooltip">Младшие параллели</span>
                </div>
                <div className="icon-box" onClick={() => {
                    window.location.href = "/?contestId=" + CONTEST_IDS[1]
                }}>
                   <TfiCup size={30} />
                   
                   <span className="tooltip">Старшие параллели</span>
                </div>
            
                </div>
            {/* <div className='time'>
                <Time_cont/>
            </div> */}
        </aside>
    )
}
