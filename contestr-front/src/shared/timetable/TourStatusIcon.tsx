import {
    Calendar,
    Check,
    Loader2,
    Play,
    type LucideProps,
} from "lucide-react";
import type { ReactNode } from "react";
import type { TimelineSegment } from "@/client/types.gen";
import { statusLabel } from "./statusLabels";
import type { TourVisualState } from "./tourVisualState";
import "./timetableIcons.css";

type StatusIconProps = {
    status: TimelineSegment["status"];
    visualState?: TourVisualState;
    size?: number;
    className?: string;
};

function MovingChevronsIcon({ size }: { size: number }) {
    return (
        <span
            className="tt-chevrons-icon"
            style={{ ["--tt-icon-slot-size" as string]: `${size}px` }}
            aria-hidden
        >
            <span className="tt-chevrons-icon__inner">{">>"}</span>
        </span>
    );
}

export function TourStatusIcon({
    status,
    visualState,
    size = 14,
    className = "",
}: StatusIconProps) {
    const label = statusLabel(status);

    const iconProps: LucideProps = {
        size,
        strokeWidth: 2,
        className,
        "aria-hidden": true,
    };

    let icon: ReactNode;
    switch (visualState ?? status) {
        case "past":
            icon = <Check {...iconProps} />;
            break;
        case "active":
            icon = <MovingChevronsIcon size={size} />;
            break;
        case "next":
            icon = <Play {...iconProps} />;
            break;
        case "starting":
            icon = (
                <Loader2
                    {...iconProps}
                    className={`${className} tt-status-icon--spin`.trim()}
                />
            );
            break;
        case "future":
            icon = <Calendar {...iconProps} />;
            break;
        default:
            icon = <Calendar {...iconProps} />;
    }

    return (
        <span
            className="tt-status-icon-wrap"
            style={{ ["--tt-icon-slot-size" as string]: `${size}px` }}
            role="img"
            aria-label={label}
            title={label}
        >
            {icon}
        </span>
    );
}
