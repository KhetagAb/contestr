import { useQuery } from "@tanstack/react-query"
import { getRegattaContestStandingsOptions } from "./client/@tanstack/react-query.gen"
import { useSearchParam } from "react-use";

export const useCurRegattaData = () => {
    const contestId = parseInt(useSearchParam("contestId") || "", 10);
    return useQuery(getRegattaContestStandingsOptions({
        path: {
            contest_id: contestId
        }
    }))
}