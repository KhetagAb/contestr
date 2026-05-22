import {
    useEffect,
    useId,
    useRef,
    type ComponentType,
    type CSSProperties,
} from "react";
import { createPortal } from "react-dom";
import {
    CalendarRange,
    CircleCheck,
    Sparkles,
    Trophy,
    Zap,
    type LucideProps,
} from "lucide-react";
import {
    REGATTA_RULES_SECTIONS,
    type RegattaRulesSectionId,
} from "./regattaRulesContent";
import styles from "./RegattaRulesModal.module.css";

type Props = {
    open: boolean;
    onClose: () => void;
};

type SectionVisual = {
    Icon: ComponentType<LucideProps>;
    badge: string;
};

const SECTION_VISUALS: Record<RegattaRulesSectionId, SectionVisual> = {
    structure: { Icon: CalendarRange, badge: "Туры" },
    bonuses: { Icon: Zap, badge: "Очки" },
};

export function RegattaRulesModal({ open, onClose }: Props) {
    const titleId = useId();
    const closeBtnRef = useRef<HTMLButtonElement>(null);

    useEffect(() => {
        if (!open) {
            return;
        }
        closeBtnRef.current?.focus();
        const onKeyDown = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                onClose();
            }
        };
        document.addEventListener("keydown", onKeyDown);
        const prevOverflow = document.body.style.overflow;
        document.body.style.overflow = "hidden";
        return () => {
            document.removeEventListener("keydown", onKeyDown);
            document.body.style.overflow = prevOverflow;
        };
    }, [open, onClose]);

    if (!open) {
        return null;
    }

    return createPortal(
        <div
            className={styles.backdrop}
            role="presentation"
            onMouseDown={(event) => {
                if (event.target === event.currentTarget) {
                    onClose();
                }
            }}
        >
            <div
                className={styles.dialog}
                role="dialog"
                aria-modal="true"
                aria-labelledby={titleId}
            >
                <header className={styles.header}>
                    <div className={styles.headerMain}>
                        <span className={styles.headerIconWrap} aria-hidden>
                            <Trophy size={28} strokeWidth={2} />
                            <Sparkles
                                size={14}
                                strokeWidth={2}
                                className={styles.headerSparkle}
                            />
                        </span>
                        <div>
                            <h2 id={titleId} className={styles.title}>
                                Правила регаты
                            </h2>
                            <p className={styles.subtitle}>
                                Короткий гайд — как читать таблицу и набирать очки
                            </p>
                        </div>
                    </div>
                </header>
                <div className={styles.body}>
                    <div className={styles.grid}>
                        {REGATTA_RULES_SECTIONS.map((section, index) => {
                            const { Icon, badge } = SECTION_VISUALS[section.id];
                            return (
                                <section
                                    key={section.id}
                                    className={styles.card}
                                    style={
                                        {
                                            "--card-delay": `${index * 60}ms`,
                                        } as CSSProperties
                                    }
                                >
                                    <div className={styles.cardTop}>
                                        <span className={styles.cardIcon} aria-hidden>
                                            <Icon size={22} strokeWidth={2} />
                                        </span>
                                        <span className={styles.cardBadge}>{badge}</span>
                                    </div>
                                    <h3 className={styles.cardTitle}>{section.title}</h3>
                                    <ul className={styles.list}>
                                        {section.items.map((item) => (
                                            <li key={item} className={styles.listItem}>
                                                <CircleCheck
                                                    size={16}
                                                    strokeWidth={2.5}
                                                    className={styles.listIcon}
                                                    aria-hidden
                                                />
                                                <span>{item}</span>
                                            </li>
                                        ))}
                                    </ul>
                                </section>
                            );
                        })}
                    </div>
                </div>
                <footer className={styles.footer}>
                    <button
                        ref={closeBtnRef}
                        type="button"
                        className={styles.closeBtn}
                        onClick={onClose}
                    >
                        <Sparkles size={18} strokeWidth={2} aria-hidden />
                        Понятно, погнали
                    </button>
                </footer>
            </div>
        </div>,
        document.body,
    );
}
