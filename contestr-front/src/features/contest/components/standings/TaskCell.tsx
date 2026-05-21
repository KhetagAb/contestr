import type { ProblemResult } from "@/client";
import styles from "@/features/contest/pages/ContestStandingsPage.module.css";
import { getTaskCellKind } from "./taskCellKind";

const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
};

type Props = {
    problem: ProblemResult | undefined;
};

export function TaskCell({ problem }: Props) {
    const cell = getTaskCellKind(problem);

    if (cell.kind === "empty" || cell.kind === "none") {
        return <span className={styles.taskEmpty}>—</span>;
    }

    if (cell.kind === "failed") {
        return <span className={styles.taskFailed}>−{cell.count}</span>;
    }

    const { problem: solved } = cell;
    return (
        <div className={styles.taskSolved}>
            <span className={styles.taskScore}>{solved.score}</span>
            {solved.last_submission_time != null && (
                <div className={styles.taskTime}>{formatTime(solved.last_submission_time)}</div>
            )}
        </div>
    );
}
