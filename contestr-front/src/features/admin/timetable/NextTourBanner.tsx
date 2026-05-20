import { Coffee, Play } from "lucide-react";
import type { TimetableView } from "@/client/types.gen";
import { AutostartToggle } from "./AutostartToggle";
import { TourStatusIcon } from "./TourStatusIcon";
import { formatDurationCompact, formatTourClock, segmentLabel } from "./time";

type Props = {
    view: TimetableView;
    busy: boolean;
    onStartNow: () => void;
    showAutostart?: boolean;
    onAutoStartChange?: (enabled: boolean) => void;
};

export function NextTourBanner({
    view,
    busy,
    onStartNow,
    showAutostart = false,
    onAutoStartChange,
}: Props) {
    const nextSegment = view.timeline_segments?.find(
        (s) => s.status === "next" || s.status === "starting",
    );
    const activeSegment = view.timeline_segments?.find((s) => s.status === "active");
    const hasPending = (view.pending_slots?.length ?? 0) > 0;
    const isLastSlot = activeSegment != null && !hasPending && !nextSegment;
    const displaySegment = nextSegment ?? (isLastSlot ? activeSegment : undefined);

    if (!displaySegment && !hasPending) {
        return null;
    }

    const sectionLabel = isLastSlot ? "Последний тур" : "Следующий";
    const label = displaySegment ? segmentLabel(displaySegment) : "Слот";
    const isPause = displaySegment?.kind === "pause";
    const autostartAvailable = view.auto_start_available ?? false;
    const autostartOn = autostartAvailable && Boolean(view.auto_start_enabled);
    const canAdvance =
        !autostartOn &&
        !isLastSlot &&
        (nextSegment?.status === "next" ||
            nextSegment?.status === "starting" ||
            (activeSegment != null && hasPending));

    return (
        <section className="tt-next-banner">
            <div className="tt-next-banner__content">
                <span className="tt-next-banner__label">{sectionLabel}</span>
                <span className="tt-next-banner__tour">
                    {displaySegment && !isPause && (
                        <span className="tt-next-banner__status-icon">
                            <TourStatusIcon
                                status={displaySegment.status}
                                visualState={displaySegment.status}
                            />
                        </span>
                    )}
                    <span className="tt-next-banner__tour-name">
                        {isPause ? (
                            <Coffee
                                size={16}
                                className={`tt-next-banner__pause-icon${displaySegment.status === "active" ? " tt-next-banner__pause-icon--active" : ""}`}
                                aria-label={label}
                            />
                        ) : (
                            label
                        )}
                    </span>
                </span>
                {displaySegment && (
                    <span className="tt-next-banner__facts">
                        <span className="tt-next-banner__fact">
                            <span className="tt-next-banner__fact-label">старт:</span>
                            <span className="tt-next-banner__fact-value">
                                {formatTourClock(view.contest_start_time, displaySegment.start_time)}
                            </span>
                        </span>
                        <span className="tt-next-banner__fact">
                            <span className="tt-next-banner__fact-label">длительность:</span>
                            <span className="tt-next-banner__fact-value">
                                {formatDurationCompact(displaySegment.duration)}
                            </span>
                        </span>
                    </span>
                )}
            </div>
            <div className="tt-next-banner__actions">
                {showAutostart && onAutoStartChange && (
                    <AutostartToggle
                        enabled={autostartOn}
                        available={autostartAvailable}
                        busy={busy}
                        onChange={onAutoStartChange}
                    />
                )}
                {canAdvance && (
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
            </div>
        </section>
    );
}
