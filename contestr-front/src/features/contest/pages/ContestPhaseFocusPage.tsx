import { ArrowLeft } from "lucide-react";
import { useSearchParam } from "react-use";
import { useAppPath } from "@/app/AppPath";
import { ContestPhaseDisplay } from "@/features/contest/components/phase/ContestPhaseDisplay";
import { useContestStandings } from "@/features/contest/hooks/useContestStandings";
import { useContestTimetable } from "@/features/contest/hooks/useContestTimetable";
import styles from "./ContestPhaseFocusPage.module.css";

export default function ContestPhaseFocusPage() {
    const { navigate } = useAppPath();
    const contestId = parseInt(useSearchParam("contestId") || "", 10);
    const hasContest = Number.isFinite(contestId) && contestId > 0;
    const { data: standings } = useContestStandings();
    const { data, isLoading, isError, isFetching } = useContestTimetable();
    const contestTitle = standings?.contest_name?.trim() || "Контест";

    return (
        <section className={styles.focusPage} aria-label="Расписание тура">
            <button
                type="button"
                className={styles.backBtn}
                onClick={() => navigate("/")}
            >
                <ArrowLeft size={18} aria-hidden />
                К таблице
            </button>

            <div className={styles.focusCenter}>
                <header className={styles.focusHeader}>
                    {hasContest ? (
                        <h1 className={styles.contestTitle}>{contestTitle}</h1>
                    ) : null}
                </header>

                <div className={styles.focusMain}>
                    {!hasContest ? (
                        <p className={styles.statusText}>
                            Не указан контест. Откройте страницу из таблицы результатов.
                        </p>
                    ) : null}

                    {hasContest && isLoading && !data ? (
                        <p className={styles.statusText}>Загрузка расписания…</p>
                    ) : null}

                    {hasContest && isError ? (
                        <p className={styles.statusText}>Не удалось загрузить расписание</p>
                    ) : null}

                    {hasContest && !isLoading && !isFetching && !isError && !data ? (
                        <p className={styles.statusText}>Расписание недоступно</p>
                    ) : null}

                    {data ? <ContestPhaseDisplay view={data} variant="focus" /> : null}
                </div>
            </div>
        </section>
    );
}
