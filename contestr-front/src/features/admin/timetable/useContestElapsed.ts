import { useEffect, useState } from "react";

/** Секунды с момента старта контеста; обновляется раз в секунду. */
export function useContestElapsed(contestStartIso?: string, serverElapsed = 0): number {
    const [elapsed, setElapsed] = useState(serverElapsed);

    useEffect(() => {
        if (!contestStartIso) {
            setElapsed(0);
            return;
        }
        const startMs = new Date(contestStartIso).getTime();
        if (Number.isNaN(startMs)) {
            setElapsed(0);
            return;
        }
        const tick = () => {
            setElapsed(Math.max(0, Math.floor((Date.now() - startMs) / 1000)));
        };
        tick();
        const id = window.setInterval(tick, 1000);
        return () => window.clearInterval(id);
    }, [contestStartIso]);

    useEffect(() => {
        if (contestStartIso) {
            return;
        }
        setElapsed(serverElapsed);
    }, [contestStartIso, serverElapsed]);

    return elapsed;
}
