import { useMemo } from "react";
import { useQueries } from "@tanstack/react-query";
import { getRegattaContestParticipantsOptions } from "@/client/@tanstack/react-query.gen";
import type { ContestParticipantItem } from "@/client/types.gen";

export function useContestParticipantCounts(contestIds: number[], enabled: boolean) {
    const queries = useQueries({
        queries: contestIds.map((contestId) => ({
            ...getRegattaContestParticipantsOptions({
                path: { contest_id: contestId },
            }),
            enabled: enabled && contestId > 0,
            staleTime: 60_000,
        })),
    });

    const countByContestId = useMemo(() => {
        const map = new Map<number, number>();
        contestIds.forEach((contestId, index) => {
            const data = queries[index]?.data as ContestParticipantItem[] | undefined;
            if (Array.isArray(data)) {
                map.set(contestId, data.length);
            }
        });
        return map;
    }, [contestIds, queries]);

    const isLoading = queries.some((query) => query.isLoading);

    return { countByContestId, isLoading };
}
