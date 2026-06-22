import { Check } from "lucide-react";
import { useCallback, useState, type KeyboardEvent, type MouseEvent } from "react";
import type { RegisteredContestItem } from "@/client/types.gen";
import { CopyIcon } from "@/shared/icons/CopyIcon";
import { formatContestStartLabel } from "@/shared/utils/contestStartTime";

type Props = {
    contest: RegisteredContestItem;
    active: boolean;
    disabled: boolean;
    startTime?: string;
    participantCount?: number;
    onSelect: () => void;
};

async function copyText(text: string): Promise<boolean> {
    try {
        await navigator.clipboard.writeText(text);
        return true;
    } catch {
        try {
            const textarea = document.createElement("textarea");
            textarea.value = text;
            textarea.style.position = "fixed";
            textarea.style.left = "-9999px";
            document.body.appendChild(textarea);
            textarea.select();
            const ok = document.execCommand("copy");
            document.body.removeChild(textarea);
            return ok;
        } catch {
            return false;
        }
    }
}

export function ContestListRow({
    contest,
    active,
    disabled,
    startTime,
    participantCount,
    onSelect,
}: Props) {
    const [copied, setCopied] = useState(false);

    const handleCopyId = useCallback(
        async (event: MouseEvent | KeyboardEvent) => {
            event.stopPropagation();
            event.preventDefault();
            if (disabled) {
                return;
            }
            const ok = await copyText(String(contest.contest_id));
            if (ok) {
                setCopied(true);
                window.setTimeout(() => setCopied(false), 1500);
            }
        },
        [contest.contest_id, disabled],
    );

    const onIdKeyDown = useCallback(
        (event: KeyboardEvent) => {
            if (event.key === "Enter" || event.key === " ") {
                void handleCopyId(event);
            }
        },
        [handleCopyId],
    );

    const startLabel = formatContestStartLabel(startTime) ?? "—";
    const participantsLabel =
        participantCount === undefined
            ? "Участников: …"
            : `Участников: ${participantCount}`;

    return (
        <div
            className={`cf-contest-list-item${
                active ? " cf-contest-list-item--active" : ""
            }`}
        >
            <div className="cf-contest-list-item__aside">
                <span
                    className={`cf-contest-list-item__id${
                        copied ? " cf-contest-list-item__id--copied" : ""
                    }`}
                    role="button"
                    tabIndex={disabled ? -1 : 0}
                    aria-disabled={disabled || undefined}
                    onClick={handleCopyId}
                    onKeyDown={onIdKeyDown}
                    title={copied ? "Скопировано" : "Скопировать ID"}
                    aria-label={
                        copied
                            ? `ID ${contest.contest_id} скопирован`
                            : `Скопировать ID ${contest.contest_id}`
                    }
                >
                    <span className="cf-contest-list-item__id-mark" aria-hidden>
                        {copied ? (
                            <Check size={13} className="cf-contest-list-item__id-icon" />
                        ) : (
                            <CopyIcon className="cf-contest-list-item__id-copy-icon" />
                        )}
                    </span>
                    {contest.contest_id}
                </span>
                <span className="cf-contest-list-item__start">{startLabel}</span>
            </div>
            <button
                type="button"
                className="cf-contest-list-item__main"
                onClick={onSelect}
                disabled={disabled}
            >
                <span className="cf-contest-list-item__name">{contest.name}</span>
                <span className="cf-contest-list-item__participants">{participantsLabel}</span>
            </button>
        </div>
    );
}
