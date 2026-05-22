import type { TimelineSegment } from "@/client/types.gen";

/** Как в AxisSegment: уже этого — показываем только иконку. */
export const SEGMENT_EXPAND_MIN_PX = 80;

/** Оценка ширины колонки сегмента на шкале (px). */
export function estimateSegmentWidthPx(
    segment: TimelineSegment,
    segments: TimelineSegment[],
    scaleWidthPx: number,
): number {
    const total = segments.reduce((sum, s) => sum + s.duration, 0);
    if (total <= 0 || scaleWidthPx <= 0) {
        return 0;
    }
    return (segment.duration / total) * scaleWidthPx;
}

export function isNarrowTimelineSegment(
    segment: TimelineSegment,
    segments: TimelineSegment[],
    scaleWidthPx: number,
): boolean {
    const width = estimateSegmentWidthPx(segment, segments, scaleWidthPx);
    return width > 0 && width < SEGMENT_EXPAND_MIN_PX;
}
