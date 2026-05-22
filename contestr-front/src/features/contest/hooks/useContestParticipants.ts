import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { useSearchParam } from "react-use";
import type { ContestParticipantItem } from "@/client/types.gen";
import { getRegattaContestParticipantsOptions } from "@/client/@tanstack/react-query.gen";
import { useContestStandings } from "@/features/contest/hooks/useContestStandings";

export function useContestParticipants() {
    const contestId = parseInt(useSearchParam("contestId") || "", 10);
    const enabled = Number.isFinite(contestId) && contestId > 0;

    const participantsQuery = useQuery({
        ...getRegattaContestParticipantsOptions({
            path: { contest_id: contestId },
        }),
        enabled,
        staleTime: 60_000,
    });

    const { data: standings } = useContestStandings();

    const participants = useMemo((): ContestParticipantItem[] => {
        const fromApi = participantsQuery.data ?? [];
        if (fromApi.length > 0) {
            return [...fromApi].sort((a, b) =>
                a.display_name.localeCompare(b.display_name, "ru"),
            );
        }

        const rows = standings?.rows ?? [];
        const seen = new Set<string>();
        const fromRows: ContestParticipantItem[] = [];
        for (const row of rows) {
            if (!row.user_id || seen.has(row.user_id)) {
                continue;
            }
            seen.add(row.user_id);
            fromRows.push({
                participant_id: row.user_id,
                display_name: row.display_name?.trim() || row.user_id,
            });
        }
        return fromRows.sort((a, b) =>
            a.display_name.localeCompare(b.display_name, "ru"),
        );
    }, [participantsQuery.data, standings?.rows]);

    return {
        participants,
        isLoading: participantsQuery.isLoading && participants.length === 0,
        isError: participantsQuery.isError,
    };
}
