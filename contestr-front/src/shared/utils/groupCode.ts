import { md5 } from "js-md5";

/** Двухсимвольный код группы: первые 2 символа MD5 от номера (0-9, a-f). */
export function formatGroupCode(teamNumber: number): string {
    if (!Number.isFinite(teamNumber) || teamNumber <= 0) {
        return "—";
    }
    return md5(String(teamNumber)).slice(0, 2).toLowerCase();
}
