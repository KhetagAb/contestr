import type { TimelineSegment } from "@/client/types.gen";

export type TourVisualState = "past" | "active" | "next" | "starting" | "future";

export const TOUR_VISUAL_BG: Record<TourVisualState, string> = {
    past: "#a8b8a6",
    active: "#bdd7bf",
    next: "#e8dfd0",
    starting: "#e8dfd0",
    future: "#e8dfd0",
};

export function resolveSegmentVisualState(segment: TimelineSegment): TourVisualState {
    if (segment.status === "past") {
        return "past";
    }
    if (segment.status === "active") {
        return "active";
    }
    if (segment.status === "next") {
        return "next";
    }
    if (segment.status === "starting") {
        return "starting";
    }
    return "future";
}

export function visualStateClass(state: TourVisualState): string {
    return `tt-visual--${state}`;
}
