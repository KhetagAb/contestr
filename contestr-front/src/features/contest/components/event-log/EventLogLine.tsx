import type { RegattaEvent } from "@/client";
import { formatContestTime } from "@/shared/utils/eventLog";
import { formatGroupCode } from "@/shared/utils/groupCode";
import { EventBonusBadges } from "./EventBonusBadges";
import styles from "./event-log.module.css";

type Props = {
    event: RegattaEvent;
};

export function EventLogLine({ event }: Props) {
    const time = formatContestTime(event.time_sec);
    const teamNumber = event.team_number ?? 0;
    const groupCode = formatGroupCode(teamNumber);
    const isOvertakeEvent = event.type === "problem_overtake";

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
                {isOvertakeEvent ? (
                    <>
                        получил бонус за обгон в задаче{" "}
                        <span className={styles.eventChip}>{event.problem_code}</span>
                    </>
                ) : (
                    <>
                        решил задачу{" "}
                        <span className={styles.eventChip}>{event.problem_code}</span>
                    </>
                )}
            </span>
            <span className={styles.eventLogBadgesCell}>
                {!isOvertakeEvent && (
                    <EventBonusBadges first_in_group={event.first_in_group} />
                )}
            </span>
            <span className={styles.eventLogPoints}>
                <span className={styles.eventChip}>+{event.points}</span>
            </span>
        </div>
    );
}
