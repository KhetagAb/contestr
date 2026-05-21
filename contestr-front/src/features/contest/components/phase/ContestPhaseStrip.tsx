import { CalendarClock, Coffee } from "lucide-react";
import type { TimetableView } from "@/client/types.gen";
import { useContestTimetable } from "@/features/contest/hooks/useContestTimetable";
import { derivePhaseDisplay } from "@/shared/timetable/phaseState";
import { MiniTimeline } from "@/shared/timetable/MiniTimeline";
import { SEGMENT_SLOT_ICON_PX } from "@/shared/timetable/segmentIcons";
import { formatCountdown } from "@/shared/timetable/segmentTime";
import { TourStatusIcon } from "@/shared/timetable/TourStatusIcon";
import {
    resolveSegmentVisualState,
    TOUR_VISUAL_BG,
} from "@/shared/timetable/tourVisualState";
import { useContestElapsed } from "@/shared/timetable/useContestElapsed";
import styles from "./ContestPhaseStrip.module.css";

function PhaseStripContent({ view }: { view: TimetableView }) {
    const elapsed = useContestElapsed(
        view.contest_start_time,
        view.elapsed_seconds ?? 0,
    );
    const phase = derivePhaseDisplay(view, elapsed);
    const segments = view.timeline_segments ?? [];
    const heroVisual = phase.heroSegment
        ? resolveSegmentVisualState(phase.heroSegment)
        : "future";
    const heroBg = phase.heroSegment ? TOUR_VISUAL_BG[heroVisual] : "#e8dfd0";

    if (phase.empty) {
        return (
            <div className={styles.strip} aria-live="polite">
                <CalendarClock size={18} className={styles.emptyIcon} aria-hidden />
                <span className={styles.emptyText}>{phase.heroLabel}</span>
            </div>
        );
    }

    const timerText =
        phase.remainingSeconds != null
            ? formatCountdown(phase.remainingSeconds)
            : null;

    return (
        <div className={styles.strip} aria-live="polite">
            <div className={styles.main}>
                <span
                    className={styles.heroChip}
                    style={{ backgroundColor: heroBg }}
                >
                    {phase.isPause ? (
                        <Coffee
                            size={SEGMENT_SLOT_ICON_PX}
                            strokeWidth={2}
                            className={styles.heroPauseIcon}
                            aria-hidden
                        />
                    ) : phase.heroSegment ? (
                        <span className={styles.heroIcon}>
                            <TourStatusIcon
                                status={phase.heroSegment.status}
                                visualState={heroVisual}
                                size={SEGMENT_SLOT_ICON_PX}
                            />
                        </span>
                    ) : null}
                    <span className={styles.heroLabel}>{phase.heroLabel}</span>
                </span>

                {phase.statusText ? (
                    <span className={styles.statusText}>{phase.statusText}</span>
                ) : null}

                {timerText ? (
                    <span className={styles.timer} aria-label="Оставшееся время">
                        {timerText}
                    </span>
                ) : null}
            </div>

            {segments.length > 0 ? (
                <div className={styles.timeline}>
                    <MiniTimeline segments={segments} elapsed={elapsed} />
                </div>
            ) : null}
        </div>
    );
}

export function ContestPhaseStrip() {
    const { data, isLoading, isError } = useContestTimetable();

    if (isLoading && !data) {
        return (
            <div className={styles.strip} aria-busy="true">
                <span className={styles.loadingText}>Загрузка расписания…</span>
            </div>
        );
    }

    if (isError || !data) {
        return null;
    }

    return <PhaseStripContent view={data} />;
}
