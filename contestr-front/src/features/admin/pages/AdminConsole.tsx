import { useState } from "react";
import TimetablePage from "@/features/admin/timetable/TimetablePage";
import ContestsAdminPage from "@/features/admin/contests/ContestsAdminPage";
import "./AdminConsole.css";

type AdminTab = "timetable" | "contests";

export default function AdminConsole() {
    const [tab, setTab] = useState<AdminTab>("timetable");

    return (
        <div className="admin-console">
            <nav className="admin-console-tabs" aria-label="Разделы админки">
                <button
                    type="button"
                    className={`admin-console-tab${tab === "timetable" ? " admin-console-tab--active" : ""}`}
                    onClick={() => setTab("timetable")}
                >
                    Расписание туров
                </button>
                <button
                    type="button"
                    className={`admin-console-tab${tab === "contests" ? " admin-console-tab--active" : ""}`}
                    onClick={() => setTab("contests")}
                >
                    Контесты Codeforces
                </button>
            </nav>
            <div className="admin-console-body">
                {tab === "timetable" ? <TimetablePage /> : <ContestsAdminPage />}
            </div>
        </div>
    );
}
