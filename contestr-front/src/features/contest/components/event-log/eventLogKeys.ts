import type { RegattaEvent } from "@/client";

/** Стабильный ключ события — совпадает между poll и первым рендером после reload. */
export function regattaEventKey(event: RegattaEvent, tourIndex: number): string {
    return [
        tourIndex,
        event.type,
        event.time_sec,
        event.participant_id,
        event.problem_code,
        event.verdict ?? "",
    ].join("|");
}

const STAGGER_MS = 45;
const STAGGER_CAP = 12;

/** Задержка каскада появления (сверху вниз внутри тура). */
export function eventEnterDelayMs(index: number): number {
    return Math.min(index, STAGGER_CAP) * STAGGER_MS;
}
