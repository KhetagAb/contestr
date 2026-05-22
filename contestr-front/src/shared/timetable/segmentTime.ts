import type { TimelineSegment } from "@/client/types.gen";

function pad2(n: number) {
    return String(n).padStart(2, "0");
}

export function segmentLabel(segment: TimelineSegment): string {
    if (segment.kind === "pause") {
        return "Перерыв";
    }
    if (segment.round != null && segment.round > 0) {
        return `Тур ${segment.round}`;
    }
    return "Тур";
}

/** Подпись для мини-таймлайна (всегда с номером тура, если это тур). */
export function miniTimelineSegmentLabel(
    segment: TimelineSegment,
    segments: TimelineSegment[],
    index: number,
): string {
    if (segment.kind === "pause") {
        return "Перерыв";
    }
    if (segment.round != null && segment.round > 0) {
        return `Тур ${segment.round}`;
    }
    let tourNum = 0;
    for (let i = 0; i <= index; i++) {
        if (segments[i].kind === "tour") {
            tourNum++;
        }
    }
    return tourNum > 0 ? `Тур ${tourNum}` : "Тур";
}

export function formatCountdown(seconds: number): string {
    const total = Math.max(0, Math.floor(seconds));
    const h = Math.floor(total / 3600);
    const m = Math.floor((total % 3600) / 60);
    const s = total % 60;
    return `${pad2(h)}:${pad2(m)}:${pad2(s)}`;
}

export function segmentRemainingSeconds(segment: TimelineSegment, elapsed: number): number {
    const end = segment.start_time + segment.duration;
    return Math.max(0, end - elapsed);
}

export function segmentUntilStartSeconds(segment: TimelineSegment, elapsed: number): number {
    return Math.max(0, segment.start_time - elapsed);
}

/** Секунды до contest_start_time; null если старт уже прошёл или время не задано. */
export function secondsUntilContestStart(
    contestStartIso?: string,
    nowMs = Date.now(),
): number | null {
    if (!contestStartIso) {
        return null;
    }
    const startMs = new Date(contestStartIso).getTime();
    if (Number.isNaN(startMs)) {
        return null;
    }
    const sec = Math.ceil((startMs - nowMs) / 1000);
    return sec > 0 ? sec : null;
}
