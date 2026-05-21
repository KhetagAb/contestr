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
};

export function MiniTimeline({ segments, elapsed }: Props) {
    if (segments.length === 0) {
        return null;
    }

    const progress = computeActiveSegmentProgress(segments, elapsed);

    return (
        <div className={styles.track} role="list" aria-label="Ход расписания контеста">
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
                                    size={SEGMENT_SLOT_ICON_PX}
                                    strokeWidth={2}
                                    className={styles.pauseIcon}
                                    aria-hidden
                                />
                            ) : (
                                <>
                                    <TourStatusIcon
                                        status={segment.status}
                                        visualState={visual}
                                        size={SEGMENT_SLOT_ICON_PX}
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
