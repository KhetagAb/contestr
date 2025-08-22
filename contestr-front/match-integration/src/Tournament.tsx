import styles from "./tour.module.css"
import {useState} from 'react'
import {ChevronDown, ChevronUp} from 'lucide-react'
import {PiStarFill} from "react-icons/pi";
import telegramLogo from '../public/telegramLogo.svg';

// import {
//     createColumnHelper,
//     getCoreRowModel,
//     useReactTable,
// } from '@tanstack/react-table'


const Tournament = () => {
    const [openId, setOpenId] = useState(null);

    // @ts-ignore
    const toggleContainer = (id) => {
        setOpenId(openId === id ? null : id);
    };

    type Teams = {
        id: number
        name: string
        commander: string
        members: string[]
    }


    const defaultData: Teams[] = [
        {
            id: 1,
            name: "Скорпионы",
            commander: "Игорь Васильев",
            members: [
                "Иван Петров", "Алексей Смирнов", "Сергей Кузнецов", "Дмитрий Попов",
                "Михаил Васильев", "Андрей Соколов", "Евгений Михайлов", "Кирилл Новиков",
                "Станислав Фёдоров", "Артём Морозов", "Владислав Волков", "Роман Алексеев",
                "Никита Лебедев", "Глеб Семёнов", "Тимофей Козлов", "Павел Егоров",
                "Константин Павлов", "Виктор Козлов", "Арсений Орлов", "Матвей Зайцев"
            ]
        },
        {
            id: 2,
            name: "Молот",
            commander: "Алексей Кожевников",
            members: [
                "Александр Иванов", "Максим Борисов", "Илья Захаров", "Даниил Яковлев",
                "Григорий Тимофеев", "Руслан Романов", "Олег Виноградов", "Юрий Давыдов",
                "Денис Макаров", "Георгий Шестаков", "Василий Блинов", "Иннокентий Крылов",
                "Платон Тихонов", "Семён Прохоров", "Фёдор Назаров", "Эдуард Самойлов"
            ]
        },
        {
            id: 3,
            name: "Витязь",
            commander: "Артем Семенов",
            members: [
                "Артём Николаев", "Вячеслав Степанов", "Ярослав Гусев", "Антон Тарасов",
                "Леонид Богданов", "Пётр Сорокин", "Филипп Рыбаков", "Демьян Филиппов",
                "Виталий Беляев", "Аркадий Котов", "Савелий Баранов", "Герман Абрамов",
                "Давид Медведев", "Игнат Ширяев", "Марк Галкин", "Лев Кошелев",
                "Савва Терентьев", "Николай Устинов", "Егор Игнатьев", "Захар Туров"
            ]
        },
        {
            id: 4,
            name: "Торнадо",
            commander: "Денис Орлов",
            members: [
                "Денис Калинин", "Тимофей Исаев", "Владлен Чернов", "Родион Быков",
                "Святослав Маслов", "Артур Еремеев", "Геннадий Денисов", "Валентин Жданов",
                "Анатолий Суханов", "Стефан Капустин", "Марат Лапин", "Дамир Лаврентьев"
            ]
        },
        {
            id: 5,
            name: "Феникс",
            commander: "Максим Белов",
            members: [
                "Матвей Афанасьев", "Данила Воронов", "Елисей Зимин", "Роберт Трофимов",
                "Семен Карпов", "Гарри Уваров", "Данила Шарапов", "Арсений Мельников",
                "Демид Соломин", "Данила Громов", "Иван Князев", "Марк Потапов",
                "Леон Савин", "Дмитрий Красильников", "Александр Носков", "Мирон Белозёров",
                "Даниэль Пестов", "Адам Одинцов", "Марк Логинов", "Лев Зыков"
            ]
        }
    ];

    // const columnHelper = createColumnHelper<Teams>()
    //
    // const columns = [
    //     columnHelper.accessor('id', {
    //         header: () => "№№",
    //         cell: info => {
    //             const row = info.row.original
    //             return (
    //                 <span>
    //                 {row.id} — {row.name}
    //             </span>
    //             )
    //         },
    //         footer: info => info.column.id,
    //     }),
    //     columnHelper.accessor('members', {
    //         header: () => <span>Участники</span>,
    //         cell: info => {
    //             const members = info.getValue()
    //             return (
    //                 <div>
    //                     {members.map((m, idx) => (
    //                         <div key={idx}>{m}</div>
    //                     ))}
    //                 </div>
    //             )
    //         },
    //         footer: info => info.column.id,
    //     }),
    // ]
    //
    // const table = useReactTable({
    //     data: defaultData,
    //     columns,
    //     getCoreRowModel: getCoreRowModel(),
    // })


    // Массив данных для турниров (можно вынести в пропсы или получать из API)
    const tournaments = [
        {
            id: 1,
            tourName: "От ICPC до IPSC: от кода к оружию",
            description: "Спецкурс пройдет в 17:30 в Клубе. Меня зовут Хет, я преподаватель по программированию и энтузиаст по практической стрельбе. Приглашаю вас на спецкурс, где мы познакомимся с миром IPSC (International Practical Shooting Confederation) – дисциплиной, которая по соревновательности ничем не уступает ICPC (International Collegiate Programming Contest).",
            organizer: "Тембай Вика",
            handle: "@dlwr_vita",
            additionalInfo: ""
        },
        {
            id: 2,
            description: "Описание турнира 2 и прочая нужная инфа о нем.",
            organizer: "Тембай Вика @dlwr_vita",
            additionalInfo: "Дополнительная информация о втором турнире"
        },
        {
            id: 3,
            description: "Описание турнира 3 и прочая нужная инфа о нем.",
            organizer: "Тембай Вика @dlwr_vita",
            additionalInfo: "Дополнительная информация о втором турнире"
        }
    ];

    return (
        <div className={styles.tournament}>
            <h1>Спортивные активности по футболу</h1>
            <div className={styles.tournamentsRow}>
                {tournaments.map((tour) => (
                    <div key={tour.id} className={styles.tourContainer}>
                        <p className={styles.tourName}>
                            {tour.tourName}
                        </p>
                        <p className={styles.tourInfo}>
                            {tour.description}
                        </p>
                        <div className={styles.tourAdmin}>
                            <p> Организатор турнира:
                                {tour.organizer}
                                <img src={telegramLogo} alt="Telegram" className={styles.logo}
                                     onClick={() => window.open(`https://t.me/${tour.handle}`, "_blank")}
                                />
                            </p>
                        </div>
                        <div className={styles.toggleSection}>
                            <button
                                onClick={() => toggleContainer(tour.id)}
                                className={styles.toggleButton}
                            >
                                <div className={styles.logo}>
                                    {openId === tour.id ? (
                                        <>
                                            <ChevronUp/>
                                        </>
                                    ) : (
                                        <>
                                            Команды-участники <ChevronDown/>
                                        </>
                                    )}</div>
                            </button>
                            {openId === tour.id && (
                                <div className={styles.tourLists}>
                                    <p>
                                        <div className={styles.teamsRow}>
                                            {defaultData.map((team) => (
                                                <div key={team.id} className={styles.teamCard}>
                                                    <h3>{team.id} — {team.name}</h3>
                                                    <p> {team.commander}<PiStarFill className={styles.logo}
                                                                                    color={"#BF9298"}/></p>
                                                    <ul>
                                                        {team.members.map((m, idx) => (
                                                            <li key={idx}>{m}</li>
                                                        ))}
                                                    </ul>
                                                </div>
                                            ))}
                                        </div>
                                        {tour.additionalInfo}</p>
                                </div>
                            )}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default Tournament;