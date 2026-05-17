import { type FormEvent, useMemo, useState } from "react";
import { Plus, RefreshCw, Save, Trash2 } from "lucide-react";
import { adminAuthHeaders } from "./adminAuth";
import { CONTESTS } from "./consts";

type TourConfig = {
    start_time: number;
    duration: number;
    started: boolean;
};

type ToursTimetable = {
    contest_id: number;
    tour_times: TourConfig[];
};

type ErrorResponse = {
    message?: string;
};

async function readError(response: Response) {
    try {
        const body = (await response.json()) as ErrorResponse;
        return body.message || response.statusText;
    } catch {
        return response.statusText;
    }
}

function normalizeTour(tour: TourConfig): TourConfig {
    return {
        start_time: Number(tour.start_time) || 0,
        duration: Number(tour.duration) || 0,
        started: Boolean(tour.started),
    };
}

function nextTourDefaults(tours: TourConfig[]): TourConfig {
    const last = tours[tours.length - 1];
    if (!last) {
        return { start_time: 0, duration: 1800, started: false };
    }
    return {
        start_time: last.start_time + last.duration,
        duration: last.duration || 1800,
        started: false,
    };
}

export default function AdminTimetable() {
    const defaultContestId = CONTESTS[0]?.id?.toString() ?? "";
    const [contestId, setContestId] = useState(defaultContestId);
    const [tours, setTours] = useState<TourConfig[]>([]);
    const [moveTourNumber, setMoveTourNumber] = useState("1");
    const [moveStartTime, setMoveStartTime] = useState("0");
    const [nextTour, setNextTour] = useState<TourConfig | null>(null);
    const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
    const [message, setMessage] = useState("");

    const numericContestId = useMemo(() => Number(contestId), [contestId]);
    const canSubmit = Number.isInteger(numericContestId) && numericContestId > 0;

    const setResult = (nextStatus: typeof status, nextMessage: string) => {
        setStatus(nextStatus);
        setMessage(nextMessage);
    };

    const loadTimetable = async () => {
        if (!canSubmit) {
            setResult("error", "ID контеста должен быть положительным целым числом.");
            return;
        }

        setStatus("loading");
        setMessage("");
        setNextTour(null);

        const response = await fetch(`/api/admin/timetables/${numericContestId}`, {
            headers: adminAuthHeaders(),
        });

        if (response.status === 404) {
            setTours([]);
            setResult("idle", "Для этого контеста расписание не задано.");
            return;
        }
        if (!response.ok) {
            setResult("error", await readError(response));
            return;
        }

        const timetable = (await response.json()) as ToursTimetable;
        setTours((timetable.tour_times || []).map(normalizeTour));
        setResult("success", "Расписание загружено.");
    };

    const saveTimetable = async (event?: FormEvent) => {
        event?.preventDefault();
        if (!canSubmit) {
            setResult("error", "ID контеста должен быть положительным целым числом.");
            return;
        }

        setStatus("loading");
        setMessage("");

        const response = await fetch(`/api/admin/timetables/${numericContestId}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                ...adminAuthHeaders(),
            },
            body: JSON.stringify({
                contest_id: numericContestId,
                tour_times: tours.map(normalizeTour),
            } satisfies ToursTimetable),
        });

        if (!response.ok) {
            setResult("error", await readError(response));
            return;
        }

        const timetable = (await response.json()) as ToursTimetable;
        setTours((timetable.tour_times || []).map(normalizeTour));
        setResult("success", "Расписание сохранено.");
    };

    const deleteTimetable = async () => {
        if (!canSubmit) {
            setResult("error", "ID контеста должен быть положительным целым числом.");
            return;
        }
        if (!window.confirm("Удалить расписание для этого контеста?")) {
            return;
        }

        setStatus("loading");
        setMessage("");

        const response = await fetch(`/api/admin/timetables/${numericContestId}`, {
            method: "DELETE",
            headers: adminAuthHeaders(),
        });

        if (!response.ok && response.status !== 404) {
            setResult("error", await readError(response));
            return;
        }

        setTours([]);
        setNextTour(null);
        setResult("success", "Расписание удалено.");
    };

    const moveTour = async () => {
        if (!canSubmit) {
            setResult("error", "ID контеста должен быть положительным целым числом.");
            return;
        }

        const tourNumber = Number(moveTourNumber);
        const startTime = Number(moveStartTime);
        if (!Number.isInteger(tourNumber) || tourNumber < 1 || !Number.isInteger(startTime) || startTime < 0) {
            setResult("error", "Номер тура должен быть положительным, а время начала - неотрицательным целым числом.");
            return;
        }

        setStatus("loading");
        setMessage("");

        const response = await fetch(`/api/admin/timetables/${numericContestId}/tours/${tourNumber}/move`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                ...adminAuthHeaders(),
            },
            body: JSON.stringify({ start_time: startTime }),
        });

        if (!response.ok) {
            setResult("error", await readError(response));
            return;
        }

        const timetable = (await response.json()) as ToursTimetable;
        setTours((timetable.tour_times || []).map(normalizeTour));
        setResult("success", "Тур перенесён.");
    };

    const loadNextTour = async () => {
        if (!canSubmit) {
            setResult("error", "ID контеста должен быть положительным целым числом.");
            return;
        }

        setStatus("loading");
        setMessage("");

        const response = await fetch(`/api/admin/timetables/${numericContestId}/first-not-started`, {
            headers: adminAuthHeaders(),
        });

        if (!response.ok) {
            setNextTour(null);
            setResult("error", await readError(response));
            return;
        }

        const tour = normalizeTour((await response.json()) as TourConfig);
        setNextTour(tour);
        setResult("success", "Первый незапущенный тур загружен.");
    };

    const updateTour = (index: number, patch: Partial<TourConfig>) => {
        setTours((current) =>
            current.map((tour, currentIndex) =>
                currentIndex === index ? normalizeTour({ ...tour, ...patch }) : tour
            )
        );
    };

    const removeTour = (index: number) => {
        setTours((current) => current.filter((_, currentIndex) => currentIndex !== index));
    };

    return (
        <section className="admin-timetable">
            <div className="admin-section-head">
                <h2>Расписание туров</h2>
                <button type="button" className="admin-icon-btn" onClick={loadTimetable} disabled={status === "loading"}>
                    <RefreshCw size={16} />
                    Загрузить
                </button>
            </div>

            <label className="admin-login-field admin-contest-field">
                <span>Контест</span>
                <div className="admin-contest-inputs">
                    <select value={contestId} onChange={(event) => setContestId(event.target.value)}>
                        {CONTESTS.map((contest) => (
                            <option key={contest.id} value={contest.id}>
                                {contest.name} ({contest.id})
                            </option>
                        ))}
                    </select>
                    <input
                        type="number"
                        min="1"
                        value={contestId}
                        onChange={(event) => setContestId(event.target.value)}
                    />
                </div>
            </label>

            <form onSubmit={saveTimetable} className="admin-timetable-form">
                <div className="admin-timetable-table">
                    <div className="admin-timetable-row admin-timetable-header">
                        <span>#</span>
                        <span>Старт, сек</span>
                        <span>Длительность, сек</span>
                        <span>Запущен</span>
                        <span />
                    </div>
                    {tours.map((tour, index) => (
                        <div className="admin-timetable-row" key={index}>
                            <span>{index + 1}</span>
                            <input
                                type="number"
                                min="0"
                                value={tour.start_time}
                                onChange={(event) => updateTour(index, { start_time: Number(event.target.value) })}
                            />
                            <input
                                type="number"
                                min="1"
                                value={tour.duration}
                                onChange={(event) => updateTour(index, { duration: Number(event.target.value) })}
                            />
                            <input
                                type="checkbox"
                                checked={tour.started}
                                onChange={(event) => updateTour(index, { started: event.target.checked })}
                            />
                            <button type="button" className="admin-row-btn" onClick={() => removeTour(index)}>
                                <Trash2 size={16} />
                            </button>
                        </div>
                    ))}
                </div>

                <div className="admin-login-actions admin-toolbar">
                    <button
                        type="button"
                        className="admin-icon-btn"
                        onClick={() => setTours((current) => [...current, nextTourDefaults(current)])}
                    >
                        <Plus size={16} />
                        Добавить тур
                    </button>
                    <button type="submit" className="admin-icon-btn admin-primary-btn" disabled={status === "loading"}>
                        <Save size={16} />
                        Сохранить
                    </button>
                    <button type="button" className="admin-icon-btn admin-danger-btn" onClick={deleteTimetable}>
                        <Trash2 size={16} />
                        Удалить
                    </button>
                </div>
            </form>

            <div className="admin-split-tools">
                <div className="admin-tool-block">
                    <h3>Перенос цепочки</h3>
                    <div className="admin-inline-fields">
                        <label className="admin-login-field">
                            <span>Тур</span>
                            <input
                                type="number"
                                min="1"
                                value={moveTourNumber}
                                onChange={(event) => setMoveTourNumber(event.target.value)}
                            />
                        </label>
                        <label className="admin-login-field">
                            <span>Новый старт</span>
                            <input
                                type="number"
                                min="0"
                                value={moveStartTime}
                                onChange={(event) => setMoveStartTime(event.target.value)}
                            />
                        </label>
                    </div>
                    <button type="button" className="admin-icon-btn" onClick={moveTour} disabled={status === "loading"}>
                        Перенести
                    </button>
                </div>

                <div className="admin-tool-block">
                    <h3>Следующий тур</h3>
                    {nextTour ? (
                        <p className="admin-result-text">
                            старт={nextTour.start_time}с длительность={nextTour.duration}с
                        </p>
                    ) : (
                        <p className="admin-result-text">Тур не выбран.</p>
                    )}
                    <button type="button" className="admin-icon-btn" onClick={loadNextTour} disabled={status === "loading"}>
                        Найти
                    </button>
                </div>
            </div>

            {message && (
                <p className={`admin-login-message admin-login-message--${status}`}>
                    {message}
                </p>
            )}
        </section>
    );
}
