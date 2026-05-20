import { DURATION_TEMPLATE_LABEL } from "./durationTemplate";

type Props = {
    onApplyTemplate: () => void;
    disabled?: boolean;
};

export function EmptySchedule({ onApplyTemplate, disabled }: Props) {
    return (
        <div className="tt-empty">
            <p className="tt-empty__title">Расписание не задано</p>
            <p className="tt-empty__hint">Примените шаблон или добавьте туры кнопкой ниже.</p>
            <button
                type="button"
                className="admin-icon-btn tt-empty__template"
                disabled={disabled}
                onClick={onApplyTemplate}
            >
                {DURATION_TEMPLATE_LABEL}
            </button>
        </div>
    );
}
