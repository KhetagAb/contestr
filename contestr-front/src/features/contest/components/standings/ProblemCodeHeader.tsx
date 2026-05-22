import styles from "./ProblemCodeHeader.module.css";

type Props = {
    code: string;
    url?: string;
};

export function ProblemCodeHeader({ code, url }: Props) {
    if (!url) {
        return <span>{code}</span>;
    }
    return (
        <a
            href={url}
            target="_blank"
            rel="noopener noreferrer"
            className={styles.problemLink}
            onClick={(e) => e.stopPropagation()}
            title="Открыть условие задачи (PDF)"
        >
            {code}
        </a>
    );
}
