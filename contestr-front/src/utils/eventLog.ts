function pad2(n: number): string {
    return String(n).padStart(2, "0");
}

/** Секунды от старта контеста → «HH:MM:SS». */
export function formatContestTime(timeSec: number): string {
    const total = Math.max(0, Math.floor(timeSec));
    const h = Math.floor(total / 3600);
    const m = Math.floor((total % 3600) / 60);
    const s = total % 60;
    return `${pad2(h)}:${pad2(m)}:${pad2(s)}`;
}

/** Номер тура из кода задачи (например «2A» → 2). */
export function tourIndexFromProblemCode(problemCode: string): number {
    const match = /^(\d+)/.exec(problemCode);
    return match ? Number.parseInt(match[1], 10) : 0;
}
