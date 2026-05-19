import { Pencil, Plus, Trash2 } from "lucide-react";
import { useEffect, useRef, useState, type MouseEvent } from "react";

/** Меньше — в заголовке только номер тура (по фактической ширине блока, не по % длительности). */
const TITLE_FULL_MIN_PX = 58;
const DURATION_COMPACT_MAX_PX = 64;
import type { TourConfig, TourMeta } from "../client/types.gen";
import { confirmEditStartedDuration } from "./ConfirmDialogs";
import { statusClass } from "./statusLabels";
import { TourStatusIcon } from "./TourStatusIcon";
import {
    formatDuration,
    formatDurationInputValue,
    formatTourClock,
    parseDurationInput,
} from "./time";

type Props = {
    index: number;
    tour: TourConfig;
    meta?: TourMeta;
    leftPct: number;
    widthPct: number;
    isLast?: boolean;
    busy?: boolean;
    onDurationChange: (index: number, duration: number) => void;
    onRemove: (index: number) => void;
    onAdd?: () => void;
    contestStartTime?: string;
};

export function AxisSegment({
    index,
    tour,
    meta,
    leftPct,
    widthPct,
    isLast = false,
    busy = false,
    onDurationChange,
    onRemove,
    onAdd,
    contestStartTime,
}: Props) {
    const [editing, setEditing] = useState(false);
    const [draft, setDraft] = useState(formatDurationInputValue(tour.duration));
    const blockRef = useRef<HTMLDivElement>(null);
    const [blockWidthPx, setBlockWidthPx] = useState(0);

    useEffect(() => {
        if (!editing) {
            setDraft(formatDurationInputValue(tour.duration));
        }
    }, [tour.duration, editing]);

    useEffect(() => {
        const el = blockRef.current;
        if (!el) {
            return;
        }
        const sync = () => setBlockWidthPx(el.getBoundingClientRect().width);
        sync();
        const ro = new ResizeObserver(sync);
        ro.observe(el);
        return () => ro.disconnect();
    }, [leftPct, widthPct]);

    const status = meta?.status ?? (tour.started ? "started" : "planned");
    const canDelete = !tour.started;
    const narrowTitle = blockWidthPx > 0 && blockWidthPx < TITLE_FULL_MIN_PX;
    const compact = blockWidthPx > 0 && blockWidthPx < DURATION_COMPACT_MAX_PX;
    const actionSlots = (isLast && onAdd ? 1 : 0) + 1 + (canDelete ? 1 : 0);
    const tourLabel = narrowTitle ? String(index + 1) : `Тур ${index + 1}`;

    const commitDuration = () => {
        const seconds = parseDurationInput(draft);
        if (seconds === null || seconds <= 0) {
            setDraft(formatDurationInputValue(tour.duration));
            setEditing(false);
            return;
        }
        if (tour.started && !confirmEditStartedDuration(index + 1)) {
            setDraft(formatDurationInputValue(tour.duration));
            setEditing(false);
            return;
        }
        onDurationChange(index, seconds);
        setEditing(false);
    };

    const startEdit = (e: MouseEvent) => {
        e.stopPropagation();
        setDraft(formatDurationInputValue(tour.duration));
        setEditing(true);
    };

    return (
        <div
            ref={blockRef}
            className={`tt-axis__block ${statusClass(status)}${editing ? " tt-axis__block--editing" : ""}`}
            style={{ left: `${leftPct}%`, width: `${widthPct}%` }}
            title={`Тур ${index + 1}: ${formatTourClock(contestStartTime, tour.start_time)} — ${formatDuration(tour.duration)}`}
        >
            <div
                className="tt-axis__block-inner"
                data-actions={String(Math.min(actionSlots, 3))}
            >
                <div className="tt-axis__block-head">
                    <span
                        className={`tt-axis__status-icon${compact ? " tt-axis__status-icon--compact" : ""}`}
                    >
                        <TourStatusIcon status={status} />
                    </span>
                    <span
                        className={`tt-axis__block-title${narrowTitle ? " tt-axis__block-title--num" : ""}`}
                        aria-label={`Тур ${index + 1}`}
                    >
                        {tourLabel}
                    </span>
                </div>

                <div className={`tt-axis__block-duration${compact ? " tt-axis__block-duration--compact" : ""}`}>
                    {editing ? (
                        <span className="tt-axis__duration-text tt-axis__duration-text--editing">
                            <input
                                type="text"
                                inputMode="numeric"
                                className="tt-axis__duration-value"
                                size={Math.max(2, draft.length || 1)}
                                maxLength={4}
                                value={draft}
                                onChange={(e) => setDraft(e.target.value.replace(/\D/g, ""))}
                                onBlur={commitDuration}
                                onFocus={(e) => e.currentTarget.select()}
                                onKeyDown={(e) => {
                                    e.stopPropagation();
                                    if (e.key === "Enter") {
                                        commitDuration();
                                    }
                                    if (e.key === "Escape") {
                                        setDraft(formatDurationInputValue(tour.duration));
                                        setEditing(false);
                                    }
                                }}
                                onClick={(e) => e.stopPropagation()}
                                autoFocus
                                aria-label={`Длительность тура ${index + 1} в минутах`}
                            />
                            {" мин."}
                        </span>
                    ) : (
                        <span className="tt-axis__duration-text">{formatDuration(tour.duration)}</span>
                    )}
                </div>

                <div className="tt-axis__block-actions">
                    {isLast && onAdd && (
                        <button
                            type="button"
                            className="tt-axis__icon-btn"
                            disabled={busy}
                            onClick={(e) => {
                                e.stopPropagation();
                                onAdd();
                            }}
                            aria-label="Добавить тур"
                        >
                            <Plus size={10} />
                        </button>
                    )}
                    {!editing && (
                        <button
                            type="button"
                            className="tt-axis__icon-btn tt-duration-edit"
                            onClick={startEdit}
                            aria-label={`Изменить длительность тура ${index + 1}`}
                        >
                            <Pencil size={10} />
                        </button>
                    )}
                    {canDelete && (
                        <button
                            type="button"
                            className="tt-axis__icon-btn tt-axis__icon-btn--danger"
                            onClick={(e) => {
                                e.stopPropagation();
                                onRemove(index);
                            }}
                            aria-label={`Удалить тур ${index + 1}`}
                        >
                            <Trash2 size={10} />
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}
