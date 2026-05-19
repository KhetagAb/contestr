import type { TourMeta } from "../client/types.gen";

export function statusLabel(status: TourMeta["status"]): string {
    switch (status) {
        case "started":
            return "Запущен";
        case "next":
            return "Следующий";
        case "starting":
            return "Запускается…";
        case "planned":
            return "Запланирован";
        default:
            return status;
    }
}

export function statusClass(status: TourMeta["status"]): string {
    switch (status) {
        case "started":
            return "tt-status--started";
        case "next":
            return "tt-status--next";
        case "starting":
            return "tt-status--starting";
        default:
            return "tt-status--planned";
    }
}
