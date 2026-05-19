import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { adminAuthHeaders } from "../adminAuth";
import type { TimetableView, TourConfig } from "../client/types.gen";
import { CONTESTS } from "../consts";
import { readApiError } from "./errors";
import { DURATION_TEMPLATE_MINUTES } from "./durationTemplate";
import {
    createTourFromPrevious,
    rebuildChain,
    toursFingerprint,
} from "./time";

const POLL_INTERVAL_MS = 10_000;

type LoadState = "idle" | "loading" | "error";

function normalizeTour(tour: TourConfig): TourConfig {
    return {
        start_time: Number(tour.start_time) || 0,
        duration: Number(tour.duration) || 0,
        started: Boolean(tour.started),
    };
}

function normalizeTours(tours: TourConfig[]): TourConfig[] {
    return rebuildChain(tours.map(normalizeTour));
}

export function useTimetable() {
    const defaultContestId = CONTESTS[0]?.id ?? 0;
    const [contestId, setContestId] = useState<number>(defaultContestId);
    const [view, setView] = useState<TimetableView | null>(null);
    const [draftTours, setDraftTours] = useState<TourConfig[]>([]);
    const [savedTours, setSavedTours] = useState<TourConfig[]>([]);
    const [savedFingerprint, setSavedFingerprint] = useState("");
    const [loadState, setLoadState] = useState<LoadState>("idle");
    const [actionLoading, setActionLoading] = useState(false);
    const [message, setMessage] = useState("");
    const [messageKind, setMessageKind] = useState<"error" | "success" | "info">("info");
    const [hasSchedule, setHasSchedule] = useState(false);
    const loadRequestId = useRef(0);

    const contest = useMemo(
        () => CONTESTS.find((c) => c.id === contestId),
        [contestId],
    );

    const dirty = useMemo(
        () => toursFingerprint(draftTours) !== savedFingerprint,
        [draftTours, savedFingerprint],
    );

    const applyView = useCallback((next: TimetableView) => {
        const tours = normalizeTours(next.tour_times ?? []);
        const hasTours = tours.length > 0;
        setView(next);
        setDraftTours(tours);
        setSavedTours(tours);
        setSavedFingerprint(toursFingerprint(tours));
        setHasSchedule(hasTours);
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
        return draftTours.some((t) => !t.started);
    }, [dirty, draftTours, hasSchedule, view?.contest_start_time]);

    useEffect(() => {
        if (!shouldPoll) {
            return;
        }
        const id = window.setInterval(() => {
            void loadTimetable(true);
        }, POLL_INTERVAL_MS);
        return () => window.clearInterval(id);
    }, [loadTimetable, shouldPoll]);

    const setDuration = useCallback(
        (index: number, duration: number) => {
            setDraftTours((current) => {
                const next = current.map((t, i) =>
                    i === index ? { ...t, duration } : { ...t },
                );
                return rebuildChain(next);
            });
            setMessage("");
        },
        [],
    );

    const addTour = useCallback(() => {
        setDraftTours((current) => {
            const last = current[current.length - 1];
            return rebuildChain([...current, createTourFromPrevious(last)]);
        });
        setMessage("");
    }, []);

    const removeTour = useCallback((index: number) => {
        setDraftTours((current) => rebuildChain(current.filter((_, i) => i !== index)));
        setMessage("");
    }, []);

    const revertChanges = useCallback(() => {
        setDraftTours(normalizeTours(savedTours.map((t) => ({ ...t }))));
        setMessage("");
    }, [savedTours]);

    const applyDurationTemplate = useCallback(() => {
        const tours: TourConfig[] = [];
        for (const minutes of DURATION_TEMPLATE_MINUTES) {
            const prev = tours[tours.length - 1];
            tours.push({
                start_time: prev ? prev.start_time + prev.duration : 0,
                duration: minutes * 60,
                started: false,
            });
        }
        setDraftTours(rebuildChain(tours));
        setMessage("");
    }, []);

    const saveTimetable = useCallback(async () => {
        setActionLoading(true);
        setMessage("");

        const payload = {
            tour_durations: draftTours.map((t) => t.duration),
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
    }, [applyView, contestId, draftTours, view?.auto_start_enabled]);

    const setAutoStartEnabled = useCallback(
        async (enabled: boolean) => {
            if (!hasSchedule) {
                return false;
            }

            setActionLoading(true);
            setMessage("");

            const payload = {
                tour_durations: draftTours.map((t) => t.duration),
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
        [applyView, contestId, draftTours, hasSchedule],
    );

    const startNextTour = useCallback(async () => {
        const tourNumber = view?.next_tour_number;
        if (!tourNumber) {
            setMessageKind("error");
            setMessage("Нет тура для запуска.");
            return false;
        }

        setActionLoading(true);
        setMessage("");

        const response = await fetch(
            `/api/admin/timetables/${contestId}/tours/${tourNumber}/start`,
            { method: "POST", headers: adminAuthHeaders() },
        );

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
    }, [applyView, contestId, view?.next_tour_number]);

    const busy = loadState === "loading" || actionLoading;

    return {
        contestId,
        setContestId,
        contest,
        view,
        draftTours,
        dirty,
        loadState,
        busy,
        message,
        messageKind,
        hasSchedule,
        loadTimetable,
        setDuration,
        addTour,
        removeTour,
        applyDurationTemplate,
        revertChanges,
        saveTimetable,
        setAutoStartEnabled,
        startNextTour,
    };
}
