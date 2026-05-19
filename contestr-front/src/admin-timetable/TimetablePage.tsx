import { Save, Undo2 } from "lucide-react";
import { ContestContextBar } from "./ContestContextBar";
import { confirmStartNow } from "./ConfirmDialogs";
import { EmptySchedule } from "./EmptySchedule";
import { NextTourBanner } from "./NextTourBanner";
import { ScheduleAxis } from "./ScheduleAxis";
import { useTimetable } from "./useTimetable";
import "./TimetablePage.css";

export default function TimetablePage() {
    const tt = useTimetable();

    const handleContestChange = (id: number) => {
        if (Number.isInteger(id) && id > 0) {
            tt.setContestId(id);
        }
    };

    const handleStartNow = async () => {
        const n = tt.view?.next_tour_number;
        if (!n || !confirmStartNow(n)) {
            return;
        }
        await tt.startNextTour();
    };

    const handleSave = async () => {
        await tt.saveTimetable();
    };

    const showTimeline = tt.draftTours.length > 0 || tt.hasSchedule;
    const showEmpty = !showTimeline && tt.loadState !== "loading";

    return (
        <section className="admin-timetable tt-page">
            <ContestContextBar
                contestId={tt.contestId}
                onContestChange={handleContestChange}
                view={tt.view}
                hasSchedule={tt.hasSchedule}
                busy={tt.busy}
                onAutoStartChange={(enabled) => void tt.setAutoStartEnabled(enabled)}
            />

            {tt.view && tt.view.next_tour_number && (
                <NextTourBanner view={tt.view} busy={tt.busy} onStartNow={() => void handleStartNow()} />
            )}

            {tt.loadState === "loading" && !showTimeline && (
                <p className="tt-message tt-message--info">Загрузка…</p>
            )}

            {showEmpty && (
                <EmptySchedule onApplyTemplate={tt.applyDurationTemplate} disabled={tt.busy} />
            )}

            {showTimeline && (
                <div className="tt-schedule-panel">
                    <ScheduleAxis
                        tours={tt.draftTours}
                        view={tt.view}
                        dirty={tt.dirty}
                        busy={tt.busy}
                        onDurationChange={tt.setDuration}
                        onRemove={tt.removeTour}
                        onAddTour={tt.addTour}
                    />
                    <footer className="tt-footer admin-toolbar">
                        <button
                            type="button"
                            className="admin-icon-btn admin-primary-btn tt-icon-action-btn"
                            onClick={() => void handleSave()}
                            disabled={tt.busy || !tt.dirty}
                            aria-label="Сохранить изменения"
                            title="Сохранить изменения"
                        >
                            <Save size={16} aria-hidden />
                        </button>
                        <button
                            type="button"
                            className="admin-icon-btn tt-icon-action-btn"
                            onClick={tt.revertChanges}
                            disabled={tt.busy || !tt.dirty}
                            aria-label="Отменить изменения"
                            title="Отменить изменения"
                        >
                            <Undo2 size={16} aria-hidden />
                        </button>
                    </footer>
                </div>
            )}

            {tt.message && tt.messageKind === "error" && (
                <p className="admin-login-message admin-login-message--error tt-message" role="alert">
                    Ошибка: {tt.message}
                </p>
            )}
        </section>
    );
}
