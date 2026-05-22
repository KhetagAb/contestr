import { Coffee, Plus, Trash2 } from "lucide-react";
import { useEffect, useRef, useState, type CSSProperties, type MouseEvent } from "react";
import type { TimelineSegment } from "@/client/types.gen";
import {
    compactSlotBlockClass,
    SegmentCompactMarker,
    shouldShowCompactSegmentMarker,
} from "./SegmentCompactMarker";
import { SEGMENT_SLOT_ICON_PX } from "./segmentIcons";
import { TourStatusIcon } from "./TourStatusIcon";
import { resolveSegmentVisualState, visualStateClass } from "./tourVisualState";
import {
    formatDuration,
    formatDurationInputValue,
    formatTourClock,
    parseDurationInput,
    segmentLabel,
} from "./time";

type Props = {
    segment: TimelineSegment;
    busy?: boolean;
    onDurationChange?: (pendingIndex: number, duration: number) => void;
    onActiveDurationChange?: (duration: number) => void | Promise<unknown>;
    onKindChange?: (pendingIndex: number, kind: "tour" | "pause") => void;
    onRemove?: (pendingIndex: number) => void;
    onAdd?: () => void;
    contestStartTime?: string;
    progressFill?: number;
    progressColorFrom?: string;
    progressColorTo?: string;
    gridColumn: number;
    isNarrowSlot?: boolean;
};

