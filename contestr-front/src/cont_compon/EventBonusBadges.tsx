import type { RegattaEvent } from "../client";
import { SlClock, SlStar } from "react-icons/sl";
import styles from "./tables.module.css";

type Props = Pick<RegattaEvent, "solved_in_time" | "first_in_group">;

export function EventBonusBadges({ solved_in_time, first_in_group }: Props) {
    if (!solved_in_time && !first_in_group) {
        return null;
    }

    return (
        <span className={styles.eventBonusBadges}>
            {solved_in_time && (
                <span
                    className={`${styles.eventBonusChip} ${styles.eventBonusChipInTime}`}
                    title="Бонус за решение в пределах тура"
                    aria-label="Бонус: во время тура"
                >
                    <SlClock className={styles.eventBonusIcon} aria-hidden />
                    <span>во время тура</span>
                </span>
            )}
            {first_in_group && (
                <span
                    className={`${styles.eventBonusChip} ${styles.eventBonusChipFirst}`}
                    title="Бонус за первое решение среди соперников в группе"
                    aria-label="Бонус: обогнав соперников"
                >
                    <SlStar className={styles.eventBonusIcon} aria-hidden />
                    <span>обогнав соперников</span>
                </span>
            )}
        </span>
    );
}
