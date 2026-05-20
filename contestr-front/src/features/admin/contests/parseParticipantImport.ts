import type { CodeforcesHandleItem } from "@/client/types.gen";

export type ParticipantImportResult = {
    entries: CodeforcesHandleItem[];
    skippedOtherContest: number;
    invalidLines: string[];
};

/** Строка: «id[,id…] | handle | password | name» (разделитель — |). */
export function parseParticipantList(text: string, contestId: number): ParticipantImportResult {
    const entries: CodeforcesHandleItem[] = [];
    const invalidLines: string[] = [];
    let skippedOtherContest = 0;

    const lines = text.split(/\r?\n/);
    for (let i = 0; i < lines.length; i++) {
        const raw = lines[i].trim();
        if (!raw || raw.startsWith("#")) {
            continue;
        }

        const parts = raw.split("|").map((p) => p.trim());
        if (parts.length < 4) {
            invalidLines.push(`Строка ${i + 1}: нужно 4 поля (id | хэндл | пароль | имя)`);
            continue;
        }

        const [idsPart, handle, , ...nameParts] = parts;
        const name = nameParts.join(" | ").trim();
        if (!handle) {
            invalidLines.push(`Строка ${i + 1}: пустой хэндл`);
            continue;
        }

        const ids = idsPart
            .split(",")
            .map((s) => Number.parseInt(s.trim(), 10))
            .filter((n) => Number.isInteger(n) && n > 0);

        if (ids.length === 0) {
            invalidLines.push(`Строка ${i + 1}: не распознан id контеста`);
            continue;
        }

        if (!ids.includes(contestId)) {
            skippedOtherContest += 1;
            continue;
        }

        entries.push({ handle, name });
    }

    return { entries, skippedOtherContest, invalidLines };
}

export function mergeHandleDraft(
    prev: CodeforcesHandleItem[],
    imported: CodeforcesHandleItem[],
): CodeforcesHandleItem[] {
    const byHandle = new Map<string, CodeforcesHandleItem>();

    for (const row of prev) {
        const key = row.handle.trim().toLowerCase();
        if (key) {
            byHandle.set(key, { handle: row.handle.trim(), name: row.name });
        }
    }

    for (const row of imported) {
        const key = row.handle.trim().toLowerCase();
        if (!key) {
            continue;
        }
        byHandle.set(key, { handle: row.handle.trim(), name: row.name.trim() });
    }

    return [...byHandle.values()];
}
