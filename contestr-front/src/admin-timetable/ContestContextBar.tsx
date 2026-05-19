import { useEffect, useState } from "react";
import { CONTESTS } from "../consts";
import type { TimetableView } from "../client/types.gen";
import { formatContestClock, formatElapsed } from "./time";

type Props = {
    contestId: number;
    onContestChange: (id: number) => void;
    view: TimetableView | null;
    hasSchedule?: boolean;
    busy?: boolean;
    onAutoStartChange?: (enabled: boolean) => void;
};

export function ContestContextBar({
    contestId,
    onContestChange,
    view,
    hasSchedule = false,
    busy = false,
    onAutoStartChange,
}: Props) {
    const [elapsed, setElapsed] = useState(view?.elapsed_seconds ?? 0);

    useEffect(() => {
        if (!view?.contest_start_time) {
            setElapsed(0);
            return;
        }
        const startMs = new Date(view.contest_start_time).getTime();
        const tick = () => {
            setElapsed(Math.max(0, Math.floor((Date.now() - startMs) / 1000)));
        };
        tick();
        const id = window.setInterval(tick, 1000);
        return () => window.clearInterval(id);
    }, [view?.contest_start_time, view?.elapsed_seconds]);

    const showAutostart = hasSchedule && Boolean(onAutoStartChange);
    const autostartAvailable = view?.auto_start_available ?? false;
    const autostartOn =
        autostartAvailable && Boolean(view?.auto_start_enabled);

    return (
        <header className="tt-context-bar">
            <div className="tt-context-bar__top">
                <h2 className="tt-title">Расписание туров</h2>
                <div className="tt-context-bar__controls">
                    <label className="tt-contest-select">
                        <span className="visually-hidden">Контест</span>
                        <select
                            value={contestId}
                            onChange={(e) => onContestChange(Number(e.target.value))}
                        >
                            {CONTESTS.map((c) => (
                                <option key={c.id} value={c.id}>
                                    {c.name} ({c.id})
                                </option>
                            ))}
                        </select>
                    </label>
                    {view?.contest_start_time && (
                        <p className="tt-elapsed" aria-live="polite">
                            Прошло: <strong>{formatElapsed(elapsed)}</strong>
                        </p>
                    )}
                </div>
            </div>
            <div className="tt-context-bar__meta">
                <p className="tt-contest-start">
                    Старт контеста: {formatContestClock(view?.contest_start_time)}
                </p>
                {showAutostart && (
                    <label
                        className="tt-autostart-toggle"
                        title={
                            autostartAvailable
                                ? undefined
                                : "На сервере не настроен фоновый timetable_sync"
                        }
                    >
                        <span className="tt-autostart-toggle__label">
                            {autostartOn ? "Автозапуск включён" : "Автозапуск выключен"}
                        </span>
                        <input
                            type="checkbox"
                            role="switch"
                            className="tt-autostart-toggle__input"
                            checked={autostartOn}
                            disabled={busy || !autostartAvailable}
                            onChange={(e) => onAutoStartChange?.(e.target.checked)}
                        />
                        <span className="tt-autostart-switch" aria-hidden />
                    </label>
                )}
            </div>
        </header>
    );
}
