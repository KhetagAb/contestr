import { useEffect, useId, useRef } from "react";
import { createPortal } from "react-dom";
import { Search, X } from "lucide-react";
import { useFollowedParticipant } from "@/features/contest/follow/FollowedParticipantContext";
import { ParticipantPickerContent } from "./ParticipantPickerContent";
import styles from "./ParticipantPickerModal.module.css";

export function ParticipantPickerModal() {
    const { isParticipantPickerOpen, closeParticipantPicker } = useFollowedParticipant();
    const titleId = useId();
    const searchRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (!isParticipantPickerOpen) {
            return;
        }
        const timer = window.setTimeout(() => searchRef.current?.focus(), 0);
        const onKeyDown = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                closeParticipantPicker();
            }
        };
        document.addEventListener("keydown", onKeyDown);
        const prevOverflow = document.body.style.overflow;
        document.body.style.overflow = "hidden";
        return () => {
            window.clearTimeout(timer);
            document.removeEventListener("keydown", onKeyDown);
            document.body.style.overflow = prevOverflow;
        };
    }, [isParticipantPickerOpen, closeParticipantPicker]);

    if (!isParticipantPickerOpen) {
        return null;
    }

    return createPortal(
        <div
            className={styles.backdrop}
            role="presentation"
            onMouseDown={(event) => {
                if (event.target === event.currentTarget) {
                    closeParticipantPicker();
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
                    <span className={styles.headerIcon} aria-hidden>
                        <Search size={22} strokeWidth={2} />
                    </span>
                    <div className={styles.headerText}>
                        <h2 id={titleId} className={styles.title}>
                            Кто вы?
                        </h2>
                        <p className={styles.subtitle}>
                            Найдите себя в таблице — подсветим вашу строку и группу
                        </p>
                    </div>
                    <button
                        type="button"
                        className={styles.closeIconBtn}
                        aria-label="Закрыть"
                        onClick={closeParticipantPicker}
                    >
                        <X size={20} strokeWidth={2} aria-hidden />
                    </button>
                </header>
                <div className={styles.body}>
                    <ParticipantPickerContent searchInputRef={searchRef} />
                </div>
            </div>
        </div>,
        document.body,
    );
}
