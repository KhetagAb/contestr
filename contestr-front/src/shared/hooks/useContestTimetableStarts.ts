import { useMemo } from "react";
import { useQueries } from "@tanstack/react-query";
import { getRegattaContestTimetableOptions } from "@/client/@tanstack/react-query.gen";
import type { TimetableView } from "@/client/types.gen";

/** contest_start_time из TimetableView — тот же источник, что и в расписании туров. */
export function useContestTimetableStarts(contestIds: number[], enabled: boolean) {
    const queries = useQueries({
        queries: contestIds.map((contestId) => ({
            ...getRegattaContestTimetableOptions({
                path: { contest_id: contestId },
            }),
            enabled: enabled && contestId > 0,
            staleTime: 60_000,
        })),
    });

    const startTimeByContestId = useMemo(() => {
        const map = new Map<number, string>();
        contestIds.forEach((contestId, index) => {
            const view = queries[index]?.data as TimetableView | undefined;
            if (view?.contest_start_time) {
                map.set(contestId, view.contest_start_time);
            }
        });
        return map;
    }, [contestIds, queries]);

    return { startTimeByContestId };
}
