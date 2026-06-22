import { useEffect } from "react";
import {
    ADMIN_CONTESTS_PATH,
    ADMIN_ROOT_PATH,
    ADMIN_TIMETABLE_PATH,
    useAppPath,
} from "@/app/AppPath";
import TimetablePage from "@/features/admin/timetable/TimetablePage";
import ContestsAdminPage from "@/features/admin/contests/ContestsAdminPage";
import "./AdminConsole.css";

type AdminTab = "timetable" | "contests";

function resolveAdminTab(path: string): AdminTab | null {
    if (path === ADMIN_TIMETABLE_PATH) {
        return "timetable";
    }
    if (path === ADMIN_CONTESTS_PATH) {
        return "contests";
    }
    return null;
}

export default function AdminConsole() {
    const { path, navigate } = useAppPath();
    const tab = resolveAdminTab(path);

    useEffect(() => {
        if (path === ADMIN_ROOT_PATH || (path.startsWith(`${ADMIN_ROOT_PATH}/`) && tab === null)) {
            navigate(ADMIN_TIMETABLE_PATH);
        }
    }, [path, tab, navigate]);

    if (tab === null) {
        return null;
    }

    return (
        <div className="admin-console">
            <nav className="admin-console-tabs" aria-label="Разделы админки">
                <a
                    href={ADMIN_TIMETABLE_PATH}
                    className={`admin-console-tab${tab === "timetable" ? " admin-console-tab--active" : ""}`}
                    aria-current={tab === "timetable" ? "page" : undefined}
                    onClick={(event) => {
                        event.preventDefault();
                        navigate(ADMIN_TIMETABLE_PATH);
                    }}
                >
                    Расписание туров
                </a>
                <a
                    href={ADMIN_CONTESTS_PATH}
                    className={`admin-console-tab${tab === "contests" ? " admin-console-tab--active" : ""}`}
                    aria-current={tab === "contests" ? "page" : undefined}
                    onClick={(event) => {
                        event.preventDefault();
                        navigate(ADMIN_CONTESTS_PATH);
                    }}
                >
                    Контесты Codeforces
                </a>
            </nav>
            <div className="admin-console-body">
                {tab === "timetable" ? <TimetablePage /> : <ContestsAdminPage />}
            </div>
        </div>
    );
}
