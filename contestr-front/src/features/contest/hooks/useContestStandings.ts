import { useQuery } from "@tanstack/react-query";
import { useSearchParam } from "react-use";
import { getRegattaContestStandingsOptions } from "@/client/@tanstack/react-query.gen";

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
    });
}
