import { useMemo, type CSSProperties } from "react";
import type { TimetableView, TourConfig } from "../client/types.gen";
import { AxisSegment } from "./AxisSegment";
import { formatClockHHMM, formatTourClock } from "./time";

type Props = {
    tours: TourConfig[];
    view: TimetableView | null;
    dirty?: boolean;
    busy?: boolean;
    onDurationChange: (index: number, duration: number) => void;
    onRemove: (index: number) => void;
    onAddTour: () => void;
};

function boundaryTicks(tours: TourConfig[]): number[] {
    const ticks = new Set<number>([0]);
    for (const tour of tours) {
        ticks.add(tour.start_time + tour.duration);
    }
    return [...ticks].sort((a, b) => a - b);
}

type AxisTickLayout = {
    seconds: number;
    showLabel: boolean;
    /** Нечётные границы — подпись над шкалой; чётные — под шкалой. */
    above: boolean;
};

/**
 * Подписи в двух рядах: сначала чередование по индексу границы,
 * затем скрытие только при наложении на той же стороне (верх / низ).
 */
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
    }
    if (isLast) {
        cls += " tt-axis__mark--end";
    }
    return cls;
}

export function ScheduleAxis({
    tours,
    view,
    dirty = false,
    busy = false,
    onDurationChange,
    onRemove,
    onAddTour,
}: Props) {
    const layout = useMemo(() => {
        if (tours.length === 0) {
            return null;
        }
        const totalEnd = tours.reduce(
            (max, t) => Math.max(max, t.start_time + t.duration),
            0,
        );
        if (totalEnd <= 0) {
            return null;
        }

        const elapsed = view?.elapsed_seconds ?? 0;
        const nowPct =
            view?.contest_start_time && elapsed > 0
                ? Math.min(100, (elapsed / totalEnd) * 100)
                : null;

        const ticks = boundaryTicks(tours);
        const tickLabels = layoutAxisTickLabels(ticks, totalEnd, 5.5);

        return { totalEnd, nowPct, tickLabels };
    }, [tours, view]);

    if (!layout) {
        return null;
    }

    const totalEnd = layout.totalEnd;
    const scaleStyle = {
        "--tt-tour-count": tours.length,
    } as CSSProperties;

    return (
        <section
            className={`tt-axis${dirty ? " tt-axis--dirty" : ""}`}
            aria-label="Временная шкала туров"
        >
            <p className="tt-axis__caption">Временная шкала туров</p>
            <div className="tt-axis__scale" style={scaleStyle}>
                <div className="tt-axis__ruler tt-axis__ruler--above">
                    {layout.tickLabels.map(({ seconds, showLabel, above }, i) => {
                        if (!above || !showLabel) {
                            return null;
                        }
                        const pct = (seconds / totalEnd) * 100;
                        const isFirst = i === 0;
                        const isLast = i === layout.tickLabels.length - 1;
                        return (
                            <div
                                key={`above-${seconds}`}
                                className={tickMarkClass(isFirst, isLast, true)}
                                style={{ left: `${pct}%` }}
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
                </div>
                <div className="tt-axis__track-wrap">
                    <div className="tt-axis__track">
                        {layout.nowPct !== null && (
                            <div
                                className="tt-axis__now"
                                style={{ left: `${layout.nowPct}%` }}
                                title={`Сейчас: ${formatClockHHMM(view?.elapsed_seconds ?? 0)}`}
                            >
                                <span className="tt-axis__now-line" />
                                <span className="tt-axis__now-label">сейчас</span>
                            </div>
                        )}
                        {tours.map((tour, index) => {
                            const leftPct = (tour.start_time / totalEnd) * 100;
                            const widthPct = (tour.duration / totalEnd) * 100;
                            const isLast = index === tours.length - 1;
                            return (
                                <AxisSegment
                                    key={index}
                                    index={index}
                                    tour={tour}
                                    meta={view?.tours_meta[index]}
                                    leftPct={leftPct}
                                    widthPct={widthPct}
                                    isLast={isLast}
                                    busy={busy}
                                    onDurationChange={onDurationChange}
                                    onRemove={onRemove}
                                    onAdd={isLast ? onAddTour : undefined}
                                    contestStartTime={view?.contest_start_time}
                                />
                            );
                        })}
                    </div>
                </div>
                <div className="tt-axis__ruler tt-axis__ruler--below">
                    {layout.tickLabels.map(({ seconds, showLabel, above }, i) => {
                        if (above || !showLabel) {
                            return null;
                        }
                        const pct = (seconds / totalEnd) * 100;
                        const isFirst = i === 0;
                        const isLast = i === layout.tickLabels.length - 1;
                        return (
                            <div
                                key={`below-${seconds}`}
                                className={tickMarkClass(isFirst, isLast, false)}
                                style={{ left: `${pct}%` }}
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
            </div>
        </section>
    );
}
