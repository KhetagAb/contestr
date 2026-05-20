import { useMemo, type CSSProperties } from "react";
import type { TimelineSegment, TimetableView } from "@/client/types.gen";
import { AxisSegment } from "./AxisSegment";
import { computeActiveSegmentProgress } from "./tourProgress";
import { canInsertPendingSlotAfter, formatTourClock, pendingInsertIndexAfter } from "./time";
import { useContestElapsed } from "./useContestElapsed";

type Props = {
    segments: TimelineSegment[];
    view: TimetableView | null;
    dirty?: boolean;
    busy?: boolean;
    onPendingDurationChange: (pendingIndex: number, duration: number) => void;
    onActiveDurationChange: (duration: number) => void | Promise<unknown>;
    onPendingKindChange: (pendingIndex: number, kind: "tour" | "pause") => void;
    onPendingRemove: (pendingIndex: number) => void;
    onAddSlotAfter: (insertIndex: number) => void;
};

function boundaryTicks(segments: TimelineSegment[]): number[] {
    const ticks = new Set<number>([0]);
    for (const segment of segments) {
        ticks.add(segment.start_time + segment.duration);
    }
    return [...ticks].sort((a, b) => a - b);
}

type AxisTickLayout = {
    seconds: number;
    showLabel: boolean;
    above: boolean;
};

function layoutAxisTickLabels(
    ticks: number[],
    totalEnd: number,
    minGapPct: number,
): AxisTickLayout[] {
    const n = ticks.length;
    if (n === 0) {
        return [];
    }

    const items = ticks.map((seconds, index) => ({
        seconds,
        index,
        pct: totalEnd > 0 ? (seconds / totalEnd) * 100 : 0,
        row: index % 2,
        showLabel: true,
    }));

    for (const row of [0, 1]) {
        const onRow = items.filter((item) => item.row === row);
        for (let r = 1; r < onRow.length; r++) {
            const prev = onRow[r - 1];
            const cur = onRow[r];
            if (!prev.showLabel) {
                continue;
            }
            if (cur.pct - prev.pct >= minGapPct) {
                continue;
            }
            const isLastOnRow = r === onRow.length - 1;
            if (!isLastOnRow) {
                cur.showLabel = false;
            }
        }
    }

    if (n > 0) {
        items[0].showLabel = true;
        items[n - 1].showLabel = true;
    }

    return items.map(({ seconds, row, showLabel }) => ({
        seconds,
        showLabel,
        above: row === 1,
    }));
}

function tickMarkClass(isFirst: boolean, isLast: boolean, above: boolean): string {
    let cls = above ? "tt-axis__mark tt-axis__mark--above" : "tt-axis__mark tt-axis__mark--below";
    if (isFirst) {
        cls += " tt-axis__mark--start";
    } else if (isLast) {
        cls += " tt-axis__mark--end";
    } else {
        cls += " tt-axis__mark--mid";
    }
    return cls;
}

function buildAxisGridTemplateColumns(segments: TimelineSegment[]): string {
    return segments
        .map((s, i) =>
            i < segments.length - 1
                ? `${s.duration}fr var(--tt-tour-gap)`
                : `${s.duration}fr`,
        )
        .join(" ");
}

function segmentGridColumn(index: number): number {
    return index * 2 + 1;
}

function tickGridPlacement(
    tickIndex: number,
    tickCount: number,
    segmentCount: number,
): Pick<CSSProperties, "gridColumn" | "justifySelf"> {
    if (tickIndex === 0) {
        return { gridColumn: 1, justifySelf: "start" };
    }
    if (tickIndex === tickCount - 1) {
        return { gridColumn: segmentCount * 2 - 1, justifySelf: "end" };
    }
    return { gridColumn: tickIndex * 2, justifySelf: "center" };
}

