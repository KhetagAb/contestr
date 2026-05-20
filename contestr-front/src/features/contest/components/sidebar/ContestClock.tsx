import { useEffect, useState } from "react";
import { useContestStandings } from "@/features/contest/hooks/useContestStandings";

export function ContestClock() {
    const { data } = useContestStandings();
    const [now, setNow] = useState(Date.now());

    useEffect(() => {
        const interval = setInterval(() => {
            setNow(Date.now());
        }, 1000);
        return () => clearInterval(interval);
    }, []);

    if (!data) {
        return null;
    }

    const startTimeSeconds = data.current_tour_start_time;
    if (startTimeSeconds === undefined || startTimeSeconds === null || isNaN(startTimeSeconds)) {
        return <div>00:00:00</div>;
    }

    const startTime = startTimeSeconds * 1000;
    if (isNaN(startTime)) {
        return <div>00:00:00</div>;
    }

    const tourDurationMs = (data.current_tour_duration ?? 0) * 60 * 1000;
    const endTime = startTime + tourDurationMs;
    const remaining = Math.max(0, endTime - now);
    if (isNaN(remaining)) {
        return <div>00:00:00</div>;
    }

    const formatTime = (ms: number) => {
        if (isNaN(ms) || ms < 0) {
            return "00:00:00";
        }
        const totalSeconds = Math.floor(ms / 1000);
        const hours = Math.floor(totalSeconds / 3600);
        const minutes = Math.floor((totalSeconds % 3600) / 60);
        const seconds = totalSeconds % 60;
        return `${hours.toString().padStart(2, "0")}:${minutes
            .toString()
            .padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
    };

    return <div>{formatTime(remaining)}</div>;
}
