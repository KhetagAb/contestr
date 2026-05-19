import { Play } from "lucide-react";
import type { TimetableView } from "../client/types.gen";
import { statusLabel } from "./statusLabels";
import { TourStatusIcon } from "./TourStatusIcon";
import { formatDurationCompact, formatTourClock } from "./time";

type Props = {
    view: TimetableView;
    busy: boolean;
    onStartNow: () => void;
};

export function NextTourBanner({ view, busy, onStartNow }: Props) {
    const tourNumber = view.next_tour_number;
    if (!tourNumber) {
        return null;
    }

    const index = tourNumber - 1;
    const tour = view.tour_times[index];
    const meta = view.tours_meta[index];
    if (!tour || !meta) {
        return null;
    }

    return (
        <section className="tt-next-banner">
            <div className="tt-next-banner__content">
                <span className="tt-next-banner__label">Следующий</span>
                <span className="tt-next-banner__tour">
                    <span className="tt-next-banner__status-icon" title={statusLabel(meta.status)}>
                        <TourStatusIcon status={meta.status} />
                        <span className="visually-hidden">{statusLabel(meta.status)}</span>
                    </span>
                    <span className="tt-next-banner__tour-name">Тур {tourNumber}</span>
                </span>
                <span className="tt-next-banner__facts">
                    <span className="tt-next-banner__fact">
                        <span className="tt-next-banner__fact-label">старт:</span>
                        <span className="tt-next-banner__fact-value">
                            {formatTourClock(view.contest_start_time, tour.start_time)}
                        </span>
                    </span>
                    <span className="tt-next-banner__fact">
                        <span className="tt-next-banner__fact-label">длит:</span>
                        <span className="tt-next-banner__fact-value">
                            {formatDurationCompact(tour.duration)}
                        </span>
                    </span>
                </span>
            </div>
            {(meta.status === "next" || meta.status === "starting") && (
                <button
                    type="button"
                    className="admin-icon-btn admin-primary-btn tt-start-now-btn"
                    onClick={onStartNow}
                    disabled={busy}
                >
                    <Play size={16} />
                    Запустить сейчас
                </button>
            )}
        </section>
    );
}
