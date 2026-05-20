import { useContests } from "@/shared/hooks/useContests";
import { ContestTimetableBlock } from "./ContestTimetableBlock";
import { useTimetableHiddenContests } from "./useTimetableHiddenContests";
import "./TimetablePage.css";

export default function TimetablePage() {
    const { contests, isLoading } = useContests();
    const { isHidden, toggleHidden } = useTimetableHiddenContests();

    return (
        <section className="admin-timetable tt-page">
            {isLoading && contests.length === 0 && (
                <p className="tt-message tt-message--info">Загрузка контестов…</p>
            )}

            {!isLoading && contests.length === 0 && (
                <p className="tt-message tt-message--info">
                    Контесты не настроены. Добавьте их в разделе «Контесты».
                </p>
            )}

            {contests.map((contest, index) => (
                <div key={contest.contest_id} className="tt-contest-block-wrap">
                    {index > 0 && <hr className="tt-contest-divider" aria-hidden />}
                    <ContestTimetableBlock
                        contestId={contest.contest_id}
                        contestName={contest.name}
                        collapsed={isHidden(contest.contest_id)}
                        onToggleCollapsed={() => toggleHidden(contest.contest_id)}
                    />
                </div>
            ))}
        </section>
    );
}
