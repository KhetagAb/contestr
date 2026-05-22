import { isValidProblemCode } from "@/shared/utils/problemCode";

/** Regatta problem codes for rounds 1..tourCount (e.g. 1A, 1B, 2A). */
export function problemCodesForContest(
    tourCount: number,
    problemsPerTour: number,
): string[] {
    const codes: string[] = [];
    const n = Math.max(0, tourCount);
    const p = Math.max(1, problemsPerTour);
    for (let round = 1; round <= n; round++) {
        for (let i = 0; i < p; i++) {
            const code = `${round}${String.fromCharCode(65 + i)}`;
            if (isValidProblemCode(code)) {
                codes.push(code);
            }
        }
    }
    return codes;
}
