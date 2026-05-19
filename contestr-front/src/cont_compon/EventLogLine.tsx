import type { RegattaEvent } from "../client";
import { formatContestTime } from "../utils/eventLog";
import { formatGroupCode } from "../utils/groupCode";
import { EventBonusBadges } from "./EventBonusBadges";
import styles from "./tables.module.css";

type Props = {
    event: RegattaEvent;
};

export function EventLogLine({ event }: Props) {
    const time = formatContestTime(event.time_sec);
    const teamNumber = event.team_number ?? 0;
    const groupCode = formatGroupCode(teamNumber);

    return (
        <div className={styles.eventLogLine}>
            <span className={styles.eventLogTime}>{time}</span>
            <span className={styles.eventLogName}>
                <span
                    className={`${styles.eventChip} ${styles.eventGroupChip}`}
                    title={teamNumber > 0 ? `Группа ${teamNumber}` : undefined}
                >
                    {groupCode}
                </span>
                <span className={`${styles.eventChip} ${styles.eventNameChip}`}>
                    {event.display_name}
                </span>
            </span>
            <span className={styles.eventLogAction}>
                решил задачу{" "}
                <span className={styles.eventChip}>{event.problem_code}</span>
                <EventBonusBadges
                    solved_in_time={event.solved_in_time}
                    first_in_group={event.first_in_group}
                />
            </span>
            <span className={styles.eventLogPoints}>
                <span className={styles.eventChip}>+{event.points}</span>
            </span>
        </div>
    );
}
