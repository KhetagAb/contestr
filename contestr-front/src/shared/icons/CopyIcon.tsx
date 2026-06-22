type Props = {
    className?: string;
};

/** Две наложенные рамки — иконка «скопировать». */
export function CopyIcon({ className }: Props) {
    return (
        <svg
            className={className}
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden
        >
            <rect x="9" y="3" width="12" height="12" rx="2.5" strokeDasharray="3.5 2.5" />
            <rect x="3" y="9" width="12" height="12" rx="2.5" />
        </svg>
    );
}
