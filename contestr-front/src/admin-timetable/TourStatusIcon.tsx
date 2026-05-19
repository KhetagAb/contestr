import {
    Calendar,
    Check,
    CirclePlay,
    Loader2,
    type LucideProps,
} from "lucide-react";
import type { ReactNode } from "react";
import type { TourMeta } from "../client/types.gen";
import { statusLabel } from "./statusLabels";

type StatusIconProps = {
    status: TourMeta["status"];
    size?: number;
    className?: string;
};

export function TourStatusIcon({
    status,
    size = 13,
    className = "",
}: StatusIconProps) {
    const label = statusLabel(status);

    const iconProps: LucideProps = {
        size,
        className,
        "aria-hidden": true,
    };

    let icon: ReactNode;
    switch (status) {
        case "started":
            icon = <Check {...iconProps} strokeWidth={2.5} />;
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
        default:
            icon = <Calendar {...iconProps} />;
    }

    return (
        <span className="tt-status-icon-wrap" role="img" aria-label={label} title={label}>
            {icon}
        </span>
    );
}
