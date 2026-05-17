import styles from "./tables.module.css";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFacetedMinMaxValues,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
// import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
import { SlArrowDown, SlArrowUp } from "react-icons/sl";
import type { ProblemResult, RegattaContestRow } from "../client";
import { useCurRegattaData } from "../data";

const columnHelper = createColumnHelper<RegattaContestRow>();

const hightlighted_user_id = "0";

const columns = [
  columnHelper.accessor("team_number", {
    cell: (info) => info.getValue(),
    header: () => "Команда",
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

const teamColors = [
  "#cdd7d9dd",
  "#e1e4f2dd",
  "#c9cce2dd",
  "#F5CBCBdd",
  "#f5dce5dd",
  "#FFEAEAdd",
  "#D4E6F1dd",
  "#E8F5E9dd",
  "#f9f3e8dd",
  "#FCE4ECdd",
  "#E0F7FAdd",
  "#F1F8E9dd",
  "#f9eededd",
  "#D7CCC8dd",
  "#e5d0e8dd",
  "#ced1e3dd",
  "#c9dedcdd",
  "#f8dfd7dd",
  "#e1ead6dd",
  "#fad8e4dd",
  "#d2dce5dd",
  "#CFD8DCdd",
];

function scoreToGreenColor(score: number) {
  if (score <= 0) return "transparent";
  const clamped = Math.min(1, Math.max(0, score));
  const alpha = clamped;
  const g = Math.floor(255 * clamped);
  return `rgba(0, ${g}, 0, ${alpha})`;
}

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins.toString().padStart(2, "0")}:${secs
      .toString()
      .padStart(2, "0")}`;
};

function colorTeam(
    row: RegattaContestRow,
    hightlighted_user_team_number: number | undefined
) {
  if (row.user_id === hightlighted_user_id) {
    return {
      backgroundColor: "#6466fdff",
      border: "2px solid #7a76e1ff",
    };
  } else if (
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

const ResultsTable = () => {
  // когда нужен сервер ьудет это раскоментить
  const { data, isSuccess } = useCurRegattaData()

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
              data.rows.find((r) => r.user_id === hightlighted_user_id)
                  ?.team_number) ||
          undefined,
      [data, isSuccess]
  );

  const table = useReactTable({
    columns: [
      ...columns,
      ...(tasks
          ? [
            columnHelper.group({
              id: "tasks",
              header: () => "Задачи",
              columns: tasks.map((taskName, i) =>
                  columnHelper.accessor(`problem_results`, {
                    id: `task_${taskName}`,
                    header: taskName,
                    cell: (props) => {
                      const problem = props.row.original.problem_results[i];
                      return (
                          <div>
                            {problem ? problem.score : "—"}
                            {problem?.last_submission_time && (
                                <div className={styles.taskTime}>
                                  {formatTime(problem.last_submission_time)}
                                </div>
                            )}
                          </div>
                      );
                    },
                  })
              ),
            }),
          ]
          : []),
    ],
    data: data?.rows || [],
    debugTable: true,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getFacetedMinMaxValues: getFacetedMinMaxValues(),
  });

  return (
      <div>
        <h2>
          Таблица результатов
        </h2>
        <table className={styles.table}>
          <thead>
          {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id} className={styles.tableHeader}>
                {headerGroup.headers.map((header) => {
                  const columnRelativeDepth = header.depth - header.column.depth;
                  if (columnRelativeDepth > 1) return null;
                  let rowSpan = 1;
                  if (header.isPlaceholder) {
                    const leafs = header.getLeafHeaders();
                    rowSpan = leafs[leafs.length - 1].depth - header.depth;
                  }
                  return (
                      <th
                          key={header.id}
                          colSpan={header.colSpan}
                          rowSpan={2 - rowSpan}
                          className={
                            header.column.getCanSort() ? styles.canSort : ""
                          }
                          onClick={header.column.getToggleSortingHandler()}
                      >
                        <div className={styles.headerColumn}>
                          <div>
                            {flexRender(
                                header.column.columnDef.header,
                                header.getContext()
                            )}
                          </div>
                          {{
                            asc: <SlArrowUp size="12px" />,
                            desc: <SlArrowDown size="12px" />,
                          }[header.column.getIsSorted() as string] ?? null}
                        </div>
                      </th>
                  );
                })}
              </tr>
          ))}
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
                    const columnValues =
                        table
                            .getColumn(cell.column.id)
                            ?.getFacetedRowModel()
                            .rows.map((row) => {
                          const problemResults = row.original.problem_results;
                          const problem = problemResults?.find(
                              (p) => p.problem_code === taskId
                          );
                          return problem?.score ?? 0;
                        }) ?? [0];

                    const minmaxValues = columnValues
                        ? Math.max(...columnValues)
                        : 0;

                    const ddd = cell.getValue() as
                        | Array<ProblemResult>
                        | undefined;
                    const curScore =
                        (ddd?.find?.((p) => p.problem_code === taskId)?.score) ||
                        0;
                    return (
                        <td
                            key={cell.id}
                            style={{
                              backgroundColor:
                                  ["team_number", "display_name", "user_id"].includes(
                                      cell.column.id
                                  )
                                      ? rowHightlight.backgroundColor
                                      : scoreToGreenColor(curScore / minmaxValues),
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
        <div />
      </div>
  );
};

export default ResultsTable;
