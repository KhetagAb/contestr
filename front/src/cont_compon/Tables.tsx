import styles from "./tables.module.css";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFacetedMinMaxValues,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
  // type RowData,
} from "@tanstack/react-table";
import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
// import { Filter } from "./filter";
// import 'react-tabulator/css/tabulator.min.css';
import { SlArrowDown, SlArrowUp } from "react-icons/sl";
import type { ProblemResult, RegattaContestRow } from "../client";
import Color from "colorjs.io";
import { getRegattaContestStandingsOptions } from "../client/@tanstack/react-query.gen";
import { useSearchParam } from "react-use";

const columnHelper = createColumnHelper<RegattaContestRow>();

const hightlighted_user_id = 20773;

const columns = [
  columnHelper.accessor("user_id", {
    cell: (info) => info.getValue(),
    header: () => "id",
  }),
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
  columnHelper.accessor("solved_problems", {
    cell: (info) => info.getValue() + "",
    header: () => "Решено",
  }),
];

const teamColors = [
  "#cdd7d9dd", // оригинальный
  "#e1e4f2dd", // оригинальный
  "#c9cce2dd", // оригинальный
  "#F5CBCBdd", // оригинальный
  "#f5dce5dd", // оригинальный
  "#FFEAEAdd", // оригинальный
  "#D4E6F1dd", // светло-голубой
  "#E8F5E9dd", // мятно-зеленый
  "#FFF3E0dd", // теплый бежевый
  "#FCE4ECdd", // розоватый
  "#E0F7FAdd", // бирюзовый
  "#F1F8E9dd", // салатовый
  "#f9eededd", // персиковый
  "#D7CCC8dd", // светло-коричневый
  "#E1BEE7dd", // лавандовый
  "#C5CAE9dd", // сиреневый
  "#B2DFDBdd", // морской волны
  "#f8dfd7dd", // коралловый
  "#DCEDC8dd", // светло-зеленый
  "#fad8e4dd", // розовый
  "#BBDEFBdd", // васильковый
  "#CFD8DCdd", // серо-голубой
];
function scoreToGreenColor(score: number) {
  if (score <= 0) return "transparent";
  const clamped = Math.min(1, Math.max(0, score));
  const vividGreen = new Color("lch(75% 100 136)");
  vividGreen.lch.c *= clamped;
  vividGreen.alpha = clamped;
  return vividGreen.to("srgb").toString({ format: "rgba" });
}

function colorTeam(
  row: RegattaContestRow,
  hightlighted_user_team_number: number | undefined
) {
  if (row.user_id === hightlighted_user_id) {
    return {
      backgroundColor: "#f25577ff",
      border: "2px solid #eba6b3ff",
    };
  } else if (
    hightlighted_user_team_number &&
    row.team_number == hightlighted_user_team_number
  ) {
    return {
      backgroundColor: "#f75e7fff",
      // border: '1px solid #ee9900',
    };
  }
  return {
    backgroundColor: teamColors[(row.team_number - 1) % teamColors.length],
  };
}

// const getScores = () => {
//   return Object.fromEntries([..."ABCDFGHIJKLMNOPQRSTUV"].map((name) => [name, Math.round(Math.random() * 100)]))
// }

