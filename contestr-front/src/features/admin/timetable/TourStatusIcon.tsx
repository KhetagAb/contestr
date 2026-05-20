import {
    Calendar,
    Check,
    CirclePlay,
    Loader2,
    type LucideProps,
} from "lucide-react";
import type { ReactNode } from "react";
import type { TimelineSegment } from "@/client/types.gen";
import { statusLabel } from "./statusLabels";
import type { TourVisualState } from "./tourVisualState";

type StatusIconProps = {
    status: TimelineSegment["status"];
    visualState?: TourVisualState;
    size?: number;
    className?: string;
};

function MovingChevronsIcon({ className = "" }: { className?: string }) {
    return (
        <span className={`tt-chevrons-icon ${className}`.trim()} aria-hidden>
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
        className,
        "aria-hidden": true,
    };

    let icon: ReactNode;
    switch (visualState ?? status) {
        case "past":
            icon = <Check {...iconProps} strokeWidth={2} />;
            break;
        case "active":
            icon = <MovingChevronsIcon className={className} />;
            break;
        case "next":
            icon = <CirclePlay {...iconProps} />;
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
        <span className="tt-status-icon-wrap" role="img" aria-label={label} title={label}>
            {icon}
        </span>
    );
}
