import { RefreshCw } from "lucide-react";
import { useContests } from "@/shared/hooks/useContests";
import type { TimetableView } from "@/client/types.gen";
import { formatContestClock, formatElapsed } from "./time";
import { useContestElapsed } from "./useContestElapsed";

type Props = {
    contestId: number;
    onContestChange: (id: number) => void;
    view: TimetableView | null;
    onRefresh: () => void;
    busy: boolean;
};

export function ContestContextBar({ contestId, onContestChange, view, onRefresh, busy }: Props) {
    const { contests, isLoading } = useContests();
    const elapsed = useContestElapsed(view?.contest_start_time, view?.elapsed_seconds ?? 0);

    return (
        <div className="tt-context-bar">
            <div className="tt-context-bar__row">
                <p className="tt-contest-start">
                    Старт контеста: {formatContestClock(view?.contest_start_time)}
                </p>
                <div className="tt-context-bar__controls">
                    {view?.contest_start_time && (
                        <p className="tt-elapsed" aria-live="polite">
                            Прошло: <strong>{formatElapsed(elapsed)}</strong>
                        </p>
                    )}
                    <label className="tt-contest-select">
                        <span className="visually-hidden">Контест</span>
                        <select
                            value={contestId || ""}
                            onChange={(e) => onContestChange(Number(e.target.value))}
                            disabled={isLoading || contests.length === 0}
                        >
                            {contests.length === 0 && (
                                <option value="">Контесты не настроены</option>
                            )}
                            {contests.map((c) => (
                                <option key={c.contest_id} value={c.contest_id}>
                                    {c.name} ({c.contest_id})
                                </option>
                            ))}
                        </select>
                    </label>
                    <button
                        type="button"
                        className="tt-refresh-contest-btn"
                        title="Обновить данные контеста"
                        aria-label="Обновить данные контеста"
                        onClick={onRefresh}
                        disabled={busy || !contestId}
                    >
                        <RefreshCw size={16} aria-hidden />
                    </button>
                </div>
            </div>
        </div>
    );
}
