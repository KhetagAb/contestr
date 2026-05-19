import type { TourConfig } from "../client/types.gen";

function pad2(n: number) {
    return String(n).padStart(2, "0");
}

function formatHHMMParts(h: number, m: number): string {
    return `${pad2(h)}:${pad2(m)}`;
}

/** Секунды от нулевой точки → «ЧЧ:ММ» (секунды отбрасываются). */
export function formatClockHHMM(seconds: number): string {
    const totalMinutes = Math.floor(Math.max(0, seconds) / 60);
    const h = Math.floor(totalMinutes / 60) % 24;
    const m = totalMinutes % 60;
    return formatHHMMParts(h, m);
}

/** Абсолютное время старта тура по ISO старта контеста и смещению в секундах. */
export function formatTourClock(contestStartIso: string | undefined, offsetSeconds: number): string {
    if (!contestStartIso) {
        return formatClockHHMM(offsetSeconds);
    }
    const start = new Date(contestStartIso);
    if (Number.isNaN(start.getTime())) {
        return formatClockHHMM(offsetSeconds);
    }
    const at = new Date(start.getTime() + offsetSeconds * 1000);
    return formatHHMMParts(at.getHours(), at.getMinutes());
}

/** Секунды от старта → «+N мин.» (для подсказок и т.п.). */
export function formatOffset(seconds: number): string {
    const sign = seconds < 0 ? "-" : "+";
    const minutes = Math.floor(Math.abs(seconds) / 60);
    return `${sign}${minutes} мин.`;
}

/** Длительность тура для отображения: «30 мин.» */
export function formatDuration(seconds: number): string {
    const minutes = Math.max(1, Math.round(Math.max(0, seconds) / 60));
    return `${minutes} мин.`;
}

/** Длительность без точки: «30 мин» */
export function formatDurationCompact(seconds: number): string {
    const minutes = Math.max(1, Math.round(Math.max(0, seconds) / 60));
    return `${minutes} мин`;
}

/** Число минут для поля ввода (без суффикса). */
export function formatDurationInputValue(seconds: number): string {
    const minutes = Math.max(1, Math.round(Math.max(0, seconds) / 60));
    return String(minutes);
}

/** Парсит целое число минут. */
export function parseDurationInput(input: string): number | null {
    const trimmed = input.trim().replace(/\s*мин\.?$/i, "").trim();
    if (!/^\d+$/.test(trimmed)) {
        return null;
    }
    const minutes = Number(trimmed);
    if (minutes <= 0) {
        return null;
    }
    return minutes * 60;
}

/** Прошедшее время контеста «ЧЧ:ММ» (секунды отбрасываются). */
export function formatElapsed(seconds: number): string {
    return formatClockHHMM(seconds);
}

export function formatContestClock(iso?: string): string {
    if (!iso) {
        return "—";
    }
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) {
        return "—";
    }
    return date.toLocaleString("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
    });
}

/** Пересчёт start_time: запущенные не трогаем, цепочка от первого незапущенного. */
export function rebuildChain(tourTimes: TourConfig[]): TourConfig[] {
    const tours = tourTimes.map((t) => ({ ...t }));
    if (tours.length === 0) {
        return tours;
    }

    let anchor = tours.length;
    for (let i = 0; i < tours.length; i++) {
        if (!tours[i].started) {
            anchor = i;
            break;
        }
    }

    if (anchor >= tours.length) {
        return tours;
    }

    if (anchor === 0) {
        tours[0].start_time = 0;
    } else {
        const prev = tours[anchor - 1];
        tours[anchor].start_time = prev.start_time + prev.duration;
    }

    for (let i = anchor + 1; i < tours.length; i++) {
        const prev = tours[i - 1];
        tours[i].start_time = prev.start_time + prev.duration;
    }

    return tours;
}

export function defaultTourDuration(previous?: TourConfig): number {
    return previous?.duration && previous.duration > 0 ? previous.duration : 1800;
}

export function createTourFromPrevious(previous?: TourConfig): TourConfig {
    const duration = defaultTourDuration(previous);
    const start_time = previous ? previous.start_time + previous.duration : 0;
    return { start_time, duration, started: false };
}

export function toursFingerprint(tours: TourConfig[]): string {
    return JSON.stringify(tours.map((t) => ({ duration: t.duration, started: t.started })));
}
