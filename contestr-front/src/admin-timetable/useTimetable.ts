import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { adminAuthHeaders } from "../adminAuth";
import type { ScheduleSlot, TimetableView } from "../client/types.gen";
import { useContests } from "../useContests";
import { readApiError } from "./errors";
import { DURATION_TEMPLATE_MINUTES } from "./durationTemplate";
import {
    buildDraftTimelineSegments,
    createSlotFromPrevious,
    pendingSlotsFingerprint,
} from "./time";

const POLL_INTERVAL_MS = 10_000;

type LoadState = "idle" | "loading" | "error";

function normalizeSlot(slot: ScheduleSlot): ScheduleSlot {
    return {
        duration: Number(slot.duration) || 0,
        kind: slot.kind === "pause" ? "pause" : "tour",
    };
}

function normalizeSlots(slots: ScheduleSlot[]): ScheduleSlot[] {
    return slots.map(normalizeSlot);
}

export function useTimetable() {
    const { contests } = useContests();
    const [contestId, setContestId] = useState<number>(0);
    const [view, setView] = useState<TimetableView | null>(null);
    const [draftPending, setDraftPending] = useState<ScheduleSlot[]>([]);
    const [savedFingerprint, setSavedFingerprint] = useState("");
    const [loadState, setLoadState] = useState<LoadState>("idle");
    const [actionLoading, setActionLoading] = useState(false);
    const [message, setMessage] = useState("");
    const [messageKind, setMessageKind] = useState<"error" | "success" | "info">("info");
    const [hasSchedule, setHasSchedule] = useState(false);
    const loadRequestId = useRef(0);

    const contest = useMemo(
        () => contests.find((c) => c.contest_id === contestId),
        [contestId, contests],
    );

    const canEditSchedule = contestId > 0;

    const resetScheduleState = useCallback(() => {
        setView(null);
        setDraftPending([]);
        setSavedFingerprint("");
        setHasSchedule(false);
        setLoadState("idle");
        setMessage("");
    }, []);

    useEffect(() => {
        if (contests.length === 0) {
            setContestId(0);
            resetScheduleState();
            return;
        }
        if (!contests.some((c) => c.contest_id === contestId)) {
            setContestId(contests[0].contest_id);
        }
    }, [contests, contestId, resetScheduleState]);

    useEffect(() => {
        if (!canEditSchedule) {
            resetScheduleState();
        }
    }, [canEditSchedule, resetScheduleState]);

    const dirty = useMemo(
        () => pendingSlotsFingerprint(draftPending) !== savedFingerprint,
        [draftPending, savedFingerprint],
    );

    const displaySegments = useMemo(() => {
        if (!view) {
            return [];
        }
        if (!dirty) {
            return view.timeline_segments ?? [];
        }
        return buildDraftTimelineSegments(view, draftPending);
    }, [view, dirty, draftPending]);

    const applyView = useCallback((next: TimetableView) => {
        const pending = normalizeSlots(next.pending_slots ?? []);
        const has = (next.timeline_segments?.length ?? 0) > 0 || pending.length > 0;
        setView(next);
        setDraftPending(pending);
        setSavedFingerprint(pendingSlotsFingerprint(pending));
        setHasSchedule(has);
    }, []);

    const loadTimetable = useCallback(
        async (silent = false) => {
            if (!Number.isInteger(contestId) || contestId <= 0) {
                return;
            }

            const requestId = ++loadRequestId.current;
            if (!silent) {
                setLoadState("loading");
                setMessage("");
            }

            const response = await fetch(`/api/admin/timetables/${contestId}`, {
                headers: adminAuthHeaders(),
            });

            if (requestId !== loadRequestId.current) {
                return;
            }

            if (!response.ok) {
                setLoadState("error");
                setMessageKind("error");
                setMessage(await readApiError(response));
                return;
            }

            const body = (await response.json()) as TimetableView;
            if (requestId !== loadRequestId.current) {
                return;
            }

            applyView(body);
            setLoadState("idle");
            if (!silent) {
                setMessage("");
            }
        },
        [applyView, contestId],
    );

    useEffect(() => {
        void loadTimetable();
    }, [loadTimetable]);

    const shouldPoll = useMemo(() => {
        if (dirty || !hasSchedule || !view?.contest_start_time) {
            return false;
        }
        return (view.pending_slots?.length ?? 0) > 0;
    }, [dirty, hasSchedule, view?.contest_start_time, view?.pending_slots?.length]);

    useEffect(() => {
        if (!shouldPoll) {
            return;
        }
        const id = window.setInterval(() => {
            void loadTimetable(true);
        }, POLL_INTERVAL_MS);
        return () => window.clearInterval(id);
    }, [loadTimetable, shouldPoll]);

    const setPendingDuration = useCallback((pendingIndex: number, duration: number) => {
        setDraftPending((current) =>
            current.map((s, i) => (i === pendingIndex ? { ...s, duration } : s)),
        );
        setMessage("");
    }, []);

    const setPendingKind = useCallback((pendingIndex: number, kind: ScheduleSlot["kind"]) => {
        setDraftPending((current) =>
            current.map((s, i) => (i === pendingIndex ? { ...s, kind } : s)),
        );
        setMessage("");
    }, []);

    const addSlot = useCallback(() => {
        if (contestId <= 0) {
            return;
        }
        setDraftPending((current) => {
            const last = current[current.length - 1];
            return [...current, createSlotFromPrevious(last)];
        });
        setMessage("");
    }, [contestId]);

    const removeSlot = useCallback((pendingIndex: number) => {
        setDraftPending((current) => current.filter((_, i) => i !== pendingIndex));
        setMessage("");
    }, []);

    const revertChanges = useCallback(() => {
        if (!view) {
            return;
        }
        setDraftPending(normalizeSlots(view.pending_slots ?? []));
        setMessage("");
    }, [view]);

    const applyDurationTemplate = useCallback(() => {
        if (contestId <= 0) {
            return;
        }
        const slots: ScheduleSlot[] = DURATION_TEMPLATE_MINUTES.map((minutes) => ({
            duration: minutes * 60,
            kind: "tour" as const,
        }));
        setDraftPending(slots);
        setMessage("");
    }, [contestId]);

    const saveTimetable = useCallback(async () => {
        setActionLoading(true);
        setMessage("");

        const payload = {
            pending_slots: draftPending,
            auto_start_enabled: view?.auto_start_enabled ?? true,
        };

        const response = await fetch(`/api/admin/timetables/${contestId}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                ...adminAuthHeaders(),
            },
            body: JSON.stringify(payload),
        });

        setActionLoading(false);

        if (!response.ok) {
            setMessageKind("error");
            setMessage(await readApiError(response));
            return false;
        }

        const body = (await response.json()) as TimetableView;
        applyView(body);
        setMessage("");
        return true;
    }, [applyView, contestId, draftPending, view?.auto_start_enabled]);

    const setAutoStartEnabled = useCallback(
        async (enabled: boolean) => {
            if (!hasSchedule) {
                return false;
            }

            setActionLoading(true);
            setMessage("");

            const payload = {
                pending_slots: draftPending,
                auto_start_enabled: enabled,
            };

            const response = await fetch(`/api/admin/timetables/${contestId}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    ...adminAuthHeaders(),
                },
                body: JSON.stringify(payload),
            });

            setActionLoading(false);

            if (!response.ok) {
                setMessageKind("error");
                setMessage(await readApiError(response));
                return false;
            }

            const body = (await response.json()) as TimetableView;
            applyView(body);
            setMessage("");
            return true;
        },
        [applyView, contestId, draftPending, hasSchedule],
    );

    const advance = useCallback(async () => {
        setActionLoading(true);
        setMessage("");

        const response = await fetch(`/api/admin/timetables/${contestId}/advance`, {
            method: "POST",
            headers: adminAuthHeaders(),
        });

        setActionLoading(false);

        if (!response.ok) {
            setMessageKind("error");
            setMessage(await readApiError(response));
            return false;
        }

        const body = (await response.json()) as TimetableView;
        applyView(body);
        setMessage("");
        return true;
    }, [applyView, contestId]);

    const busy = loadState === "loading" || actionLoading;

    return {
        contestId,
        setContestId,
        contest,
        view,
        displaySegments,
        draftPending,
        dirty,
        loadState,
        busy,
        message,
        messageKind,
        hasSchedule,
        canEditSchedule,
        loadTimetable,
        setPendingDuration,
        setPendingKind,
        addSlot,
        removeSlot,
        applyDurationTemplate,
        revertChanges,
        saveTimetable,
        setAutoStartEnabled,
        advance,
    };
}
