import type { CSSProperties } from "react";
import type { RegattaEvent } from "@/client";
import { AppointmentMissedIcon } from "@/shared/icons/AppointmentMissedIcon";
import { formatContestTime } from "@/shared/utils/eventLog";
import { formatGroupCode } from "@/shared/utils/groupCode";
import { EventBonusBadges } from "./EventBonusBadges";
import styles from "./event-log.module.css";

type Props = {
    event: RegattaEvent;
    /** Каскад при первом появлении строки в DOM (reload или новое событие с poll). */
    enterDelayMs?: number;
};

export function EventLogLine({ event, enterDelayMs = 0 }: Props) {
    const time = formatContestTime(event.time_sec);
    const teamNumber = event.team_number ?? 0;
    const groupCode = formatGroupCode(teamNumber);
    const isOvertakeEvent = event.type === "problem_overtake";
    const isRejectedEvent = event.type === "problem_rejected";
    const isOutsideTour = event.solved_in_time === false;

    const lineStyle = {
        "--event-enter-delay": `${enterDelayMs}ms`,
    } as CSSProperties;

    return (
        <div
            className={`${styles.eventLogLine} ${styles.eventLogLineEnter}${isRejectedEvent ? ` ${styles.eventLogLineRejected}` : ""}`}
            style={lineStyle}
        >
            <span className={styles.eventLogTimeCell}>
                <span className={styles.eventLogTime}>{time}</span>
                {isOutsideTour && (
                    <span
                        className={styles.eventOutsideTourChip}
                        title="Вне тура"
                        aria-label="Посылка вне тура"
                    >
                        <AppointmentMissedIcon className={styles.eventOutsideTourIcon} />
                    </span>
                )}
            </span>
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
                {isRejectedEvent ? (
                    <>
                        получает{" "}
                        <span className={styles.eventChip}>
                            {event.verdict ?? "?"}
                        </span>{" "}
                        по задаче{" "}
                        <span className={styles.eventChip}>{event.problem_code}</span>
                    </>
                ) : isOvertakeEvent ? (
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
                {!isOvertakeEvent && !isRejectedEvent && (
                    <EventBonusBadges first_in_group={event.first_in_group} />
                )}
            </span>
            <span className={styles.eventLogPoints}>
                {!isRejectedEvent && (
                    <span className={styles.eventChip}>+{event.points}</span>
                )}
            </span>
        </div>
    );
}
