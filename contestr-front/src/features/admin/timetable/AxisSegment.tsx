import { Coffee, Pencil, Plus, Trash2 } from "lucide-react";
import { useEffect, useRef, useState, type CSSProperties, type MouseEvent } from "react";
import type { TimelineSegment } from "@/client/types.gen";
import { confirmEditPendingDuration } from "./ConfirmDialogs";
import { TourStatusIcon } from "./TourStatusIcon";
import { resolveSegmentVisualState, visualStateClass } from "./tourVisualState";
import {
    formatDuration,
    formatDurationInputValue,
    formatTourClock,
    parseDurationInput,
    segmentLabel,
} from "./time";

const TITLE_FULL_MIN_PX = 58;
const DURATION_COMPACT_MAX_PX = 64;

type Props = {
    segment: TimelineSegment;
    isLast?: boolean;
    busy?: boolean;
    onDurationChange?: (pendingIndex: number, duration: number) => void;
    onKindChange?: (pendingIndex: number, kind: "tour" | "pause") => void;
    onRemove?: (pendingIndex: number) => void;
    onAdd?: () => void;
    contestStartTime?: string;
    progressFill?: number;
    progressColorFrom?: string;
    progressColorTo?: string;
    gridColumn: number;
};

export function AxisSegment({
    segment,
    isLast = false,
    busy = false,
    onDurationChange,
    onKindChange,
    onRemove,
    onAdd,
    contestStartTime,
    progressFill,
    progressColorFrom,
    progressColorTo,
    gridColumn,
}: Props) {
    const [editing, setEditing] = useState(false);
    const [draft, setDraft] = useState(formatDurationInputValue(segment.duration));
    const blockRef = useRef<HTMLDivElement>(null);
    const [blockWidthPx, setBlockWidthPx] = useState(0);

    useEffect(() => {
        if (!editing) {
            setDraft(formatDurationInputValue(segment.duration));
        }
    }, [segment.duration, editing]);

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
    }, [segment.duration]);

    const visualState = resolveSegmentVisualState(segment);
    const label = segmentLabel(segment);
    const narrowTitle = blockWidthPx > 0 && blockWidthPx < TITLE_FULL_MIN_PX;
    const compact = blockWidthPx > 0 && blockWidthPx < DURATION_COMPACT_MAX_PX;
    const canEdit = segment.editable && segment.pending_index != null;
    const canToggleKind = canEdit && Boolean(onKindChange);
    const canDelete = canEdit;
    const actionSlots =
        (isLast && onAdd ? 1 : 0) +
        (canToggleKind ? 1 : 0) +
        (canEdit ? 1 : 0) +
        (canDelete ? 1 : 0);
    const displayLabel =
        narrowTitle && segment.kind === "tour" && segment.round
            ? String(segment.round)
            : label;

    const inProgress = progressFill !== undefined && progressColorFrom && progressColorTo;
    const progressStyle: CSSProperties | undefined = inProgress
        ? ({
              "--tt-fill": progressFill,
              "--tt-from": progressColorFrom,
              "--tt-to": progressColorTo,
          } as CSSProperties)
        : undefined;

    const elapsedInSegment = inProgress
        ? Math.round(progressFill * segment.duration)
        : 0;
    const blockTitle = inProgress
        ? `${label}: прошло ${formatDuration(elapsedInSegment)} из ${formatDuration(segment.duration)}`
        : `${label}: ${formatTourClock(contestStartTime, segment.start_time)} — ${formatDuration(segment.duration)}`;

    const commitDuration = () => {
        const seconds = parseDurationInput(draft);
        if (seconds === null || seconds <= 0) {
            setDraft(formatDurationInputValue(segment.duration));
            setEditing(false);
            return;
        }
        if (!confirmEditPendingDuration()) {
            setDraft(formatDurationInputValue(segment.duration));
            setEditing(false);
            return;
        }
        if (segment.pending_index != null && onDurationChange) {
            onDurationChange(segment.pending_index, seconds);
        }
        setEditing(false);
    };

    const startDurationEdit = (e: MouseEvent) => {
        e.stopPropagation();
        if (!canEdit) {
            return;
        }
        setEditing(true);
    };

    const togglePause = (e: MouseEvent) => {
        e.stopPropagation();
        if (segment.pending_index == null || !onKindChange) {
            return;
        }
        const next = segment.kind === "pause" ? "tour" : "pause";
        onKindChange(segment.pending_index, next);
    };

    return (
        <div
            ref={blockRef}
            className={`tt-axis__block ${visualStateClass(visualState)}${segment.kind === "pause" ? " tt-axis__block--pause" : ""}${editing ? " tt-axis__block--editing" : ""}${inProgress ? " tt-axis__block--in-progress" : ""}`}
            style={{ gridColumn, gridRow: 2, ...progressStyle }}
            title={blockTitle}
        >
            <div className="tt-axis__block-inner" data-actions={String(Math.min(actionSlots, 4))}>
                <div className="tt-axis__block-head">
                    {segment.kind !== "pause" && (
                        <span
                            className={`tt-axis__status-icon${compact ? " tt-axis__status-icon--compact" : ""}`}
                        >
                            <TourStatusIcon status={segment.status} visualState={visualState} />
                        </span>
                    )}
                    <span
                        className={`tt-axis__block-title${narrowTitle && segment.kind === "tour" ? " tt-axis__block-title--num" : ""}${segment.kind === "pause" ? " tt-axis__block-title--pause" : ""}`}
                        aria-label={label}
                    >
                        {segment.kind === "pause" ? (
                            <Coffee
                                size={compact ? 10 : 12}
                                className={`tt-axis__pause-icon${segment.status === "active" ? " tt-axis__pause-icon--active" : ""}`}
                                aria-hidden
                            />
                        ) : (
                            displayLabel
                        )}
                    </span>
                </div>

                <div
                    className={`tt-axis__block-duration${compact ? " tt-axis__block-duration--compact" : ""}`}
                >
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
                                        setDraft(formatDurationInputValue(segment.duration));
                                        setEditing(false);
                                    }
                                }}
                                onClick={(e) => e.stopPropagation()}
                                autoFocus
                                aria-label={`Длительность: ${label}`}
                            />
                            {" мин."}
                        </span>
                    ) : (
                        <span className="tt-axis__duration-text">{formatDuration(segment.duration)}</span>
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
                            aria-label="Добавить слот"
                        >
                            <Plus size={10} />
                        </button>
                    )}
                    {canToggleKind && !editing && (
                        <button
                            type="button"
                            className={`tt-axis__icon-btn tt-axis__icon-btn--pause${segment.kind === "pause" ? " tt-axis__icon-btn--pause-on" : ""}`}
                            onClick={togglePause}
                            aria-label={segment.kind === "pause" ? "Сделать туром" : "Сделать паузой"}
                            aria-pressed={segment.kind === "pause"}
                        >
                            <Coffee size={10} />
                        </button>
                    )}
                    {canEdit && !editing && (
                        <button
                            type="button"
                            className="tt-axis__icon-btn tt-duration-edit"
                            onClick={startDurationEdit}
                            aria-label={`Изменить длительность: ${label}`}
                        >
                            <Pencil size={10} />
                        </button>
                    )}
                    {canDelete && onRemove && segment.pending_index != null && (
                        <button
                            type="button"
                            className="tt-axis__icon-btn tt-axis__icon-btn--danger"
                            onClick={(e) => {
                                e.stopPropagation();
                                onRemove(segment.pending_index!);
                            }}
                            aria-label={`Удалить ${label}`}
                        >
                            <Trash2 size={10} />
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}
