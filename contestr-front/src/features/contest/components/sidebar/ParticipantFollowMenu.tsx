import { useMemo, useState } from "react";
import {
    useFollowedParticipant,
} from "@/features/contest/follow/FollowedParticipantContext";
import { useContestParticipants } from "@/features/contest/hooks/useContestParticipants";
import type { AdminSidebarSession } from "./Sidebar";
import styles from "./ParticipantFollowMenu.module.css";

type Props = {
    adminSession?: AdminSidebarSession | null;
};

export function ParticipantFollowMenu({ adminSession }: Props) {
    const {
        contestId,
        followedParticipantId,
        setFollowedParticipantId,
        clearFollowedParticipant,
        closeParticipantPicker,
    } = useFollowedParticipant();

    const selectParticipant = (participantId: string) => {
        setFollowedParticipantId(participantId);
        closeParticipantPicker();
    };
    const { participants, isLoading } = useContestParticipants();
    const [filter, setFilter] = useState("");

    const filtered = useMemo(() => {
        const q = filter.trim().toLowerCase();
        if (!q) {
            return participants;
        }
        return participants.filter(
            (p) =>
                p.display_name.toLowerCase().includes(q) ||
                p.participant_id.toLowerCase().includes(q),
        );
    }, [participants, filter]);

    return (
        <div className={styles.popover}>
            <p className={styles.title}>Кто вы?</p>
            {contestId == null ? (
                <p className={styles.hint}>Откройте контест из списка слева</p>
            ) : isLoading && participants.length === 0 ? (
                <p className={styles.hint}>Загрузка участников…</p>
            ) : participants.length === 0 ? (
                <p className={styles.hint}>Участники не настроены</p>
            ) : (
                <>
                    <input
                        type="search"
                        className={styles.search}
                        placeholder="Поиск по имени"
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        aria-label="Поиск участника"
                    />
                    <ul className={styles.list} role="listbox" aria-label="Участники">
                        {filtered.map((p) => {
                            const selected = followedParticipantId === p.participant_id;
                            return (
                                <li key={p.participant_id}>
                                    <button
                                        type="button"
                                        role="option"
                                        aria-selected={selected}
                                        className={
                                            selected
                                                ? `${styles.item} ${styles.itemSelected}`
                                                : styles.item
                                        }
                                        onClick={() => selectParticipant(p.participant_id)}
                                    >
                                        <span className={styles.itemName}>{p.display_name}</span>
                                        <span className={styles.itemHandle}>
                                            {p.participant_id}
                                        </span>
                                    </button>
                                </li>
                            );
                        })}
                    </ul>
                    {filtered.length === 0 ? (
                        <p className={styles.hint}>Никого не найдено</p>
                    ) : null}
                </>
            )}
            {followedParticipantId ? (
                <button
                    type="button"
                    className={styles.clearBtn}
                    onClick={clearFollowedParticipant}
                >
                    Сбросить
                </button>
            ) : null}

            {adminSession ? (
                <>
                    <hr className={styles.divider} />
                    <p className={styles.adminGreeting}>
                        Вы вошли как{" "}
                        <span className={styles.adminName}>{adminSession.username}</span>
                    </p>
                    <a href="/admin" className={styles.adminLink}>
                        Админ-панель
                    </a>
                    <button
                        type="button"
                        className={styles.adminLogout}
                        onClick={adminSession.onLogout}
                    >
                        Выйти
                    </button>
                </>
            ) : null}
        </div>
    );
}
