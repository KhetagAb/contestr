import type { TimelineSegment, TimetableView } from "@/client/types.gen";
import {
    effectiveUntilStartSeconds,
    isFirstCompetitiveTour,
    segmentLabel,
    segmentRemainingSeconds,
    secondsUntilContestStart,
} from "./segmentTime";

export type PhaseDisplay = {
    empty: boolean;
    finished: boolean;
    heroSegment: TimelineSegment | null;
    heroLabel: string;
    statusText: string;
    remainingSeconds: number | null;
    isPause: boolean;
};

function lastPastSegment(segments: TimelineSegment[]): TimelineSegment | null {
    for (let i = segments.length - 1; i >= 0; i--) {
        if (segments[i].status === "past") {
            return segments[i];
        }
    }
    return null;
}

function nextUpSegment(segments: TimelineSegment[]): TimelineSegment | null {
    return (
        segments.find((s) => s.status === "starting") ??
        segments.find((s) => s.status === "next") ??
        segments.find((s) => s.status === "future") ??
        null
    );
}

export function derivePhaseDisplay(
    view: TimetableView | null | undefined,
    elapsed: number,
): PhaseDisplay {
    const segments = view?.timeline_segments ?? [];
    const hasSchedule =
        segments.length > 0 || (view?.pending_slots?.length ?? 0) > 0;

    if (!view || !hasSchedule) {
        return {
            empty: true,
            finished: false,
            heroSegment: null,
            heroLabel: "Расписание не задано",
            statusText: "",
            remainingSeconds: null,
            isPause: false,
        };
    }

    const active = segments.find((s) => s.status === "active") ?? null;
    if (active) {
        return {
            empty: false,
            finished: false,
            heroSegment: active,
            heroLabel: segmentLabel(active),
            statusText: "идёт",
            remainingSeconds: segmentRemainingSeconds(active, elapsed),
            isPause: active.kind === "pause",
        };
    }

    const hasUpcoming =
        segments.some((s) =>
            ["next", "starting", "future"].includes(s.status),
        ) || (view.pending_slots?.length ?? 0) > 0;

    const allPast =
        segments.length > 0 &&
        segments.every((s) => s.status === "past") &&
        (view.pending_slots?.length ?? 0) === 0;

    if (allPast || (!hasUpcoming && segments.length > 0)) {
        const last = lastPastSegment(segments);
        return {
            empty: false,
            finished: true,
            heroSegment: last,
            heroLabel: last ? segmentLabel(last) : "Контест",
            statusText: "завершён",
            remainingSeconds: null,
            isPause: last?.kind === "pause",
        };
    }

    const next = nextUpSegment(segments);
    const last = lastPastSegment(segments);
    const hero = next ?? last;

    const untilContest = secondsUntilContestStart(view.contest_start_time);

    if (segments.length === 0 && (view.pending_slots?.length ?? 0) > 0) {
        const round = view.next_tour_number;
        const waitingForContest = untilContest != null;
        return {
            empty: false,
            finished: false,
            heroSegment: null,
            heroLabel: round != null ? `Тур ${round}` : "Скоро",
            statusText: waitingForContest ? "" : "ожидание старта",
            remainingSeconds: untilContest,
            isPause: false,
        };
    }

    let statusText = "между турами";
    let remaining: number | null = null;

    if (next) {
        const untilSegment = effectiveUntilStartSeconds(next, elapsed, view.contest_start_time);
        const firstTour = isFirstCompetitiveTour(next);

        if (firstTour && untilContest != null) {
            remaining = untilSegment;
            statusText = "";
        } else if (untilSegment > 0) {
            remaining = untilSegment;
            statusText = firstTour ? "" : "ожидание";
        } else if (next.status === "starting" || next.start_time <= elapsed) {
            statusText =
                next.round != null && next.kind === "tour"
                    ? `скоро ${segmentLabel(next)}`
                    : "скоро";
            remaining = null;
        } else {
            statusText = "ожидание";
            remaining = null;
        }
    } else if (view.next_tour_number != null) {
        const upcoming = segments.find((s) =>
            ["next", "starting", "future"].includes(s.status),
        );
        const untilSegment = upcoming
            ? effectiveUntilStartSeconds(upcoming, elapsed, view.contest_start_time)
            : view.next_tour_number === 1
              ? (untilContest ?? 0)
              : 0;

        if (view.next_tour_number === 1 && untilContest != null) {
            remaining = untilSegment;
            statusText = "";
        } else if (untilSegment > 0) {
            remaining = untilSegment;
            statusText =
                view.next_tour_number === 1
                    ? ""
                    : `скоро Тур ${view.next_tour_number}`;
        } else {
            statusText = `скоро Тур ${view.next_tour_number}`;
        }
    }

    return {
        empty: false,
        finished: false,
        heroSegment: hero,
        heroLabel: hero ? segmentLabel(hero) : "Между турами",
        statusText,
        remainingSeconds: remaining,
        isPause: hero?.kind === "pause",
    };
}