const ResultsTable = () => {
  // const [pagination, setPagination] = React.useState<PaginationState>({
  //   pageIndex: 0,
  //   pageSize: 100,
  // });
  const contestId = parseInt(useSearchParam("contestId") || "");

  const { data, isSuccess } = useQuery({
    ...getRegattaContestStandingsOptions({
      path: {
        contest_id: contestId,
      },
    }),
  });
  const tasks = useMemo(() => {
    return isSuccess &&
      data.rows &&
      data.rows.length > 0 &&
      data.rows[0].problem_results
      ? data.rows[0].problem_results.map((r) => r.problem_code)
      : undefined;
  }, [data, isSuccess]); // TODO: fixme.
  const hightlighted_user_team_number = useMemo(
    () =>
      (isSuccess &&
        data.rows &&
        data.rows.find((r) => r.user_id === hightlighted_user_id)
          ?.team_number) ||
      undefined,
    [data, isSuccess]
  ); // TODO: fixme.

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
                  // Доступ ко всему массиву
                  id: `task_${taskName}`, // Уникальный id для колонки
                  header: taskName,
                  cell: (props) => {
                    const problem = props.row.original.problem_results[i];
                    return (
                      <div>
                        {problem ? problem.score : "—"}
                        {problem?.last_submission_time && (
                          <>
                            <div className={styles.taskTime}>
                              {problem.last_submission_time}
                            </div>
                          </>
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
    getFacetedMinMaxValues: getFacetedMinMaxValues(), //if you need min/max values

    // getPaginationRowModel: getPaginationRowModel(),
    // onPaginationChange: setPagination,
    //no need to pass pageCount or rowCount with client-side pagination as it is calculated automatically
    // state: {
    //   pagination,
    // },
    // autoResetPageIndex: false, // turn off page index reset when sorting or filtering
  });

  return (
    <div>
      <h2>Таблица результатов </h2>
      <table className={styles.table}>
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id} className={styles.tableHeader}>
              {headerGroup.headers.map((header) => {
                const columnRelativeDepth = header.depth - header.column.depth;

                if (columnRelativeDepth > 1) {
                  return null;
                }

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
                    {...{
                      className: header.column.getCanSort()
                        ? styles.canSort
                        : "",
                      onClick: header.column.getToggleSortingHandler(),
                    }}
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
                    {/* {header.column.getCanFilter() ? (
                        <div>
                          <Filter column={header.column} table={table} />
                        </div>
                      ) : null} */}
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
                  const taskId = cell.column.id.replace("task_", "") // wtf
                  const columnValues = table
                    .getColumn(cell.column.id)
                    ?.getFacetedRowModel()
                    .rows.map(row => {
                      const problemResults = row.original.problem_results;
                      const problem = problemResults.find(p => p.problem_code === taskId);
                      return problem?.score ?? 0;
                    }) ?? [0];

                    const minmaxValues = columnValues ? Math.max(...columnValues) : 0;


                  const ddd = cell.getValue() as Array<ProblemResult> | undefined;
                  const curScore = ddd && ddd.find && ddd.find(p => p.problem_code === taskId)?.score || 0;
                  return (
                    <td
                      key={cell.id}
                      style={{
                        backgroundColor:
                          ["team_number", "display_name"].indexOf(
                            cell.column.id
                          ) !== -1
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
      {/* <div >
        <button
          onClick={() => table.firstPage()}
          disabled={!table.getCanPreviousPage()}
        >
          {"<<"}
        </button>
        <button
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          {<SlArrowLeftCircle/>}
        </button>
        <button
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          {<SlArrowRightCircle/>}
        </button>
        <button
          onClick={() => table.lastPage()}
          disabled={!table.getCanNextPage()}
        >
          {">>"}
        </button>
        <span>
          <div>Page</div>
          <strong>
            {table.getState().pagination.pageIndex + 1} of{" "}
            {table.getPageCount().toLocaleString()}
          </strong>
        </span>
        <span>
          | Go to page:
          <input 
            type="number"
            min="1"
            max={table.getPageCount()}
            defaultValue={table.getState().pagination.pageIndex + 1}
            onChange={(e) => {
              const page = e.target.value ? Number(e.target.value) - 1 : 0;
              table.setPageIndex(page);
            }}
        
          />
        </span>
        <select
        className="select-page"
          value={table.getState().pagination.pageSize}
          onChange={(e) => {
            table.setPageSize(Number(e.target.value));
          }}
        >
          {[10, 20, 30, 40, 50].map((pageSize) => (
            <option key={pageSize} value={pageSize}>
              Show {pageSize}
            </option>
          ))}
        </select>
      </div>
      <div>
        Showing {table.getRowModel().rows.length.toLocaleString()} of{" "}
        {table.getRowCount().toLocaleString()} Rows
      </div> */}
    </div>
  );
};

export default ResultsTable;
