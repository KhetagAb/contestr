import styles from "./tables.module.css";
import {
    createColumnHelper,
    flexRender,
    getCoreRowModel,
    getSortedRowModel,
    useReactTable,
} from "@tanstack/react-table";
import type { RegattaContestRow, ProblemResult } from "../client";
import {useEffect, useState} from "react";

// Описываем тип строки таблицы (одна посылка)
type RecentParcelRow = RegattaContestRow & ProblemResult;

export const createColumns = (contestStartTime: number) => {
    const columnHelper = createColumnHelper<RecentParcelRow>();

    return [
        columnHelper.accessor("last_submission_time", {
            header: () => "Время последней посылки",
            cell: (info) => {
                const lastSubmission = info.getValue() || 0;
                const elapsed = lastSubmission - contestStartTime;

                if (elapsed < 0) return "00:00:00";

                const hours = Math.floor(elapsed / 3600);
                const mins = Math.floor((elapsed % 3600) / 60);
                const secs = elapsed % 60;

                return `${hours.toString().padStart(2, "0")}:${mins
                    .toString()
                    .padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
            },
            sortingFn: "alphanumeric",
        }),

        columnHelper.accessor("display_name", {
            header: () => "Фамилия и имя",
            cell: (info) => info.getValue(),
        }),

        columnHelper.accessor("team_number", {
            header: () => "№№ команды",
            cell: (info) => info.getValue(),
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

const RecentParcelsTable = ({
                                data,
                                contestStartTime,
                            }: {
    data: RecentParcelRow[];
    contestStartTime: number;
}) => {
    const columns = createColumns(contestStartTime);

    data = data.sort((a, b) => {
        return b.last_submission_time - a.last_submission_time
    })

    // const [sorting, setSorting] = useState<SortingState>([]);
    const table = useReactTable({
        data,
        columns,
        // state: { sorting },
        // onSortingChange: setSorting,
        getCoreRowModel: getCoreRowModel(),
        getSortedRowModel: getSortedRowModel(),
    });


    const [visibleRows, setVisibleRows] = useState<RecentParcelRow[]>([]);
    const [newRowIds, setNewRowIds] = useState<Set<string>>(new Set());

    useEffect(() => {
        const newIds: string[] = [];
        data.forEach((row) => {
            const rowKey = `${row.user_id}-${row.problem_code}`;
            if (!visibleRows.find(r => `${r.user_id}-${r.problem_code}` === rowKey)) {
                newIds.push(rowKey);
            }
        });

        if (newIds.length > 0) {
            setVisibleRows(data);
            setNewRowIds(prev => new Set([...prev, ...newIds]));

            setTimeout(() => {
                setNewRowIds(prev => {
                    const copy = new Set(prev);
                    newIds.forEach(id => copy.delete(id));
                    return copy;
                });
            }, 10000);
        }
    }, [data]);


    return (
        <div>
            <h2>Последние посылки:</h2>
            <table className={styles.table}>
                <tbody>
                {table.getRowModel().rows.map((row) => {
                    const rowKey = `${row.original.user_id}-${row.original.problem_code}`;
                    const isNew = newRowIds.has(rowKey);

                    return (
                        <tr
                            key={row.id}
                            className={isNew ? styles.newRow : ""}
                        >
                            {row.getVisibleCells().map((cell) => {
                                const color =
                                    cell.column.id === "score"
                                        ? cell.getValue() > 0
                                            ? "green"
                                            : cell.getValue() < 0
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
