import { Save, Undo2 } from "lucide-react";
import { ContestContextBar } from "./ContestContextBar";
import { confirmAdvance } from "./ConfirmDialogs";
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
        if (!confirmAdvance()) {
            return;
        }
        await tt.advance();
    };

    const handleSave = async () => {
        await tt.saveTimetable();
    };

    const showTimeline =
        tt.canEditSchedule && (tt.displaySegments.length > 0 || tt.hasSchedule);
    const showEmpty = tt.canEditSchedule && !showTimeline && tt.loadState !== "loading";
    const scheduleLocked = !tt.canEditSchedule;
    const showBanner =
        tt.canEditSchedule &&
        tt.hasSchedule &&
        tt.view != null;

    return (
        <section className="admin-timetable tt-page">
            <ContestContextBar
                contestId={tt.contestId}
                onContestChange={handleContestChange}
                view={tt.view}
            />

            {showBanner && (
                <NextTourBanner
                    view={tt.view!}
                    busy={tt.busy}
                    onStartNow={() => void handleStartNow()}
                    showAutostart={tt.hasSchedule}
                    onAutoStartChange={(enabled) => void tt.setAutoStartEnabled(enabled)}
                />
            )}

            {scheduleLocked && (
                <p className="tt-message tt-message--info">
                    Выберите контест, чтобы настроить расписание туров.
                </p>
            )}

            {tt.canEditSchedule && tt.loadState === "loading" && !showTimeline && (
                <p className="tt-message tt-message--info">Загрузка…</p>
            )}

            {showEmpty && (
                <EmptySchedule onApplyTemplate={tt.applyDurationTemplate} disabled={tt.busy} />
            )}

            {showTimeline && (
                <div className="tt-schedule-panel">
                    <ScheduleAxis
                        segments={tt.displaySegments}
                        view={tt.view}
                        dirty={tt.dirty}
                        busy={tt.busy}
                        onPendingDurationChange={tt.setPendingDuration}
                        onPendingKindChange={tt.setPendingKind}
                        onPendingRemove={tt.removeSlot}
                        onAddSlot={tt.addSlot}
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
