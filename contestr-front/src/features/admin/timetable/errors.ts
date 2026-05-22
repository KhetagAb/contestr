import { formatDurationCompact, minActiveDurationSeconds } from "./time";

export function activeSegmentMinDurationError(minSeconds: number): string {
    return `Нельзя сделать текущий сегмент короче уже прошедшего времени. Минимум сейчас: ${formatDurationCompact(minSeconds)}.`;
}

function translateElapsedDurationError(details: string): string | null {
    const match = details.match(
        /duration cannot be less than elapsed time in segment \((\d+)s\)/,
    );
    if (!match) {
        return null;
    }
    const elapsedIn = Number(match[1]);
    if (!Number.isFinite(elapsedIn)) {
        return null;
    }
    return activeSegmentMinDurationError(minActiveDurationSeconds(elapsedIn, 0));
}

const ERROR_TRANSLATIONS: Record<string, string> = {
    "bad request": "Некорректный запрос.",
    conflict: "Конфликт состояния.",
    "contest not found": "Контест не найден. Сначала дождитесь синхронизации контеста.",
    "contest has no participants": "В контесте нет участников.",
    "contest already registered": "Этот контест уже зарегистрирован.",
    "contest not registered": "Контест не зарегистрирован.",
    "handle mapping not found": "Маппинг handle не найден.",
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
    "manual start disabled while auto start is enabled":
        "Ручной запуск недоступен при включённом автозапуске. Отключите автозапуск.",
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
        if (details.includes("overlap")) {
            return "Некорректное расписание: туры пересекаются по времени.";
        }
        if (details.includes("duration cannot be less than elapsed")) {
            return (
                translateElapsedDurationError(details) ??
                activeSegmentMinDurationError(60)
            );
        }
        return "Некорректное расписание.";
    }

    if (key.includes("only the first not started tour can be started")) {
        return "Можно запустить только первый незапущенный тур.";
    }
    if (key.includes("manual start disabled while auto start is enabled")) {
        return "Ручной запуск недоступен при включённом автозапуске. Отключите автозапуск.";
    }
    if (key.includes("contest has not started yet")) {
        return "Контест ещё не начался.";
    }
    if (key.startsWith("failed to get contest")) {
        return "Контест не найден. Сначала дождитесь синхронизации контеста.";
    }
    if (key.includes("codeforces contest") || key.includes("error getting contest standings")) {
        return "Не удалось получить данные контеста из Codeforces. Проверьте ID контеста, доступ менеджера и API-ключи.";
    }
    if (
        key.includes("object storage") ||
        key.includes("put object") ||
        key.includes("save metadata") ||
        key.includes("загрузка pdf") ||
        key.includes("сохранение метаданных")
    ) {
        return value;
    }
    if (key.startsWith("failed to ")) {
        return "Не удалось выполнить действие.";
    }

    return value.length > 80 ? value : "Не удалось выполнить действие.";
}

export async function readApiError(response: Response) {
    try {
        const body = (await response.json()) as { message?: string };
        return translateAdminMessage(body.message || response.statusText);
    } catch {
        return translateAdminMessage(response.statusText);
    }
}
