import { Coffee } from "lucide-react";
import type { TimelineSegment } from "@/client/types.gen";
import { SEGMENT_SLOT_ICON_PX } from "./segmentIcons";
import { TourStatusIcon } from "./TourStatusIcon";
import { formatDuration } from "./time";
import type { TourVisualState } from "./tourVisualState";

/** Показывать ли маркер типа слота в узкой колонке (без hover). */
export function shouldShowCompactSegmentMarker(
    needsExpand: boolean,
    editing: boolean,
): boolean {
    return needsExpand && !editing;
}

export function compactSegmentMarkerTitle(
    segment: TimelineSegment,
    label: string,
): string {
    return `${label} · ${formatDuration(segment.duration)}`;
}

export function compactSlotBlockClass(needsExpand: boolean): string {
    return needsExpand ? "tt-axis__block--compact-slot" : "";
}

type Props = {
    segment: TimelineSegment;
    visualState: TourVisualState;
    label: string;
};

/** Иконка типа слота по центру узкой колонки (перерыв или тур). */
export function SegmentCompactMarker({ segment, visualState, label }: Props) {
    const title = compactSegmentMarkerTitle(segment, label);
    const isPause = segment.kind === "pause";

    return (
        <span
            className={`tt-axis__segment-marker${isPause ? " tt-axis__segment-marker--pause" : " tt-axis__segment-marker--tour"}`}
            title={title}
            aria-hidden
        >
            {isPause ? (
                <Coffee
                    size={SEGMENT_SLOT_ICON_PX}
                    strokeWidth={2}
                    className={`tt-axis__pause-icon${segment.status === "active" ? " tt-axis__pause-icon--active" : ""}`}
                />
            ) : (
                <TourStatusIcon
                    status={segment.status}
                    visualState={visualState}
                    size={SEGMENT_SLOT_ICON_PX}
                />
            )}
        </span>
    );
}
