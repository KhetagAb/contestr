import type { TimelineSegment } from "../client/types.gen";

export function statusLabel(status: TimelineSegment["status"]): string {
    switch (status) {
        case "past":
            return "Завершён";
        case "active":
            return "Идёт";
        case "next":
            return "Следующий";
        case "starting":
            return "Запускается…";
        case "future":
            return "Запланирован";
        default:
            return status;
    }
}
