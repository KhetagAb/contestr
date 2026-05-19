import type { ScheduleSlot, TimelineSegment, TimetableView } from "../client/types.gen";

function pad2(n: number) {
    return String(n).padStart(2, "0");
}

function formatHHMMParts(h: number, m: number): string {
    return `${pad2(h)}:${pad2(m)}`;
}

export function formatClockHHMM(seconds: number): string {
    const totalMinutes = Math.floor(Math.max(0, seconds) / 60);
    const h = Math.floor(totalMinutes / 60) % 24;
    const m = totalMinutes % 60;
    return formatHHMMParts(h, m);
}

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

export function formatDuration(seconds: number): string {
    const minutes = Math.max(1, Math.round(Math.max(0, seconds) / 60));
    return `${minutes} мин.`;
}

export function formatDurationCompact(seconds: number): string {
    const minutes = Math.max(1, Math.round(Math.max(0, seconds) / 60));
    return `${minutes} мин`;
}

export function formatDurationInputValue(seconds: number): string {
    const minutes = Math.max(1, Math.round(Math.max(0, seconds) / 60));
    return String(minutes);
}

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

export function formatElapsed(seconds: number): string {
    const total = Math.max(0, Math.floor(seconds));
    const h = Math.floor(total / 3600);
    const m = Math.floor((total % 3600) / 60);
    const s = total % 60;
    return `${pad2(h)}:${pad2(m)}:${pad2(s)}`;
}

export function formatContestClock(iso?: string): string {
    if (!iso) {
        return "—";
    }
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) {
        return "—";
    }
    const dd = pad2(date.getDate());
    const mm = pad2(date.getMonth() + 1);
    const yy = pad2(date.getFullYear() % 100);
    const hh = pad2(date.getHours());
    const min = pad2(date.getMinutes());
    const ss = pad2(date.getSeconds());
    return `${dd}.${mm}.${yy}, ${hh}:${min}:${ss}`;
}

export function defaultSlotDuration(previous?: ScheduleSlot): number {
    return previous?.duration && previous.duration > 0 ? previous.duration : 1800;
}

export function createSlotFromPrevious(previous?: ScheduleSlot): ScheduleSlot {
    return { duration: defaultSlotDuration(previous), kind: "tour" };
}

export function pendingSlotsFingerprint(slots: ScheduleSlot[]): string {
    return JSON.stringify(slots.map((s) => ({ duration: s.duration, kind: s.kind })));
}

function competitiveRoundsInFact(segments: TimelineSegment[]): number {
    return segments.filter((s) => s.sequence != null && s.kind === "tour").length;
}

/** Черновик pending-слотов на оси (пока не сохранён). */
export function buildDraftTimelineSegments(
    view: TimetableView,
    draftPending: ScheduleSlot[],
): TimelineSegment[] {
    const fact = (view.timeline_segments ?? []).filter((s) => s.sequence != null);
    let anchor = 0;
    for (const seg of fact) {
        anchor = Math.max(anchor, seg.start_time + seg.duration);
    }

    let cursor = anchor;
    let competitive = competitiveRoundsInFact(fact);
    const pendingSegments: TimelineSegment[] = [];

    for (let i = 0; i < draftPending.length; i++) {
        const slot = draftPending[i];
        const kind = slot.kind === "pause" ? "pause" : "tour";
        let round: number | null = null;
        if (kind === "tour") {
            competitive += 1;
            round = competitive;
        }
        pendingSegments.push({
            pending_index: i,
            kind,
            round,
            duration: slot.duration,
            start_time: cursor,
            status: "future",
            editable: true,
        });
        cursor += slot.duration;
    }

    return [...fact, ...pendingSegments];
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
