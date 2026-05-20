import { useCallback, useEffect, useMemo, useState } from "react";
import { adminAuthHeaders } from "@/features/admin/auth/adminAuth";
import type {
    CodeforcesHandleItem,
    RegisteredContestItem,
    ScoringSettings,
    TourSettings,
} from "@/client/types.gen";
import { readApiError } from "@/features/admin/timetable/errors";
import { mergeHandleDraft, parseParticipantList } from "./parseParticipantImport";

type LoadState = "idle" | "loading" | "error";

export const DEFAULT_SCORING_SETTINGS: ScoringSettings = {
    mode: "binary",
    binary_overtake_mode: "retrospective",
    full_solve_bonus: 100,
    solve_in_time_bonus: 100,
    overtake_bonus: 100,
};

export const DEFAULT_TOUR_SETTINGS: TourSettings = {
    group_size: 3,
    problems_per_tour: 2,
    group_shuffle_percent: 40,
};

function normalizeScoringSettings(settings?: Partial<ScoringSettings> | null): ScoringSettings {
    return {
        ...DEFAULT_SCORING_SETTINGS,
        ...(settings ?? {}),
        mode: settings?.mode === "partial" ? "partial" : "binary",
        binary_overtake_mode:
            settings?.binary_overtake_mode === "during_tour_only"
                ? "during_tour_only"
                : "retrospective",
        full_solve_bonus: Number(settings?.full_solve_bonus ?? DEFAULT_SCORING_SETTINGS.full_solve_bonus),
        solve_in_time_bonus: Number(
            settings?.solve_in_time_bonus ?? DEFAULT_SCORING_SETTINGS.solve_in_time_bonus,
        ),
        overtake_bonus: Number(settings?.overtake_bonus ?? DEFAULT_SCORING_SETTINGS.overtake_bonus),
    };
}

function normalizeTourSettings(settings?: Partial<TourSettings> | null): TourSettings {
    return {
        group_size: Math.max(1, Number(settings?.group_size ?? DEFAULT_TOUR_SETTINGS.group_size)),
        problems_per_tour: Math.max(
            1,
            Number(settings?.problems_per_tour ?? DEFAULT_TOUR_SETTINGS.problems_per_tour),
        ),
        group_shuffle_percent: Math.min(
            100,
            Math.max(
                0,
                Number(settings?.group_shuffle_percent ?? DEFAULT_TOUR_SETTINGS.group_shuffle_percent),
            ),
        ),
    };
}

