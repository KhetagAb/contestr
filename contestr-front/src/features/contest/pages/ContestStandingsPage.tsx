import { useMemo } from "react";
import { useSearchParam } from "react-use";
import {
    createColumnHelper,
    flexRender,
    getCoreRowModel,
    getSortedRowModel,
    useReactTable,
} from "@tanstack/react-table";
import { SlArrowDown, SlArrowUp } from "react-icons/sl";
import type { ProblemResult, RegattaContestRow } from "@/client";
import { ContestEventLog } from "@/features/contest/components/event-log/ContestEventLog";
import { ContestPhaseStrip } from "@/features/contest/components/phase/ContestPhaseStrip";
import { TaskCell } from "@/features/contest/components/standings/TaskCell";
import {
    hasRegattaStartedTours,
    useContestStandings,
} from "@/features/contest/hooks/useContestStandings";
import { useContestParticipants } from "@/features/contest/hooks/useContestParticipants";
import { useFollowedParticipant } from "@/features/contest/follow/FollowedParticipantContext";
import { formatGroupCode } from "@/shared/utils/groupCode";
import ContestHomePage from "@/features/contest/pages/ContestHomePage";
import styles from "./ContestStandingsPage.module.css";

const columnHelper = createColumnHelper<RegattaContestRow>();

function isOnFollowedTeam(
    row: RegattaContestRow,
    followedParticipantId: string | null,
    followedTeamNumber: number | undefined,
): boolean {
    return !!(
        followedParticipantId &&
        followedTeamNumber != null &&
        followedTeamNumber > 0 &&
        row.team_number === followedTeamNumber
    );
}

function groupChipClassName(
    row: RegattaContestRow,
    followedParticipantId: string | null,
    followedTeamNumber: number | undefined,
): string {
    if (isOnFollowedTeam(row, followedParticipantId, followedTeamNumber)) {
        return `${styles.groupChip} ${styles.groupChipFollowed}`;
    }
    return styles.groupChip;
}

/** Фоны групп — тёплые оттенки в духе #fbf8f3 / event-log, но различимые */
const teamColors = [
    "#f9f3e8ee",
    "#efe8f5ee",
    "#e8f3eaee",
    "#f5ebe8ee",
    "#e8eef5ee",
    "#f5f0e8ee",
    "#ebe8e1ee",
    "#f0ebe8ee",
    "#e8f5f0ee",
    "#f5e8f0ee",
    "#f2efe8ee",
    "#ebe5f0ee",
    "#e8f0ebee",
    "#f8f3e8ee",
    "#efeae8ee",
    "#e5eef5ee",
    "#f0e8f2ee",
    "#e8f2f5ee",
    "#f5f2e8ee",
    "#ebe8f5ee",
    "#f3ebe8ee",
    "#e8ebe8ee",
];

