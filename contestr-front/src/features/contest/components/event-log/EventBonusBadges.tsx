import type { RegattaEvent } from "@/client";
import { MedalFirstIcon } from "@/shared/icons/MedalFirstIcon";
import styles from "./event-log.module.css";

type Props = Pick<RegattaEvent, "first_in_group">;

export function EventBonusBadges({ first_in_group }: Props) {
    if (!first_in_group) {
        return null;
    }

    return (
        <span className={styles.eventBonusBadges}>
            <span
                className={`${styles.eventBonusChip} ${styles.eventBonusChipFirst}`}
                title="Бонус за первое решение среди соперников в группе"
                aria-label="Бонус за обгон в группе"
            >
                <MedalFirstIcon className={styles.eventBonusIcon} />
            </span>
        </span>
    );
}
