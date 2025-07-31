import { ReactTabulator } from 'react-tabulator';
import { useEffect, useState } from 'react';
import './tables.css';
// import 'react-tabulator/css/tabulator.min.css';

 const localData = [
    { id: 1, name: "Иван Петров", "total-score": 145, resolved: 12, "*A": 8, "*B": 4 },
    { id: 2, name: "Елена Смирнова", "total-score": 210, resolved: 18, "*A": 12, "*B": 6 },
    { id: 3, name: "Алексей Козлов", "total-score": 98, resolved: 9, "*A": 5, "*B": 4 },
    { id: 4, name: "Ольга Новикова", "total-score": 176, resolved: 15, "*A": 10, "*B": 5 },
    { id: 5, name: "Дмитрий Васнецов", "total-score": 231, resolved: 20, "*A": 14, "*B": 6 },
    { id: 6, name: "Анна Кузнецова", "total-score": 120, resolved: 11, "*A": 7, "*B": 4 },
    { id: 7, name: "Михаил Белов", "total-score": 187, resolved: 16, "*A": 11, "*B": 5 },
    { id: 8, name: "Татьяна Морозова", "total-score": 165, resolved: 14, "*A": 9, "*B": 5 },
    { id: 9, name: "Сергей Иванов", "total-score": 199, resolved: 17, "*A": 12, "*B": 5 },
    { id: 10, name: "Наталья Соколова", "total-score": 142, resolved: 13, "*A": 8, "*B": 5 },
];

const columns = [
    {
        title: "ID",
        field: "id",
        width: 70,
        frozen: true,
        hozAlign: "center"
    },
    {
        title: "Имя",
        field: "name",
        width: 150,
        frozen: true,
        // headerFilter: "input"
    },
    {
        title: "Сумма баллов",
        field: "total-score",
        width: 120,
        hozAlign: "right",
        sorter: "number",
        formatter: "progress",
        formatterParams: {
            min: 0,
            max: 250,
            color: ["#e8e5c2", "#d8bdc1", "#a4d1a4"],
            legend: "total-score",
        }
    },
    {
        title: "Дорешено",
        field: "resolved",
        width: 100,
        sorter: "number",
        hozAlign: "center"
    },
    {
        title: "*A",
        field: "*A",
        width: 70,
        sorter: "number",
        hozAlign: "center",
        headerTooltip: "Решено задач типа A"
    },
    {
        title: "*B",
        field: "*B",
        width: 70,
        sorter: "number",
        hozAlign: "center",
        headerTooltip: "Решено задач типа B"
    }
];

const tableOptions = {
    layout: "fitColumns",
    responsiveLayout: "collapse",
    pagination: "local",
    paginationSize: 10,
    paginationSizeSelector: [5, 10, 20, 50],
    movableColumns: true,
    headerFilterPlaceholder: "Фильтр...",
    initialSort: [
        {column: "total-score", dir: "desc"}
    ],
    locale: "ru",
    langs: {
        "ru": {
            // "pagination": {
            //     "first": "Первая",
            //     "last": "Последняя",
            //     "prev": "Предыдущая",
            //     "next": "Следующая",
            // }
        }
    }
};

const ResultsTable = () => {
    const [data, setData] = useState([]);

    useEffect(() => {
        setData(localData); // Просто устанавливаем данные
    }, []);

    return (
        <div className="results-container">
            <h2 className="results-title">Результаты группы В</h2>


                <ReactTabulator
                    data={data}
                    columns={columns}
                    options={tableOptions}
                    className="results-table"
                    tooltips={true}
                />

        </div>
    );
};

export default ResultsTable;