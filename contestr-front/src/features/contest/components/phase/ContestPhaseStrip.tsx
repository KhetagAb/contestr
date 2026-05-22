import { CONTEST_PHASE_FOCUS_PATH, useAppPath } from "@/app/AppPath";
import { useContestTimetable } from "@/features/contest/hooks/useContestTimetable";
import { ContestPhaseDisplay } from "./ContestPhaseDisplay";
import styles from "./ContestPhaseStrip.module.css";

export function ContestPhaseStrip() {
    const { navigate } = useAppPath();
    const { data, isLoading, isError } = useContestTimetable();

    if (isLoading && !data) {
        return (
            <div className={styles.phaseRow}>
                <div className={styles.currentTourCard} aria-busy="true">
                    <span className={styles.loadingText}>Загрузка расписания…</span>
                </div>
            </div>
        );
    }

    if (isError || !data) {
        return null;
    }

    return (
        <ContestPhaseDisplay
            view={data}
            variant="strip"
            onOpenFocusPage={() => navigate(CONTEST_PHASE_FOCUS_PATH)}
        />
    );
}
