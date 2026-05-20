import type { TimetableView } from "@/client/types.gen";
import { ChevronDown, ChevronRight } from "lucide-react";
import { formatContestClock, formatElapsed } from "./time";
import { useContestElapsed } from "./useContestElapsed";

type Props = {
    name: string;
    contestId: number;
    view: TimetableView | null;
    collapsed: boolean;
    onToggleCollapsed: () => void;
};

export function ContestTimetableHeader({
    name,
    contestId,
    view,
    collapsed,
    onToggleCollapsed,
}: Props) {
    const elapsed = useContestElapsed(view?.contest_start_time, view?.elapsed_seconds ?? 0);

    return (
        <header className={`tt-contest-header${collapsed ? " tt-contest-header--collapsed" : ""}`}>
            <h2 className="tt-contest-header__title">
                <button
                    type="button"
                    className="tt-contest-header__toggle-btn"
                    onClick={onToggleCollapsed}
                    title={collapsed ? "Развернуть" : "Свернуть"}
                    aria-label={collapsed ? "Развернуть контест" : "Свернуть контест"}
                    aria-expanded={!collapsed}
                >
                    {collapsed ? (
                        <ChevronRight size={16} strokeWidth={2.25} aria-hidden />
                    ) : (
                        <ChevronDown size={16} strokeWidth={2.25} aria-hidden />
                    )}
                </button>
                {name}
                <span className="tt-contest-header__id">#{contestId}</span>
            </h2>
            {!collapsed && (
            <div className="tt-contest-header__meta">
                <p className="tt-contest-header__start">
                    Старт контеста: {formatContestClock(view?.contest_start_time)}
                </p>
                {view?.contest_start_time && (
                    <p className="tt-contest-header__elapsed" aria-live="polite">
                        Прошло: <strong>{formatElapsed(elapsed)}</strong>
                    </p>
                )}
            </div>
            )}
        </header>
    );
}
