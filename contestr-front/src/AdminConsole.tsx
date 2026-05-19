import { useState } from "react";
import TimetablePage from "./admin-timetable/TimetablePage";
import ContestsAdminPage from "./admin-contests/ContestsAdminPage";
import "./AdminConsole.css";

type AdminTab = "timetable" | "contests";

export default function AdminConsole() {
    const [tab, setTab] = useState<AdminTab>("contests");

    return (
        <div className="admin-console">
            <nav className="admin-console-tabs" aria-label="Разделы админки">
                <button
                    type="button"
                    className={`admin-console-tab${tab === "contests" ? " admin-console-tab--active" : ""}`}
                    onClick={() => setTab("contests")}
                >
                    Контесты Codeforces
                </button>
                <button
                    type="button"
                    className={`admin-console-tab${tab === "timetable" ? " admin-console-tab--active" : ""}`}
                    onClick={() => setTab("timetable")}
                >
                    Расписание туров
                </button>
            </nav>
            <div className="admin-console-body">
                {tab === "contests" ? <ContestsAdminPage /> : <TimetablePage />}
            </div>
        </div>
    );
}
