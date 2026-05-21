import { useMemo } from "react";
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
import { useContestStandings } from "@/features/contest/hooks/useContestStandings";
import { formatGroupCode } from "@/shared/utils/groupCode";
import styles from "./ContestStandingsPage.module.css";

const columnHelper = createColumnHelper<RegattaContestRow>();

const hightlighted_user_id = "0";

const columns = [
    columnHelper.accessor("team_number", {
        cell: (info) => {
            const n = info.getValue();
            return (
                <span className={styles.groupChip} title={`Группа ${n}`}>
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
];

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

function colorTeam(
    row: RegattaContestRow,
    hightlighted_user_team_number: number | undefined
) {
    if (row.user_id === hightlighted_user_id) {
        return {
            backgroundColor: "#6466fdff",
            border: "2px solid #7a76e1ff",
        };
    }
    if (
        hightlighted_user_team_number &&
        row.team_number === hightlighted_user_team_number
    ) {
        return {
            backgroundColor: "#6466fdff",
        };
    }
    return {
        backgroundColor: teamColors[(row.team_number - 1) % teamColors.length],
    };
}

function findProblem(row: RegattaContestRow, taskId: string): ProblemResult | undefined {
    return row.problem_results?.find((p) => p.problem_code === taskId);
}

function taskSortValue(row: RegattaContestRow, taskName: string): number {
    return findProblem(row, taskName)?.score ?? 0;
}

export default function ContestStandingsPage() {
    const { data, isSuccess } = useContestStandings();

    const tasks = useMemo(() => {
        return (
            isSuccess &&
            data.rows &&
            data.rows.length > 0 &&
            data.rows[0].problem_results.map((r) => r.problem_code)
        );
    }, [data, isSuccess]);

    const hightlighted_user_team_number = useMemo(
        () =>
            (isSuccess &&
                data.rows &&
                data.rows.find((r) => r.user_id === hightlighted_user_id)?.team_number) ||
            undefined,
        [data, isSuccess]
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
        columns: [...columns, ...taskColumns],
        data: data?.rows || [],
        getCoreRowModel: getCoreRowModel(),
        getSortedRowModel: getSortedRowModel(),
    });

    const contestTitle = data?.contest_name?.trim() || "Контест";
    const taskList = tasks || [];

    const maxPositiveScore = useMemo(
        () => globalMaxPositiveScore(data?.rows),
        [data?.rows]
    );

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
            <div className={styles.phaseCard}>
                <ContestPhaseStrip />
            </div>
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
                            const rowHightlight = colorTeam(
                                row.original,
                                hightlighted_user_team_number
                            );
                            return (
                                <tr key={row.id}>
                                    {row.getVisibleCells().map((cell) => {
                                        const taskId = cell.column.id.replace("task_", "");
                                        const isTaskColumn = cell.column.id.startsWith("task_");

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
                                                    backgroundColor: [
                                                        "team_number",
                                                        "display_name",
                                                        "user_id",
                                                    ].includes(cell.column.id)
                                                        ? rowHightlight.backgroundColor
                                                        : taskCellBackground(
                                                              curScore,
                                                              maxPositiveScore
                                                          ),
                                                    borderTop: rowHightlight.border,
                                                    borderBottom: rowHightlight.border,
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
            {isSuccess && data?.events && <ContestEventLog events={data.events} />}
        </section>
    );
}
