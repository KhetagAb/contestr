import { useCallback, useEffect, useMemo, useState } from "react";
import { adminAuthHeaders } from "../adminAuth";
import type { CodeforcesHandleItem, RegisteredContestItem } from "../client/types.gen";
import { readApiError } from "../admin-timetable/errors";
import { mergeHandleDraft, parseParticipantList } from "./parseParticipantImport";

type LoadState = "idle" | "loading" | "error";

export function useAdminContests() {
    const [contests, setContests] = useState<RegisteredContestItem[]>([]);
    const [selectedContestId, setSelectedContestId] = useState<number | null>(null);
    const [handles, setHandles] = useState<CodeforcesHandleItem[]>([]);
    const [draftHandles, setDraftHandles] = useState<CodeforcesHandleItem[]>([]);
    const [loadState, setLoadState] = useState<LoadState>("idle");
    const [handlesLoadState, setHandlesLoadState] = useState<LoadState>("idle");
    const [busy, setBusy] = useState(false);
    const [contestsMessage, setContestsMessage] = useState("");
    const [contestsMessageKind, setContestsMessageKind] = useState<"error" | "success" | "info">(
        "info",
    );
    const [handlesMessage, setHandlesMessage] = useState("");
    const [handlesMessageKind, setHandlesMessageKind] = useState<"error" | "success" | "info">(
        "info",
    );

    const showContestsMessage = useCallback(
        (text: string, kind: "error" | "success" | "info" = "info") => {
            setContestsMessage(text);
            setContestsMessageKind(kind);
        },
        [],
    );

    const showHandlesMessage = useCallback(
        (text: string, kind: "error" | "success" | "info" = "info") => {
            setHandlesMessage(text);
            setHandlesMessageKind(kind);
        },
        [],
    );

    const loadContests = useCallback(async () => {
        setLoadState("loading");
        try {
            const res = await fetch("/api/admin/contests", { headers: adminAuthHeaders() });
            if (!res.ok) {
                throw new Error(await readApiError(res));
            }
            const data = (await res.json()) as RegisteredContestItem[];
            setContests(data);
            setLoadState("idle");
            return data;
        } catch (e) {
            setLoadState("error");
            showContestsMessage(e instanceof Error ? e.message : "Ошибка загрузки", "error");
            return [];
        }
    }, [showContestsMessage]);

    const loadHandles = useCallback(
        async (contestId: number) => {
            setHandlesLoadState("loading");
            try {
                const res = await fetch(`/api/admin/contests/${contestId}/handles`, {
                    headers: adminAuthHeaders(),
                });
                if (!res.ok) {
                    throw new Error(await readApiError(res));
                }
                const data = (await res.json()) as CodeforcesHandleItem[];
                setHandles(data);
                setDraftHandles(data.map((h) => ({ ...h })));
                setHandlesLoadState("idle");
            } catch (e) {
                setHandlesLoadState("error");
                showHandlesMessage(
                    e instanceof Error ? e.message : "Ошибка загрузки маппингов",
                    "error",
                );
            }
        },
        [showHandlesMessage],
    );

    useEffect(() => {
        void loadContests();
    }, [loadContests]);

    useEffect(() => {
        if (selectedContestId != null) {
            void loadHandles(selectedContestId);
        } else {
            setHandles([]);
            setDraftHandles([]);
            setHandlesMessage("");
        }
    }, [selectedContestId, loadHandles]);

    const addContest = useCallback(
        async (contestId: number, name?: string) => {
            setBusy(true);
            try {
                const body: { contest_id: number; name?: string } = { contest_id: contestId };
                if (name?.trim()) {
                    body.name = name.trim();
                }
                const res = await fetch("/api/admin/contests", {
                    method: "POST",
                    headers: { ...adminAuthHeaders(), "Content-Type": "application/json" },
                    body: JSON.stringify(body),
                });
                if (!res.ok) {
                    throw new Error(await readApiError(res));
                }
                const created = (await res.json()) as RegisteredContestItem;
                await loadContests();
                setSelectedContestId(created.contest_id);
                showContestsMessage(`Контест «${created.name}» добавлен`, "success");
            } catch (e) {
                showContestsMessage(
                    e instanceof Error ? e.message : "Не удалось добавить контест",
                    "error",
                );
            } finally {
                setBusy(false);
            }
        },
        [loadContests, showContestsMessage],
    );

    const deleteContest = useCallback(
        async (contestId: number) => {
            setBusy(true);
            try {
                const res = await fetch(`/api/admin/contests/${contestId}`, {
                    method: "DELETE",
                    headers: adminAuthHeaders(),
                });
                if (!res.ok) {
                    throw new Error(await readApiError(res));
                }
                if (selectedContestId === contestId) {
                    setSelectedContestId(null);
                }
                await loadContests();
                showContestsMessage("Контест удалён", "success");
            } catch (e) {
                showContestsMessage(
                    e instanceof Error ? e.message : "Не удалось удалить контест",
                    "error",
                );
            } finally {
                setBusy(false);
            }
        },
        [loadContests, selectedContestId, showContestsMessage],
    );

    const saveHandles = useCallback(async () => {
        if (selectedContestId == null) {
            return;
        }
        setBusy(true);
        try {
            const payload = draftHandles.filter((h) => h.handle.trim() !== "");
            const res = await fetch(`/api/admin/contests/${selectedContestId}/handles`, {
                method: "PUT",
                headers: { ...adminAuthHeaders(), "Content-Type": "application/json" },
                body: JSON.stringify({ handles: payload }),
            });
            if (!res.ok) {
                throw new Error(await readApiError(res));
            }
            const data = (await res.json()) as CodeforcesHandleItem[];
            setHandles(data);
            setDraftHandles(data.map((h) => ({ ...h })));
            showHandlesMessage("Участники сохранены", "success");
        } catch (e) {
            showHandlesMessage(
                e instanceof Error ? e.message : "Не удалось сохранить маппинги",
                "error",
            );
        } finally {
            setBusy(false);
        }
    }, [draftHandles, selectedContestId, showHandlesMessage]);

    const deleteHandle = useCallback(
        async (handle: string) => {
            if (selectedContestId == null) {
                return;
            }
            setBusy(true);
            try {
                const res = await fetch(
                    `/api/admin/contests/${selectedContestId}/handles/${encodeURIComponent(handle)}`,
                    { method: "DELETE", headers: adminAuthHeaders() },
                );
                if (!res.ok) {
                    throw new Error(await readApiError(res));
                }
                await loadHandles(selectedContestId);
                showHandlesMessage("Маппинг удалён", "success");
            } catch (e) {
                showHandlesMessage(
                    e instanceof Error ? e.message : "Не удалось удалить маппинг",
                    "error",
                );
            } finally {
                setBusy(false);
            }
        },
        [loadHandles, selectedContestId, showHandlesMessage],
    );

    const handlesDirty = useMemo(() => {
        const normalize = (rows: CodeforcesHandleItem[]) =>
            rows
                .filter((h) => h.handle.trim() !== "")
                .map((h) => ({ handle: h.handle.trim(), name: h.name.trim() }))
                .sort((a, b) => a.handle.localeCompare(b.handle, "ru"));

        return (
            JSON.stringify(normalize(draftHandles)) !== JSON.stringify(normalize(handles))
        );
    }, [draftHandles, handles]);

    const updateDraftRow = useCallback((index: number, field: "handle" | "name", value: string) => {
        setDraftHandles((prev) =>
            prev.map((row, i) => (i === index ? { ...row, [field]: value } : row)),
        );
    }, []);

    const removeDraftRow = useCallback((index: number) => {
        setDraftHandles((prev) => prev.filter((_, i) => i !== index));
    }, []);

    const appendDraftRow = useCallback((handle: string, name: string) => {
        const h = handle.trim();
        if (!h) {
            return false;
        }
        setDraftHandles((prev) => {
            const key = h.toLowerCase();
            const idx = prev.findIndex((r) => r.handle.trim().toLowerCase() === key);
            if (idx >= 0) {
                return prev.map((row, i) =>
                    i === idx ? { handle: h, name: name.trim() } : row,
                );
            }
            return [...prev, { handle: h, name: name.trim() }];
        });
        return true;
    }, []);

    const importHandlesFromList = useCallback(
        (text: string) => {
            if (selectedContestId == null) {
                return {
                    ok: false as const,
                    entriesCount: 0,
                    skippedOtherContest: 0,
                    invalidLines: [] as string[],
                    detail: "Сначала выберите контест",
                };
            }

            const { entries, skippedOtherContest, invalidLines } = parseParticipantList(
                text,
                selectedContestId,
            );

            if (entries.length === 0 && invalidLines.length === 0) {
                return {
                    ok: false as const,
                    entriesCount: 0,
                    skippedOtherContest,
                    invalidLines,
                    detail:
                        skippedOtherContest > 0
                            ? `Нет строк для контеста ID ${selectedContestId}`
                            : "Список пуст или не распознан",
                };
            }

            if (entries.length > 0) {
                setDraftHandles((prev) => mergeHandleDraft(prev, entries));
            }

            return {
                ok: entries.length > 0,
                entriesCount: entries.length,
                skippedOtherContest,
                invalidLines,
                detail: "",
            };
        },
        [selectedContestId],
    );

    return {
        contests,
        selectedContestId,
        setSelectedContestId,
        handles,
        draftHandles,
        loadState,
        handlesLoadState,
        busy,
        contestsMessage,
        contestsMessageKind,
        handlesMessage,
        handlesMessageKind,
        handlesDirty,
        addContest,
        deleteContest,
        saveHandles,
        deleteHandle,
        updateDraftRow,
        removeDraftRow,
        appendDraftRow,
        importHandlesFromList,
        showHandlesMessage,
        reloadContests: loadContests,
    };
}
