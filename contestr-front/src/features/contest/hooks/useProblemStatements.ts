import { useQuery } from "@tanstack/react-query";
import { useSearchParam } from "react-use";
import { getRegattaContestProblemStatementsOptions } from "@/client/@tanstack/react-query.gen";

export function useProblemStatements() {
    const contestId = parseInt(useSearchParam("contestId") || "", 10);
    const enabled = Number.isFinite(contestId) && contestId > 0;

    return useQuery({
        ...getRegattaContestProblemStatementsOptions({
            path: { contest_id: contestId },
        }),
        enabled,
        refetchInterval: enabled ? 15_000 : false,
        select: (data) => data.statements ?? {},
    });
}
