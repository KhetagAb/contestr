import { Save, Undo2 } from "lucide-react";
import { confirmAdvance } from "./ConfirmDialogs";
import { ContestTimetableHeader } from "./ContestTimetableHeader";
import { EmptySchedule } from "./EmptySchedule";
import { NextTourBanner } from "./NextTourBanner";
import { ScheduleAxis } from "./ScheduleAxis";
import { useContestTimetable } from "./useContestTimetable";

type Props = {
    contestId: number;
    contestName: string;
    collapsed: boolean;
    onToggleCollapsed: () => void;
};

export function ContestTimetableBlock({
    contestId,
    contestName,
    collapsed,
    onToggleCollapsed,
}: Props) {
    const tt = useContestTimetable(contestId);

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
        tt.displaySegments.length > 0 || tt.hasSchedule;
    const showEmpty = !showTimeline && tt.loadState !== "loading";
    const showBanner = tt.hasSchedule && tt.view != null;

    return (
        <article className={`tt-contest-block${collapsed ? " tt-contest-block--collapsed" : ""}`}>
            <ContestTimetableHeader
                name={contestName}
                contestId={contestId}
                view={tt.view}
                collapsed={collapsed}
                onToggleCollapsed={onToggleCollapsed}
            />

            {!collapsed && showBanner && (
                <NextTourBanner
                    view={
                        tt.dirty
                            ? {
                                  ...tt.view!,
                                  timeline_segments: tt.displaySegments,
                                  pending_slots: tt.draftPending,
                              }
                            : tt.view!
                    }
                    busy={tt.busy}
                    onStartNow={() => void handleStartNow()}
                    showAutostart={tt.hasSchedule}
                    onAutoStartChange={(enabled) => void tt.setAutoStartEnabled(enabled)}
                />
            )}

            {!collapsed && tt.loadState === "loading" && !showTimeline && (
                <p className="tt-message tt-message--info">Загрузка…</p>
            )}

            {!collapsed && showEmpty && (
                <EmptySchedule
                    onApplyTemplate={tt.applyDurationTemplate}
                    disabled={tt.busy}
                />
            )}

            {!collapsed && showTimeline && (
                <div className="tt-schedule-panel">
                    <ScheduleAxis
                        segments={tt.displaySegments}
                        view={tt.view}
                        dirty={tt.dirty}
                        busy={tt.busy}
                        onPendingDurationChange={tt.setPendingDuration}
                        onActiveDurationChange={tt.updateActiveDuration}
                        onPendingKindChange={tt.setPendingKind}
                        onPendingRemove={tt.removeSlot}
                        onAddSlotAfter={tt.addSlotAfter}
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

            {!collapsed && tt.message && tt.messageKind === "error" && (
                <p
                    className="admin-login-message admin-login-message--error tt-message"
                    role="alert"
                >
                    Ошибка: {tt.message}
                </p>
            )}
            {!collapsed && tt.message && tt.messageKind === "success" && (
                <p className="tt-message tt-message--success" role="status">
                    {tt.message}
                </p>
            )}
        </article>
    );
}
