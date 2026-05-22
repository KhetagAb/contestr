/** Mirrors contestr/pkg/problemcode (Validate + Round). */

const VALID_CODE = /^\d+[A-Z]$/;

export function isValidProblemCode(code: string): boolean {
    return VALID_CODE.test(code);
}

/** Tour round from code (e.g. "2A" → 2). null if invalid. */
export function roundFromProblemCode(code: string): number | null {
    if (!isValidProblemCode(code)) {
        return null;
    }
    const digits = code.match(/^\d+/)?.[0];
    if (!digits) {
        return null;
    }
    const n = Number.parseInt(digits, 10);
    return n > 0 ? n : null;
}
