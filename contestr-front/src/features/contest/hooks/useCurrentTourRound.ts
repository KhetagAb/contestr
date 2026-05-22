import { useMemo } from "react";
import type { TimelineSegment } from "@/client/types.gen";
import { useContestTimetable } from "@/features/contest/hooks/useContestTimetable";

function roundFromActiveTour(segments: TimelineSegment[]): number | null {
    const active = segments.find((s) => s.kind === "tour" && s.status === "active");
    if (active?.round != null && active.round > 0) {
        return active.round;
    }
    return null;
}

/** Номер тура (round), по которому идёт сейчас соревнование. Без fallback на будущие туры из расписания. */
export function useCurrentTourRound(): number | null {
    const { data: view } = useContestTimetable();

    return useMemo(() => {
        const segments = view?.timeline_segments ?? [];
        if (segments.length === 0) {
            return null;
        }
        return roundFromActiveTour(segments);
    }, [view?.timeline_segments]);
}
