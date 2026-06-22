function pad2(n: number): string {
    return String(n).padStart(2, "0");
}

/** Дата и время старта контеста для подписи в UI (как в расписании туров). */
export function formatContestStartLabel(iso?: string): string | null {
    if (!iso) {
        return null;
    }
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) {
        return null;
    }
    const dd = pad2(date.getDate());
    const mm = pad2(date.getMonth() + 1);
    const yy = pad2(date.getFullYear() % 100);
    const hh = pad2(date.getHours());
    const min = pad2(date.getMinutes());
    return `${dd}.${mm}.${yy}, ${hh}:${min}`;
}
