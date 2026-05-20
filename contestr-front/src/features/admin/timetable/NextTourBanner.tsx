import { AlertTriangle, Coffee, Play } from "lucide-react";
import type { TimetableView } from "@/client/types.gen";
import { AutostartToggle } from "./AutostartToggle";
import { SEGMENT_SLOT_ICON_PX } from "./segmentIcons";
import { TourStatusIcon } from "./TourStatusIcon";
import { useContestElapsed } from "./useContestElapsed";
import {
    formatDurationCompact,
    formatThroughMinutes,
    formatTourClock,
    segmentLabel,
} from "./time";

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
    const elapsed = useContestElapsed(
        view.contest_start_time,
        view.elapsed_seconds ?? 0,
    );

    if (!displaySegment && !hasPending) {
        return null;
    }

    const sectionLabel = isLastSlot ? "Последний тур" : "Следующий";
    const label = displaySegment ? segmentLabel(displaySegment) : "Слот";
    const isPause = displaySegment?.kind === "pause";
    const autostartAvailable = view.auto_start_available ?? false;
    const autostartOn = autostartAvailable && Boolean(view.auto_start_enabled);
    const hasAdvanceTarget =
        !isLastSlot &&
        (nextSegment?.status === "next" ||
            nextSegment?.status === "starting" ||
            (activeSegment != null && hasPending));

    const throughSegment =
        nextSegment && displaySegment && !isLastSlot ? displaySegment : null;
    const showThrough = throughSegment != null;
    const untilStartSeconds = throughSegment
        ? Math.max(0, throughSegment.start_time - elapsed)
        : 0;
    const throughMinutes = throughSegment
        ? Math.max(0, Math.round(untilStartSeconds / 60))
        : null;
    const showThroughWarning = throughMinutes === 0;

    return (
        <section className="tt-next-banner">
            <div className="tt-next-banner__content">
                <span className="tt-next-banner__lead">
                    <span className="tt-next-banner__label">{sectionLabel}</span>
                    {displaySegment && (
                        <span className="tt-next-banner__tour-chip">
                            {!isPause && (
                                <span className="tt-next-banner__tour-chip-icon">
                                    <TourStatusIcon
                                        status={displaySegment.status}
                                        visualState={displaySegment.status}
                                        size={SEGMENT_SLOT_ICON_PX}
                                    />
                                </span>
                            )}
                            {isPause ? (
                                <Coffee
                                    size={SEGMENT_SLOT_ICON_PX}
                                    strokeWidth={2}
                                    className={`tt-next-banner__pause-icon${displaySegment.status === "active" ? " tt-next-banner__pause-icon--active" : ""}`}
                                    aria-label={label}
                                />
                            ) : (
                                <span className="tt-next-banner__tour-chip-name">{label}</span>
                            )}
                        </span>
                    )}
                    {displaySegment && showThrough && (
                        <span className="tt-next-banner__meta-part tt-next-banner__through">
                            через{" "}
                            <strong>{formatThroughMinutes(untilStartSeconds)}</strong>
                            {showThroughWarning && (
                                <AlertTriangle
                                    size={15}
                                    className="tt-next-banner__through-warn"
                                    aria-label="Тур начинается сейчас"
                                />
                            )}
                        </span>
                    )}
                </span>
                {displaySegment && (
                    <>
                        <span className="tt-next-banner__sep" aria-hidden>
                            ·
                        </span>
                        <span className="tt-next-banner__details">
                            <span className="tt-next-banner__meta-part">
                                старт:{" "}
                                {formatTourClock(
                                    view.contest_start_time,
                                    displaySegment.start_time,
                                )}
                            </span>
                            <span className="tt-next-banner__sep" aria-hidden>
                                ·
                            </span>
                            <span className="tt-next-banner__meta-part">
                                длительность: {formatDurationCompact(displaySegment.duration)}
                            </span>
                        </span>
                    </>
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
                {hasAdvanceTarget && (
                    <button
                        type="button"
                        className="admin-icon-btn admin-primary-btn tt-start-now-btn"
                        onClick={onStartNow}
                        disabled={busy || autostartOn}
                        title={
                            autostartOn
                                ? "Отключите автозапуск, чтобы запустить слот вручную"
                                : undefined
                        }
                    >
                        <Play size={16} />
                        Запустить сейчас
                    </button>
                )}
            </div>
        </section>
    );
}
