const ERROR_TRANSLATIONS: Record<string, string> = {
    "bad request": "Некорректный запрос.",
    conflict: "Конфликт состояния.",
    "contest not found": "Контест не найден. Сначала дождитесь синхронизации контеста.",
    "failed to create timetable": "Не удалось создать расписание.",
    "failed to delete timetable": "Не удалось удалить расписание.",
    "failed to get contest": "Не удалось получить контест.",
    "failed to get timetable": "Не удалось получить расписание.",
    "failed to issue token": "Не удалось выпустить токен.",
    "failed to move tour": "Не удалось перенести тур.",
    "failed to start tour": "Не удалось запустить тур.",
    "failed to update timetable": "Не удалось обновить расписание.",
    "failed to update timetable after starting tour": "Не удалось обновить расписание после запуска тура.",
    forbidden: "Доступ запрещён.",
    "internal server error": "Внутренняя ошибка сервера.",
    "invalid or expired token": "Токен некорректен или истёк.",
    "invalid request body": "Некорректное тело запроса.",
    "invalid timetable": "Некорректное расписание.",
    "missing authorization header": "Не передан заголовок авторизации.",
    "not found": "Не найдено.",
    "timetable already exists": "Расписание уже существует.",
    "timetable not found": "Расписание не найдено.",
    "tour already started": "Тур уже запущен.",
    "tour not found": "Тур не найден.",
    "tour not found in timetable": "Тур не найден в расписании.",
    unauthorized: "Не авторизован.",
};

export function translateAdminMessage(message?: string) {
    const value = message?.trim();
    if (!value) {
        return "Неизвестная ошибка.";
    }

    if (/[А-Яа-яЁё]/.test(value)) {
        return value;
    }

    const key = value.toLowerCase();
    const direct = ERROR_TRANSLATIONS[key];
    if (direct) {
        return direct;
    }

    if (key.startsWith("invalid timetable:")) {
        const details = key.slice("invalid timetable:".length).trim();
        if (details.includes("duration must be positive")) {
            return "Некорректное расписание: длительность тура должна быть положительной.";
        }
        return "Некорректное расписание.";
    }

    if (key.includes("only the first not started tour can be started")) {
        return "Можно запустить только первый незапущенный тур.";
    }
    if (key.includes("contest has not started yet")) {
        return "Контест ещё не начался.";
    }
    if (key.startsWith("failed to get contest")) {
        return "Контест не найден. Сначала дождитесь синхронизации контеста.";
    }
    if (key.startsWith("failed to ")) {
        return "Не удалось выполнить действие.";
    }

    return "Не удалось выполнить действие.";
}

export async function readApiError(response: Response) {
    try {
        const body = (await response.json()) as { message?: string };
        return translateAdminMessage(body.message || response.statusText);
    } catch {
        return translateAdminMessage(response.statusText);
    }
}
