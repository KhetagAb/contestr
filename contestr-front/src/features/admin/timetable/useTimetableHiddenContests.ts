import { useCallback, useState } from "react";

const STORAGE_KEY = "admin-timetable-hidden-contests";

function readHiddenIds(): Set<number> {
    try {
        const raw = localStorage.getItem(STORAGE_KEY);
        if (!raw) {
            return new Set();
        }
        const parsed: unknown = JSON.parse(raw);
        if (!Array.isArray(parsed)) {
            return new Set();
        }
        return new Set(
            parsed.filter((id): id is number => typeof id === "number" && Number.isInteger(id)),
        );
    } catch {
        return new Set();
    }
}

function writeHiddenIds(ids: Set<number>) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify([...ids]));
}

export function useTimetableHiddenContests() {
    const [hiddenIds, setHiddenIds] = useState(readHiddenIds);

    const isHidden = useCallback((contestId: number) => hiddenIds.has(contestId), [hiddenIds]);

    const setHidden = useCallback((contestId: number, hidden: boolean) => {
        setHiddenIds((prev) => {
            const next = new Set(prev);
            if (hidden) {
                next.add(contestId);
            } else {
                next.delete(contestId);
            }
            writeHiddenIds(next);
            return next;
        });
    }, []);

    const toggleHidden = useCallback((contestId: number) => {
        setHiddenIds((prev) => {
            const next = new Set(prev);
            if (next.has(contestId)) {
                next.delete(contestId);
            } else {
                next.add(contestId);
            }
            writeHiddenIds(next);
            return next;
        });
    }, []);

    return { isHidden, setHidden, toggleHidden };
}
