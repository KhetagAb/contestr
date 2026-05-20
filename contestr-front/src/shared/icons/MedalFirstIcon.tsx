type Props = {
    className?: string;
};

/** sport-winner-svgrepo-com.svg — медаль «1 место». */
export function MedalFirstIcon({ className }: Props) {
    return (
        <svg
            className={className}
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden
        >
            <circle cx="12" cy="14.4" r="8" />
            <path d="M9.29,6.87l-3-3.65A1,1,0,0,1,7.08,1.6h9.84a1,1,0,0,1,.78,1.62l-3,3.65" />
            <polyline points="9.99 12.91 12.28 10.3 12.28 18.3" />
        </svg>
    );
}
