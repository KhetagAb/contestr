import type { ProblemResult } from "@/client";

export function getTaskCellKind(problem: ProblemResult | undefined) {
    if (!problem) {
        return { kind: "empty" as const };
    }
    if (problem.score > 0) {
        return { kind: "solved" as const, problem };
    }
    if (problem.score < 0) {
        return { kind: "failed" as const, count: Math.abs(problem.score) };
    }
    return { kind: "none" as const };
}
