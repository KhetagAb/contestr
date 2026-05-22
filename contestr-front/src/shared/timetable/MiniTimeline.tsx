import { useEffect, useMemo, useRef } from "react";
import { Coffee } from "lucide-react";
import type { TimelineSegment } from "@/client/types.gen";
import { SEGMENT_SLOT_ICON_PX } from "./segmentIcons";
import { TourStatusIcon } from "./TourStatusIcon";
import { computeActiveSegmentProgress } from "./tourProgress";
import { miniTimelineSegmentLabel } from "./segmentTime";
import {
    resolveSegmentVisualState,
    TOUR_VISUAL_BG,
    visualStateClass,
} from "./tourVisualState";
import styles from "./MiniTimeline.module.css";

type Props = {
    segments: TimelineSegment[];
    elapsed: number;
    size?: "default" | "large";
    /** Равномерно заполнить ширину; горизонтальный скролл только если не помещается */
    fitWidth?: boolean;
};

function focusSegmentIndex(segments: TimelineSegment[]): number {
    const active = segments.findIndex((s) => s.status === "active");
    if (active >= 0) {
        return active;
    }
    const upcoming = segments.findIndex((s) =>
        ["next", "starting"].includes(s.status),
    );
    if (upcoming >= 0) {
        return upcoming;
    }
    for (let i = segments.length - 1; i >= 0; i--) {
        if (segments[i].status === "past") {
            return i;
        }
    }
    return 0;
}

const LARGE_SEGMENT_ICON_PX = 22;

export function MiniTimeline({
    segments,
    elapsed,
    size = "default",
    fitWidth = false,
}: Props) {
    const trackRef = useRef<HTMLDivElement>(null);
    const focusIndex = useMemo(() => focusSegmentIndex(segments), [segments]);

    useEffect(() => {
        const track = trackRef.current;
        if (!track || focusIndex < 0) {
            return;
        }
        if (fitWidth && track.scrollWidth <= track.clientWidth + 1) {
            track.scrollLeft = 0;
            return;
        }
        const el = track.children[focusIndex] as HTMLElement | undefined;
        if (!el) {
            return;
        }
        const left = el.offsetLeft - (track.clientWidth - el.offsetWidth) / 2;
        track.scrollLeft = Math.max(0, left);
    }, [focusIndex, segments, fitWidth]);

    if (segments.length === 0) {
        return null;
    }

    const progress = computeActiveSegmentProgress(segments, elapsed);

    const iconPx = size === "large" ? LARGE_SEGMENT_ICON_PX : SEGMENT_SLOT_ICON_PX;
    const trackClass = [
        styles.track,
        size === "large" ? styles.trackLarge : "",
        fitWidth ? styles.trackFit : "",
    ]
        .filter(Boolean)
        .join(" ");

    return (
        <div
            ref={trackRef}
            className={trackClass}
            role="list"
            aria-label="Ход расписания контеста"
        >
            {segments.map((segment, index) => {
                const visual = resolveSegmentVisualState(segment);
                const isPause = segment.kind === "pause";
                const isActive = segment.status === "active";
                const fill =
                    progress?.index === index ? progress.fill : undefined;

                const bg = TOUR_VISUAL_BG[visual];
                const style =
                    isActive && fill != null && progress
                        ? {
                              backgroundImage: `linear-gradient(90deg, ${progress.colorFrom} ${fill * 100}%, ${progress.colorTo} ${fill * 100}%)`,
                          }
                        : { backgroundColor: bg };

                const tourLabel = miniTimelineSegmentLabel(segment, segments, index);
                const pauseLabel = "Перерыв";

                return (
                    <div
                        key={`${segment.sequence ?? "p"}-${segment.pending_index ?? index}-${segment.start_time}`}
                        role="listitem"
                        className={[
                            styles.segment,
                            styles[`segment--${visual}`],
                            isPause ? styles.segmentPause : "",
                            visualStateClass(visual),
                        ].join(" ")}
                        style={style}
                        title={isPause ? pauseLabel : tourLabel}
                        aria-label={isPause ? pauseLabel : tourLabel}
                    >
                        <span className={styles.segmentInner}>
                            {isPause ? (
                                <Coffee
                                    size={iconPx}
                                    strokeWidth={2}
                                    className={styles.pauseIcon}
                                    aria-hidden
                                />
                            ) : (
                                <>
                                    <TourStatusIcon
                                        status={segment.status}
                                        visualState={visual}
                                        size={iconPx}
                                    />
                                    <span className={styles.segmentLabel}>{tourLabel}</span>
                                </>
                            )}
                        </span>
                    </div>
                );
            })}
        </div>
    );
}
