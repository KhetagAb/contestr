const STORAGE_PREFIX = "contestr:follow:";

export function followStorageKey(contestId: number): string {
    return `${STORAGE_PREFIX}${contestId}`;
}

export function readFollowedParticipantId(contestId: number): string | null {
    if (!Number.isFinite(contestId) || contestId <= 0) {
        return null;
    }
    try {
        const raw = localStorage.getItem(followStorageKey(contestId));
        if (!raw?.trim()) {
            return null;
        }
        return raw.trim();
    } catch {
        return null;
    }
}

export function writeFollowedParticipantId(
    contestId: number,
    participantId: string | null,
): void {
    if (!Number.isFinite(contestId) || contestId <= 0) {
        return;
    }
    try {
        const key = followStorageKey(contestId);
        if (!participantId?.trim()) {
            localStorage.removeItem(key);
            return;
        }
        localStorage.setItem(key, participantId.trim());
    } catch {
        /* ignore quota / private mode */
    }
}