function scoreToGreenColor(ratio: number) {
    if (ratio <= 0) return "transparent";
    const t = Math.min(1, Math.max(0, ratio));
    const eased = Math.pow(t, 0.68);
    const r = Math.round(244 + (155 - 244) * eased);
    const g = Math.round(249 + (192 - 249) * eased);
    const b = Math.round(240 + (148 - 240) * eased);
    const alpha = 0.68 + 0.32 * eased;
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

function taskCellBackground(score: number, globalMaxPositive: number) {
    if (score <= 0) {
        return score < 0 ? "rgba(252, 232, 232, 0.5)" : "transparent";
    }
    if (globalMaxPositive <= 0) return "transparent";
    return scoreToGreenColor(score / globalMaxPositive);
}

function globalMaxPositiveScore(rows: RegattaContestRow[] | undefined) {
    if (!rows?.length) return 0;
    let max = 0;
    for (const row of rows) {
        for (const p of row.problem_results) {
            if (p.score > 0) {
                max = Math.max(max, p.score);
            }
        }
    }
    return max;
}

function rowBackgroundColor(row: RegattaContestRow): string {
    if (!row.team_number || row.team_number <= 0) {
        return "#fbf8f3ee";
    }
    return teamColors[(row.team_number - 1) % teamColors.length];
}

function findProblem(row: RegattaContestRow, taskId: string): ProblemResult | undefined {
    return row.problem_results?.find((p) => p.problem_code === taskId);
}

function taskSortValue(row: RegattaContestRow, taskName: string): number {
    return findProblem(row, taskName)?.score ?? 0;
}

export default function ContestStandingsPage() {
    const { followedParticipantId } = useFollowedParticipant();
    const { data, isSuccess, isLoading, isError } = useContestStandings();
    const { participants } = useContestParticipants();
    const contestId = parseInt(useSearchParam("contestId") || "", 10);
    const hasContest = Number.isFinite(contestId) && contestId > 0;

    const hasStartedTours = isSuccess && data ? hasRegattaStartedTours(data) : false;

    const tableRows = useMemo((): RegattaContestRow[] => {
        if (!isSuccess || !data) {
            return [];
        }
        if (hasStartedTours) {
            return data.rows;
        }
        if (data.rows.length > 0) {
            return [...data.rows].sort((a, b) =>
                a.display_name.localeCompare(b.display_name, "ru"),
            );
        }
        return participants.map((p) => ({
            user_id: p.participant_id,
            display_name: p.display_name,
            team_number: 0,
            total_score: 0,
            solved_problems: 0,
            problem_results: [],
        }));
    }, [data, hasStartedTours, isSuccess, participants]);

    const tasks = useMemo(() => {
        if (!hasStartedTours || !tableRows.length) {
            return false;
        }
        const problemResults = tableRows[0].problem_results;
        if (!problemResults?.length) {
            return false;
        }
        return problemResults.map((r) => r.problem_code);
    }, [hasStartedTours, tableRows]);

    const followedTeamNumber = useMemo(
        () =>
            (followedParticipantId &&
                tableRows.find((r) => r.user_id === followedParticipantId)?.team_number) ||
            undefined,
        [tableRows, followedParticipantId],
    );

    const baseColumns = useMemo(
        () => [
            columnHelper.accessor("team_number", {
                cell: (info) => {
                    const n = info.getValue();
                    const row = info.row.original;
                    return (
                        <span
                            className={groupChipClassName(
                                row,
                                followedParticipantId,
                                followedTeamNumber,
                            )}
                            title={`Группа ${n}`}
                        >
                            {formatGroupCode(n)}
                        </span>
                    );
                },
                header: () => "Группа",
            }),
            columnHelper.accessor("display_name", {
                cell: (info) => info.getValue(),
                header: () => "Имя",
            }),
            columnHelper.accessor("total_score", {
                cell: (info) => info.getValue() + "",
                header: () => "Счет",
            }),
        ],
        [followedParticipantId, followedTeamNumber],
    );

    const taskColumns = useMemo(
        () =>
            tasks
                ? tasks.map((taskName) =>
                      columnHelper.accessor((row) => taskSortValue(row, taskName), {
                          id: `task_${taskName}`,
                          header: taskName,
                          cell: (props) => {
                              const problem = findProblem(props.row.original, taskName);
                              return <TaskCell problem={problem} />;
                          },
                      })
                  )
                : [],
        [tasks]
    );

    const table = useReactTable({
        columns: [...baseColumns, ...taskColumns],
        data: tableRows,
        getCoreRowModel: getCoreRowModel(),
        getSortedRowModel: getSortedRowModel(),
    });

    const contestTitle = data?.contest_name?.trim() || "Контест";
    const taskList = tasks || [];

    const maxPositiveScore = useMemo(
        () => globalMaxPositiveScore(tableRows),
        [tableRows]
    );

    if (!hasContest) {
        return <ContestHomePage />;
    }

    if (isLoading && !data) {
        return (
            <section className={styles.standingsSection} aria-label="Таблица результатов">
                <p className={styles.statusText} role="status">
                    Загрузка таблицы…
                </p>
            </section>
        );
    }

    if (isError) {
        return (
            <section className={styles.standingsSection} aria-label="Таблица результатов">
                <p className={styles.statusText} role="alert">
                    Не удалось загрузить таблицу результатов
                </p>
            </section>
        );
    }

    const renderSortableTh = (
        columnId: string,
        label: string,
        rowSpan?: number,
        colSpan?: number,
        extraClassName?: string
    ) => {
        const column = table.getColumn(columnId);
        const sorted = column?.getIsSorted();
        const canSort = column?.getCanSort() ?? false;
        const thClass = [extraClassName, canSort ? styles.canSort : undefined]
            .filter(Boolean)
            .join(" ");
        return (
            <th
                key={columnId}
                rowSpan={rowSpan}
                colSpan={colSpan}
                className={thClass || undefined}
                onClick={canSort ? column?.getToggleSortingHandler() : undefined}
            >
                <div className={styles.headerColumn}>
                    <div>{label}</div>
                    {sorted === "asc" ? (
                        <SlArrowUp size="12px" />
                    ) : sorted === "desc" ? (
                        <SlArrowDown size="12px" />
                    ) : null}
                </div>
            </th>
        );
    };

    return (
        <section className={styles.standingsSection} aria-label={contestTitle}>
            <ContestPhaseStrip />
            <div className={styles.standingsCard}>
                <h2 className={styles.standingsCardTitle}>{contestTitle}</h2>
                <table className={styles.standingsTable}>
                    <thead>
                        {taskList.length > 0 ? (
                            <>
                                <tr className={styles.standingsTableHeader}>
                                    {renderSortableTh("team_number", "Группа", 2)}
                                    {renderSortableTh("display_name", "Имя", 2, undefined, styles.nameColumn)}
                                    {renderSortableTh("total_score", "Счет", 2)}
                                    <th colSpan={taskList.length} className={styles.tasksGroupHeader}>
                                        Задачи
                                    </th>
                                </tr>
                                <tr className={styles.standingsTableHeader}>
                                    {taskList.map((taskName) =>
                                        renderSortableTh(
                                            `task_${taskName}`,
                                            taskName,
                                            undefined,
                                            undefined,
                                            styles.taskColumn
                                        )
                                    )}
                                </tr>
                            </>
                        ) : (
                            <tr className={styles.standingsTableHeader}>
                                {renderSortableTh("team_number", "Группа")}
                                {renderSortableTh("display_name", "Имя", undefined, undefined, styles.nameColumn)}
                                {renderSortableTh("total_score", "Счет")}
                            </tr>
                        )}
                    </thead>
                    <tbody>
                        {table.getRowModel().rows.map((row) => {
                            const teamBg = rowBackgroundColor(row.original);
                            const isFollowedSelf =
                                !!followedParticipantId &&
                                row.original.user_id === followedParticipantId;
                            return (
                                <tr
                                    key={row.id}
                                    className={
                                        isFollowedSelf ? styles.rowFollowedSelf : undefined
                                    }
                                >
                                    {row.getVisibleCells().map((cell) => {
                                        const taskId = cell.column.id.replace("task_", "");
                                        const isTaskColumn = cell.column.id.startsWith("task_");
                                        const isIdentityColumn = [
                                            "team_number",
                                            "display_name",
                                            "total_score",
                                        ].includes(cell.column.id);

                                        const problem = isTaskColumn
                                            ? findProblem(row.original, taskId)
                                            : undefined;
                                        const curScore = problem?.score ?? 0;

                                        const cellClassName = isTaskColumn
                                            ? styles.taskColumn
                                            : cell.column.id === "display_name"
                                              ? styles.nameColumn
                                              : undefined;

                                        return (
                                            <td
                                                key={cell.id}
                                                className={cellClassName}
                                                style={{
                                                    backgroundColor: isIdentityColumn
                                                        ? teamBg
                                                        : taskCellBackground(
                                                              curScore,
                                                              maxPositiveScore,
                                                          ),
                                                }}
                                            >
                                                {flexRender(
                                                    cell.column.columnDef.cell,
                                                    cell.getContext()
                                                )}
                                            </td>
                                        );
                                    })}
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            </div>
            {hasStartedTours && (data?.events?.length ?? 0) > 0 ? (
                <ContestEventLog events={data!.events} />
            ) : null}
        </section>
    );
}
