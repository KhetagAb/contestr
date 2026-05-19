import styles from "./tables.module.css";
import {
    createColumnHelper,
    flexRender,
    getCoreRowModel,
    getSortedRowModel,
    useReactTable,
} from "@tanstack/react-table";
import type { RegattaContestRow, ProblemResult } from "../client";
import {useMemo} from "react";
import { useCurRegattaData } from "../data";
import { formatGroupCode } from "../utils/groupCode";

// Описываем тип строки таблицы (одна посылка)
type RecentParcelRow = RegattaContestRow & ProblemResult;

const formatSeconds = (totalSeconds: number) => {
    const seconds = Math.max(0, Math.floor(totalSeconds));
    const hours = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    return `${hours.toString().padStart(2, "0")}:${mins
        .toString()
        .padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
};

export const createColumns = (tourDurationMinutes?: number) => {
    const columnHelper = createColumnHelper<RecentParcelRow>();
    const tourDurationSec = (tourDurationMinutes ?? 0) * 60;

    return [
        columnHelper.accessor("last_submission_time", {
            header: () => "Время до конца тура",
            cell: (info) => {
                const lastSubmission = info.getValue();

                if (
                    lastSubmission === undefined ||
                    lastSubmission === null ||
                    isNaN(lastSubmission) ||
                    tourDurationSec <= 0
                ) {
                    return "00:00:00";
                }

                const remaining = Math.max(0, tourDurationSec - lastSubmission);
                return formatSeconds(remaining);
            },
            sortingFn: "alphanumeric",
        }),

        columnHelper.accessor("display_name", {
            header: () => "Фамилия и имя",
            cell: (info) => info.getValue(),
        }),

        columnHelper.accessor("team_number", {
            header: () => "Команда",
            cell: (info) => {
                const n = info.getValue();
                return <span title={`Группа ${n}`}>{formatGroupCode(n)}</span>;
            },
        }),

        columnHelper.accessor("problem_code", {
            header: () => "Задача",
            cell: (info) => info.getValue(),
        }),

        columnHelper.accessor("score", {
            header: () => "Изменения в баллах",
            cell: (info) => {
                const value = info.getValue() || 0;
                return value > 0 ? `+${value}` : value;
            },
        }),
    ];
};

const RecentParcelsTable = () => {
    const {data} = useCurRegattaData();
    const columns = useMemo(
        () => createColumns(data?.current_tour_duration),
        [data?.current_tour_duration],
    );
    const rows = useMemo(() => {
        const contestRows = data?.rows ?? [];
        // Разворачиваем данные: для каждого участника создаем отдельные строки для каждой задачи
        const flattenedRows: RecentParcelRow[] = [];
        for (const row of contestRows) {
            for (const problemResult of row.problem_results) {
                // Показываем только успешные посылки (с last_submission_time)
                if (problemResult.last_submission_time !== undefined && problemResult.last_submission_time !== null) {
                    flattenedRows.push({
                        ...row,
                        ...problemResult,
                    });
                }
            }
        }
        // Сортируем по времени последней посылки (от новых к старым)
        flattenedRows.sort((a, b) => {
            return (b.last_submission_time ?? 0) - (a.last_submission_time ?? 0);
        });
        return flattenedRows;
    }, [data])

    // const [sorting, setSorting] = useState<SortingState>([]);
    const table = useReactTable({
        data: (rows as RecentParcelRow[]), // todo: this is wrong
        columns,
        // state: { sorting },
        // onSortingChange: setSorting,
        getCoreRowModel: getCoreRowModel(),
        getSortedRowModel: getSortedRowModel(),
    });


    // const [visibleRows, setVisibleRows] = useState<RecentParcelRow[]>([]);
    // const [newRowIds, setNewRowIds] = useState<Set<string>>(new Set());
    // TODO: переделать это всё используя средства reactTable: https://tanstack.com/table/latest/docs/guide/column-filtering
    // useEffect(() => {
    //     const newIds: string[] = [];
    //     rows.forEach((row) => {
    //         const rowKey = `${row.user_id}-${row.problem_code}`;
    //         if (!visibleRows.find(r => `${r.user_id}-${r.problem_code}` === rowKey)) {
    //             newIds.push(rowKey);
    //         }
    //     });

    //     if (newIds.length > 0) {
    //         setVisibleRows(rows);
    //         setNewRowIds(prev => new Set([...prev, ...newIds]));

    //         setTimeout(() => {
    //             setNewRowIds(prev => {
    //                 const copy = new Set(prev);
    //                 newIds.forEach(id => {copy.delete(id)});
    //                 return copy;
    //             });
    //         }, 10000);
    //     }
    // }, [rows, visibleRows]);


    return (
        <div>
            <h2>Последние посылки:</h2>
            <table className={styles.table}>
                <tbody>
                {table.getRowModel().rows.map((row) => {
                    // const rowKey = `${row.original.user_id}-${row.original.problem_code}`;
                    // const isNew = newRowIds.has(rowKey);
                    const isNew = false;
                    return (
                        <tr
                            key={row.id}
                            className={isNew ? styles.newRow : ""}
                        >
                            {row.getVisibleCells().map((cell) => {
                                const color =
                                    cell.column.id === "score"
                                        ? (cell.getValue() as number) > 0
                                            ? "green"
                                            : (cell.getValue() as number) < 0
                                                ? "red"
                                                : "inherit"
                                        : "inherit";

                                return (
                                    <td key={cell.id} style={{ color }}>
                                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                    </td>
                                );
                            })}
                        </tr>
                    );
                })}
                </tbody>
            </table>
        </div>
    );
};

export default RecentParcelsTable;
