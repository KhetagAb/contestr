import { useCallback, useEffect, useMemo, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type {
    AdminProblemStatementItem,
    AdminProblemStatementsResponse,
    TourSettings,
} from "@/client/types.gen";
import { adminAuthHeaders } from "@/features/admin/auth/adminAuth";
import { readApiError } from "@/features/admin/timetable/errors";
import { problemCodesForContest } from "@/shared/utils/problemCodes";
import {
    classifyBulkFiles,
    formatBulkUploadReport,
} from "./problemStatementUpload";

type LoadState = "idle" | "loading" | "error";

function problemStatementsQueryKey(contestId: number) {
    return ["admin", "problem-statements", contestId] as const;
}

async function fetchAdminProblemStatements(
    contestId: number,
): Promise<AdminProblemStatementsResponse> {
    const response = await fetch(
        `/api/admin/contests/${contestId}/problem-statements`,
        { headers: adminAuthHeaders() },
    );
    if (!response.ok) {
        throw new Error(await readApiError(response));
    }
    return (await response.json()) as AdminProblemStatementsResponse;
}

async function putProblemStatementPdf(
    contestId: number,
    problemCode: string,
    file: File,
): Promise<void> {
    const response = await fetch(
        `/api/admin/contests/${contestId}/problem-statements/${encodeURIComponent(problemCode)}`,
        {
            method: "PUT",
            headers: {
                ...adminAuthHeaders(),
                "Content-Type": "application/pdf",
            },
            body: file,
        },
    );
    if (!response.ok) {
        throw new Error(await readApiError(response));
    }
}

export function useAdminProblemStatements(
    contestId: number | null,
    tourSettings: TourSettings,
    tourSlotCount: number,
) {
    const queryClient = useQueryClient();
    const enabled = contestId != null && contestId > 0;

    const query = useQuery({
        queryKey: contestId ? problemStatementsQueryKey(contestId) : ["admin", "problem-statements", "none"],
        queryFn: () => fetchAdminProblemStatements(contestId!),
        enabled,
    });

    const [busy, setBusy] = useState(false);
    const [message, setMessage] = useState("");
    const [messageKind, setMessageKind] = useState<"error" | "success" | "info">("info");

    const expectedCodes = useMemo(
        () => problemCodesForContest(tourSlotCount, tourSettings.problems_per_tour),
        [tourSlotCount, tourSettings.problems_per_tour],
    );

    const uploadedByCode = useMemo(() => {
        const map = new Map<string, AdminProblemStatementItem>();
        for (const item of query.data?.items ?? []) {
            map.set(item.problem_code, item);
        }
        return map;
    }, [query.data?.items]);

    const rows = useMemo(() => {
        const rowForCode = (code: string) => {
            const uploaded = uploadedByCode.get(code);
            return {
                problem_code: code,
                status: uploaded ? ("uploaded" as const) : ("missing" as const),
                public_url: uploaded?.public_url,
                uploaded_at: uploaded?.uploaded_at,
                size_bytes: uploaded?.size_bytes,
            };
        };
        const expectedRows = expectedCodes.map(rowForCode);
        const extraCodes = (query.data?.items ?? [])
            .map((i) => i.problem_code)
            .filter((code) => !expectedCodes.includes(code));
        return [...expectedRows, ...extraCodes.map(rowForCode)];
    }, [expectedCodes, uploadedByCode, query.data?.items]);

    const knownCodesSet = useMemo(
        () => new Set(rows.map((r) => r.problem_code)),
        [rows],
    );

    const loadState: LoadState = query.isLoading
        ? "loading"
        : query.isError
          ? "error"
          : "idle";

    const invalidate = useCallback(async () => {
        if (!contestId) return;
        await queryClient.invalidateQueries({
            queryKey: problemStatementsQueryKey(contestId),
        });
    }, [contestId, queryClient]);

    const uploadPdf = useCallback(
        async (problemCode: string, file: File) => {
            if (!contestId) return;
            setBusy(true);
            setMessage("");
            try {
                await putProblemStatementPdf(contestId, problemCode, file);
                await invalidate();
                setMessageKind("success");
                setMessage(`Загружено: ${problemCode}`);
            } catch (e) {
                setMessageKind("error");
                setMessage(e instanceof Error ? e.message : "Ошибка загрузки");
            } finally {
                setBusy(false);
            }
        },
        [contestId, invalidate],
    );

    const uploadBulk = useCallback(
        async (files: File[]) => {
            if (!contestId) return;

            const classified = classifyBulkFiles(files, knownCodesSet);
            const total = classified.toUpload.length;

            if (
                total === 0 &&
                classified.skippedInvalid.length === 0 &&
                classified.skippedUnknown.length === 0 &&
                classified.skippedNotPdf.length === 0
            ) {
                setMessageKind("info");
                setMessage("Нет файлов для загрузки.");
                return;
            }

            setBusy(true);
            setMessage("");

            let uploaded = 0;
            const failed: { name: string; reason: string }[] = [];

            for (const { code, file } of classified.toUpload) {
                try {
                    await putProblemStatementPdf(contestId, code, file);
                    uploaded += 1;
                } catch (e) {
                    failed.push({
                        name: file.name,
                        reason:
                            e instanceof Error ? e.message : "Ошибка загрузки",
                    });
                }
            }

            if (uploaded > 0) {
                await invalidate();
            }

            const report = formatBulkUploadReport({
                uploaded,
                total,
                failed,
                skippedInvalid: classified.skippedInvalid,
                skippedUnknown: classified.skippedUnknown,
                skippedNotPdf: classified.skippedNotPdf,
            });

            let kind: "error" | "success" | "info" = "success";
            if (uploaded === 0 && failed.length > 0) {
                kind = "error";
            } else if (uploaded === 0) {
                kind = "info";
            } else if (failed.length > 0 || uploaded < total) {
                kind = "info";
            }
            setMessageKind(kind);
            setMessage(report);
            setBusy(false);
        },
        [contestId, invalidate, knownCodesSet],
    );

    const deletePdf = useCallback(
        async (problemCode: string) => {
            if (!contestId) return;
            setBusy(true);
            setMessage("");
            try {
                const response = await fetch(
                    `/api/admin/contests/${contestId}/problem-statements/${encodeURIComponent(problemCode)}`,
                    {
                        method: "DELETE",
                        headers: adminAuthHeaders(),
                    },
                );
                if (!response.ok) {
                    throw new Error(await readApiError(response));
                }
                await invalidate();
                setMessageKind("success");
                setMessage(`Удалено: ${problemCode}`);
            } catch (e) {
                setMessageKind("error");
                setMessage(e instanceof Error ? e.message : "Ошибка удаления");
            } finally {
                setBusy(false);
            }
        },
        [contestId, invalidate],
    );

    useEffect(() => {
        setMessage("");
    }, [contestId]);

    return {
        rows,
        loadState,
        busy,
        message,
        messageKind,
        uploadPdf,
        uploadBulk,
        deletePdf,
        refetch: query.refetch,
    };
}
