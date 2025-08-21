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
import {RegattaContestRow} from "../client";

const columnHelper = createColumnHelper<RegattaContestRow>(); // тут будет обновление
// время|задача|имя|баллы|  --  time|problem|name|score|

const colums = [
    columnHelper.accessor("time", {
        cell : (info ) => info.getValue,
        header: () => "Время изменения",
    }),
    columnHelper.accessor("problem",{
        cell: (info) => info.getValue,
        header: () => "Задача",
    }),
    columnHelper.accessor("name", {
        cell: (info) => info.getValue,
        header: () => "Фамилия Имя",
    }),
    columnHelper.accessor("score", {
        cell: (info) => info.getValue,
        header: () => "Изменения в баллах",
    }),
];

const recentParcelsTable = () => {
    const table = useReactTable({
        columns = [
            ...colums,

        ]
    })


    debugTable: true,
        getCoreRowModel: getCoreRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        getFacetedMinMaxValues: getFacetedMinMaxValues(),
}



