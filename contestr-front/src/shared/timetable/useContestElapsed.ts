import { useEffect, useMemo, useRef, useState } from "react";

function elapsedFromContestStart(contestStartIso: string): number | null {
    const startMs = new Date(contestStartIso).getTime();
    if (Number.isNaN(startMs)) {
        return null;
    }
    return Math.max(0, Math.floor((Date.now() - startMs) / 1000));
}

/**
 * Секунды с момента старта контеста; обновляется раз в секунду.
 * Без contest_start_time — тикает от последнего serverElapsed с API.
 */
export function useContestElapsed(contestStartIso?: string, serverElapsed = 0): number {
    const serverAnchor = useRef({ value: serverElapsed, atMs: Date.now() });

    useEffect(() => {
        serverAnchor.current = { value: serverElapsed, atMs: Date.now() };
    }, [serverElapsed]);

    const [pulse, setPulse] = useState(0);

    useEffect(() => {
        const id = window.setInterval(() => setPulse((p) => p + 1), 1000);
        return () => window.clearInterval(id);
    }, []);

    return useMemo(() => {
        if (contestStartIso) {
            const fromStart = elapsedFromContestStart(contestStartIso);
            if (fromStart != null) {
                return fromStart;
            }
        }

        const { value, atMs } = serverAnchor.current;
        return value + Math.floor((Date.now() - atMs) / 1000);
    }, [contestStartIso, serverElapsed, pulse]);
}