export function ScheduleAxis({
    segments,
    view,
    dirty = false,
    busy = false,
    onPendingDurationChange,
    onActiveDurationChange,
    onPendingKindChange,
    onPendingRemove,
    onAddSlotAfter,
}: Props) {
    const elapsed = useContestElapsed(view?.contest_start_time, view?.elapsed_seconds ?? 0);

    const layout = useMemo(() => {
        if (segments.length === 0) {
            return null;
        }
        const totalEnd = segments.reduce(
            (max, s) => Math.max(max, s.start_time + s.duration),
            0,
        );
        if (totalEnd <= 0) {
            return null;
        }

        const ticks = boundaryTicks(segments);
        const tickLabels = layoutAxisTickLabels(ticks, totalEnd, 5.5);
        const activeProgress = computeActiveSegmentProgress(segments, elapsed);

        return { tickLabels, activeProgress };
    }, [segments, elapsed]);

    if (!layout) {
        return null;
    }

    const tickCount = layout.tickLabels.length;
    const segmentCount = segments.length;
    const scaleStyle = {
        "--tt-tour-count": segmentCount,
        gridTemplateColumns: buildAxisGridTemplateColumns(segments),
    } as CSSProperties;

    return (
        <section
            className={`tt-axis${dirty ? " tt-axis--dirty" : ""}`}
            aria-label="Временная шкала туров"
        >
            <p className="tt-axis__caption">Временная шкала туров</p>
            <div className="tt-axis__scale" style={scaleStyle}>
                {layout.tickLabels.map(({ seconds, showLabel, above }, i) => {
                    if (!above || !showLabel) {
                        return null;
                    }
                    const isFirst = i === 0;
                    const isLast = i === tickCount - 1;
                    return (
                        <div
                            key={`above-${seconds}`}
                            className={tickMarkClass(isFirst, isLast, true)}
                            style={{
                                gridRow: 1,
                                ...tickGridPlacement(i, tickCount, segmentCount),
                            }}
                        >
                            <span className="tt-axis__tick">
                                {formatTourClock(view?.contest_start_time, seconds)}
                            </span>
                            <span className="tt-axis__tick-stem" aria-hidden>
                                |
                            </span>
                        </div>
                    );
                })}
                {segments.map((segment, index) => {
                    const active = layout.activeProgress;
                    const isActive = active?.index === index;
                    const canAddAfter = canInsertPendingSlotAfter(segment);
                    return (
                        <AxisSegment
                            key={`${segment.sequence ?? "p"}-${segment.pending_index ?? index}`}
                            segment={segment}
                            busy={busy}
                            onDurationChange={
                                segment.editable && segment.pending_index != null
                                    ? onPendingDurationChange
                                    : undefined
                            }
                            onActiveDurationChange={
                                segment.editable && segment.sequence != null
                                    ? onActiveDurationChange
                                    : undefined
                            }
                            onKindChange={
                                segment.editable && segment.pending_index != null
                                    ? onPendingKindChange
                                    : undefined
                            }
                            onRemove={
                                segment.editable && segment.pending_index != null
                                    ? onPendingRemove
                                    : undefined
                            }
                            onAdd={
                                canAddAfter
                                    ? () =>
                                          onAddSlotAfter(pendingInsertIndexAfter(segment))
                                    : undefined
                            }
                            contestStartTime={view?.contest_start_time}
                            progressFill={isActive ? active.fill : undefined}
                            progressColorFrom={isActive ? active.colorFrom : undefined}
                            progressColorTo={isActive ? active.colorTo : undefined}
                            gridColumn={segmentGridColumn(index)}
                        />
                    );
                })}
                {layout.tickLabels.map(({ seconds, showLabel, above }, i) => {
                    if (above || !showLabel) {
                        return null;
                    }
                    const isFirst = i === 0;
                    const isLast = i === tickCount - 1;
                    return (
                        <div
                            key={`below-${seconds}`}
                            className={tickMarkClass(isFirst, isLast, false)}
                            style={{
                                gridRow: 3,
                                ...tickGridPlacement(i, tickCount, segmentCount),
                            }}
                        >
                            <span className="tt-axis__tick-stem" aria-hidden>
                                |
                            </span>
                            <span className="tt-axis__tick">
                                {formatTourClock(view?.contest_start_time, seconds)}
                            </span>
                        </div>
                    );
                })}
            </div>
        </section>
    );
}