function cleanHandleRows(rows: CodeforcesHandleItem[]): CodeforcesHandleItem[] {
    return rows
        .map((h) => ({ handle: h.handle.trim(), name: h.name.trim() }))
        .filter((h) => h.handle !== "");
}

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
    const [draftScoringSettings, setDraftScoringSettings] = useState<ScoringSettings>(
        DEFAULT_SCORING_SETTINGS,
    );
    const [draftTourSettings, setDraftTourSettings] = useState<TourSettings>(
        DEFAULT_TOUR_SETTINGS,
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

    const selectedContest = useMemo(
        () => contests.find((c) => c.contest_id === selectedContestId) ?? null,
        [contests, selectedContestId],
    );

    useEffect(() => {
        setDraftScoringSettings(normalizeScoringSettings(selectedContest?.scoring_settings));
        setDraftTourSettings(normalizeTourSettings(selectedContest?.tour_settings));
    }, [selectedContest?.contest_id, selectedContest?.scoring_settings, selectedContest?.tour_settings]);

    const addContest = useCallback(
        async (contestId: number, name?: string) => {
            setBusy(true);
            try {
                const body: {
                    contest_id: number;
                    name?: string;
                    scoring_settings: ScoringSettings;
                    tour_settings: TourSettings;
                } = {
                    contest_id: contestId,
                    scoring_settings: DEFAULT_SCORING_SETTINGS,
                    tour_settings: DEFAULT_TOUR_SETTINGS,
                };
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

    const persistHandles = useCallback(async (rows: CodeforcesHandleItem[], successMessage: string) => {
        if (selectedContestId == null) {
            return false;
        }
        setBusy(true);
        try {
            const payload = cleanHandleRows(rows);
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
            showHandlesMessage(successMessage, "success");
            return true;
        } catch (e) {
            showHandlesMessage(
                e instanceof Error ? e.message : "Не удалось сохранить маппинги",
                "error",
            );
            return false;
        } finally {
            setBusy(false);
        }
    }, [selectedContestId, showHandlesMessage]);

    const saveHandles = useCallback(async () => {
        await persistHandles(draftHandles, "Участники сохранены");
    }, [draftHandles, persistHandles]);

    const addParticipant = useCallback(async (handle: string, name: string) => {
        const h = handle.trim();
        if (!h) {
            return false;
        }

        const next = cleanHandleRows(draftHandles);
        const key = h.toLowerCase();
        const existingIndex = next.findIndex((row) => row.handle.toLowerCase() === key);
        if (existingIndex >= 0) {
            next[existingIndex] = { handle: h, name: name.trim() };
        } else {
            next.push({ handle: h, name: name.trim() });
        }

        return persistHandles(next, "Участник добавлен");
    }, [draftHandles, persistHandles]);

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

    const settingsDirty = useMemo(() => {
        return (
            JSON.stringify(normalizeScoringSettings(draftScoringSettings)) !==
                JSON.stringify(normalizeScoringSettings(selectedContest?.scoring_settings)) ||
            JSON.stringify(normalizeTourSettings(draftTourSettings)) !==
                JSON.stringify(normalizeTourSettings(selectedContest?.tour_settings))
        );
    }, [draftScoringSettings, draftTourSettings, selectedContest?.scoring_settings, selectedContest?.tour_settings]);

    const updateDraftScoringSetting = useCallback(
        <K extends keyof ScoringSettings>(field: K, value: ScoringSettings[K]) => {
            setDraftScoringSettings((prev) => normalizeScoringSettings({ ...prev, [field]: value }));
        },
        [],
    );

    const updateDraftTourSetting = useCallback(
        <K extends keyof TourSettings>(field: K, value: TourSettings[K]) => {
            setDraftTourSettings((prev) => normalizeTourSettings({ ...prev, [field]: value }));
        },
        [],
    );

    const saveContestSettings = useCallback(async () => {
        if (selectedContestId == null) {
            return;
        }
        setBusy(true);
        try {
            const payload = {
                scoring_settings: normalizeScoringSettings(draftScoringSettings),
                tour_settings: normalizeTourSettings(draftTourSettings),
            };
            const res = await fetch(`/api/admin/contests/${selectedContestId}/settings`, {
                method: "PATCH",
                headers: { ...adminAuthHeaders(), "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });
            if (!res.ok) {
                throw new Error(await readApiError(res));
            }
            const updated = (await res.json()) as RegisteredContestItem;
            setContests((prev) =>
                prev.map((contest) => (contest.contest_id === updated.contest_id ? updated : contest)),
            );
            setDraftScoringSettings(normalizeScoringSettings(updated.scoring_settings));
            setDraftTourSettings(normalizeTourSettings(updated.tour_settings));
            showHandlesMessage("Настройки сохранены", "success");
        } catch (e) {
            showHandlesMessage(e instanceof Error ? e.message : "Не удалось сохранить настройки", "error");
        } finally {
            setBusy(false);
        }
    }, [draftScoringSettings, draftTourSettings, selectedContestId, showHandlesMessage]);

    const refreshContest = useCallback(
        async (contestId: number) => {
            setBusy(true);
            try {
                const res = await fetch(`/api/admin/contests/${contestId}/refresh`, {
                    method: "POST",
                    headers: adminAuthHeaders(),
                });
                if (!res.ok) {
                    throw new Error(await readApiError(res));
                }
                const updated = (await res.json()) as RegisteredContestItem;
                setContests((prev) =>
                    prev.map((contest) =>
                        contest.contest_id === updated.contest_id ? updated : contest,
                    ),
                );
                await loadContests();
                showContestsMessage(`Данные контеста «${updated.name}» обновлены`, "success");
            } catch (e) {
                showContestsMessage(
                    e instanceof Error ? e.message : "Не удалось обновить данные контеста",
                    "error",
                );
            } finally {
                setBusy(false);
            }
        },
        [loadContests, showContestsMessage],
    );

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
        settingsDirty,
        draftScoringSettings,
        draftTourSettings,
        addContest,
        deleteContest,
        addParticipant,
        saveHandles,
        saveContestSettings,
        refreshContest,
        deleteHandle,
        updateDraftRow,
        updateDraftScoringSetting,
        updateDraftTourSetting,
        removeDraftRow,
        appendDraftRow,
        importHandlesFromList,
        showHandlesMessage,
        reloadContests: loadContests,
    };
}
