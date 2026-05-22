import { useQuery } from "@tanstack/react-query";
import { useSearchParam } from "react-use";
import type { GetRegattaContestStandingsResponse } from "@/client/types.gen";
import { getRegattaContestStandingsOptions } from "@/client/@tanstack/react-query.gen";
import { isReserveDisplayName } from "@/shared/utils/reserveParticipant";

/** Туры ещё не запускались (регата не начата). */
export function hasRegattaStartedTours(
    data: Pick<
        GetRegattaContestStandingsResponse,
        "current_tour_start_time" | "current_tour_duration"
    >,
): boolean {
    return (
        (data.current_tour_start_time ?? 0) > 0 ||
        (data.current_tour_duration ?? 0) > 0
    );
}

export function useContestStandings() {
    const contestId = parseInt(useSearchParam("contestId") || "", 10);
    const enabled = Number.isFinite(contestId) && contestId > 0;
    return useQuery({
        ...getRegattaContestStandingsOptions({
            path: {
                contest_id: contestId,
            },
        }),
        enabled,
        refetchInterval: enabled ? 5_000 : false,
        select: (data) => ({
            ...data,
            rows: (data.rows ?? []).filter(
                (row) => !isReserveDisplayName(row.display_name),
            ),
            events: (data.events ?? []).filter(
                (event) => !isReserveDisplayName(event.display_name),
            ),
        }),
    });
}
