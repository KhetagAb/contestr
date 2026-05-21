import { useQuery } from "@tanstack/react-query";
import { useSearchParam } from "react-use";
import { getRegattaContestTimetableOptions } from "@/client/@tanstack/react-query.gen";

export function useContestTimetable() {
    const contestId = parseInt(useSearchParam("contestId") || "", 10);
    const enabled = Number.isFinite(contestId) && contestId > 0;
    return useQuery({
        ...getRegattaContestTimetableOptions({
            path: {
                contest_id: contestId,
            },
        }),
        enabled,
        refetchInterval: enabled ? 5_000 : false,
    });
}
