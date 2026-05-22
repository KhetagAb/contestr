import { useQuery } from "@tanstack/react-query";
import { useSearchParam } from "react-use";
import { getRegattaContestStandingsOptions } from "@/client/@tanstack/react-query.gen";
import { isReserveDisplayName } from "@/shared/utils/reserveParticipant";

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
            rows: data.rows.filter((row) => !isReserveDisplayName(row.display_name)),
            events: data.events.filter(
                (event) => !isReserveDisplayName(event.display_name),
            ),
        }),
    });
}
