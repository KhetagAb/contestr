import { useQuery } from "@tanstack/react-query";
import { getContests } from "./client/sdk.gen";
import type { RegisteredContestItem } from "./client/types.gen";

export function useContests() {
    const query = useQuery({
        queryKey: ["contests"],
        queryFn: async () => {
            const { data, error } = await getContests();
            if (error) {
                throw new Error("Не удалось загрузить список контестов");
            }
            return data ?? [];
        },
        staleTime: 30_000,
    });

    return {
        contests: (query.data ?? []) as RegisteredContestItem[],
        isLoading: query.isLoading,
        isError: query.isError,
        refetch: query.refetch,
    };
}
