import type { TimelineSegment } from "../client/types.gen";
import { resolveSegmentVisualState, TOUR_VISUAL_BG } from "./tourVisualState";

export type ActiveTourProgress = {
    index: number;
    fill: number;
    colorFrom: string;
    colorTo: string;
};

export function computeActiveSegmentProgress(
    segments: TimelineSegment[],
    elapsedSeconds: number,
): ActiveTourProgress | null {
    const index = segments.findIndex((s) => s.status === "active");
    if (index < 0) {
        return null;
    }

    const segment = segments[index];
    if (segment.duration <= 0) {
        return null;
    }

    const fill = Math.min(
        1,
        Math.max(0, (elapsedSeconds - segment.start_time) / segment.duration),
    );

    const nextSegment = segments[index + 1];
    const nextVisual = nextSegment
        ? resolveSegmentVisualState(nextSegment)
        : "future";

    return {
        index,
        fill,
        colorFrom: TOUR_VISUAL_BG.active,
        colorTo: TOUR_VISUAL_BG[nextVisual],
    };
}