export function AxisSegment({
    segment,
    busy = false,
    onDurationChange,
    onActiveDurationChange,
    onKindChange,
    onRemove,
    onAdd,
    contestStartTime,
    progressFill,
    progressColorFrom,
    progressColorTo,
    gridColumn,
    isNarrowSlot = false,
}: Props) {
    const [editing, setEditing] = useState(false);
    const [draft, setDraft] = useState(formatDurationInputValue(segment.duration));
    const blockRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (!editing) {
            setDraft(formatDurationInputValue(segment.duration));
        }
    }, [segment.duration, editing]);

    const visualState = resolveSegmentVisualState(segment);
    const label = segmentLabel(segment);
    const canEditPending = Boolean(
        segment.editable && segment.pending_index != null && onDurationChange,
    );
    const canEditActive = Boolean(
        segment.editable && segment.sequence != null && onActiveDurationChange,
    );
    const canEdit = canEditPending || canEditActive;
    const canToggleKind = canEditPending && Boolean(onKindChange);
    const canDelete = canEditPending;
    const hasActions =
        Boolean(onAdd) || canToggleKind || (canDelete && onRemove);
    /** Узкий слот: иконка в колонке, раскрытие по hover (и прошедшие, и запланированные). */
    const needsExpand = isNarrowSlot && !editing;
    const showCompactMarker = shouldShowCompactSegmentMarker(isNarrowSlot, editing);

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
    const blockAriaLabel = inProgress
        ? `${label}: прошло ${formatDuration(elapsedInSegment)} из ${formatDuration(segment.duration)}`
        : `${label}: ${formatTourClock(contestStartTime, segment.start_time)} — ${formatDuration(segment.duration)}`;

    const commitDuration = () => {
        const seconds = parseDurationInput(draft);
        if (seconds === null || seconds <= 0) {
            setDraft(formatDurationInputValue(segment.duration));
            setEditing(false);
            return;
        }
        if (segment.pending_index != null && onDurationChange) {
            onDurationChange(segment.pending_index, seconds);
        } else if (segment.sequence != null && onActiveDurationChange) {
            void onActiveDurationChange(seconds);
        }
        setEditing(false);
    };

    const startDurationEdit = (e?: MouseEvent) => {
        e?.stopPropagation();
        if (!canEdit || busy) {
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

    const surfaceClass = [
        visualStateClass(visualState),
        inProgress ? "tt-axis__block--in-progress" : "",
    ]
        .filter(Boolean)
        .join(" ");

    const blockClassName = [
        "tt-axis__block",
        needsExpand ? "" : surfaceClass,
        editing ? "tt-axis__block--editing" : "",
        needsExpand ? "tt-axis__block--expandable" : "",
        compactSlotBlockClass(isNarrowSlot),
        hasActions ? "tt-axis__block--has-actions" : "",
    ]
        .filter(Boolean)
        .join(" ");

    const innerClassName = ["tt-axis__block-inner", needsExpand ? surfaceClass : ""]
        .filter(Boolean)
        .join(" ");

    return (
        <div
            ref={blockRef}
            className={blockClassName}
            style={{ gridColumn, gridRow: 2, ...progressStyle }}
            aria-label={blockAriaLabel}
        >
            <div className={innerClassName}>
                {showCompactMarker && (
                    <SegmentCompactMarker
                        segment={segment}
                        visualState={visualState}
                        label={label}
                    />
                )}
                <div
                    className={`tt-axis__block-main${canEdit ? " tt-axis__block-main--editable" : ""}`}
                    onDoubleClick={canEdit && !editing ? startDurationEdit : undefined}
                    onMouseDown={
                        canEdit && !editing
                            ? (e) => {
                                  if (e.detail > 1) {
                                      e.preventDefault();
                                  }
                              }
                            : undefined
                    }
                    title={canEdit ? "Двойной щелчок — изменить длительность" : undefined}
                >
                    <div className="tt-axis__block-head">
                        <span className="tt-axis__status-icon">
                            {segment.kind === "pause" ? (
                                <Coffee
                                    size={SEGMENT_SLOT_ICON_PX}
                                    strokeWidth={2}
                                    className={`tt-axis__pause-icon${segment.status === "active" ? " tt-axis__pause-icon--active" : ""}`}
                                    aria-label={label}
                                />
                            ) : (
                                <TourStatusIcon
                                    status={segment.status}
                                    visualState={visualState}
                                    size={SEGMENT_SLOT_ICON_PX}
                                />
                            )}
                        </span>
                        {segment.kind !== "pause" && (
                            <span className="tt-axis__block-title" title={label}>
                                {label}
                            </span>
                        )}
                    </div>

                    <div className="tt-axis__block-duration">
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
                            <span className="tt-axis__duration-text">
                                {formatDuration(segment.duration)}
                            </span>
                        )}
                    </div>
                </div>

                <div
                    className={`tt-axis__block-actions${hasActions ? "" : " tt-axis__block-actions--reserved"}`}
                    aria-hidden={!hasActions}
                >
                    {hasActions && onAdd && (
                        <button
                            type="button"
                            className="tt-axis__icon-btn"
                            disabled={busy}
                            onClick={(e) => {
                                e.stopPropagation();
                                onAdd();
                            }}
                            aria-label="Добавить слот после этого"
                            title="Добавить слот после этого"
                        >
                            <Plus size={16} />
                        </button>
                    )}
                    {hasActions && canToggleKind && !editing && (
                        <button
                            type="button"
                            className={`tt-axis__icon-btn tt-axis__icon-btn--pause${
                                segment.kind === "pause" && segment.status === "active"
                                    ? " tt-axis__icon-btn--pause-on"
                                    : segment.kind === "pause"
                                      ? " tt-axis__icon-btn--pause-set"
                                      : ""
                            }`}
                            onClick={togglePause}
                            aria-label={segment.kind === "pause" ? "Сделать туром" : "Сделать паузой"}
                            aria-pressed={segment.kind === "pause"}
                        >
                            <Coffee size={16} />
                        </button>
                    )}
                    {hasActions && canDelete && onRemove && segment.pending_index != null && (
                        <button
                            type="button"
                            className="tt-axis__icon-btn tt-axis__icon-btn--danger"
                            onClick={(e) => {
                                e.stopPropagation();
                                onRemove(segment.pending_index!);
                            }}
                            aria-label={`Удалить ${label}`}
                        >
                            <Trash2 size={16} />
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}
