type Props = {
    enabled: boolean;
    available: boolean;
    busy?: boolean;
    onChange: (enabled: boolean) => void;
};

export function AutostartToggle({ enabled, available, busy = false, onChange }: Props) {
    return (
        <label
            className="tt-autostart-toggle"
            title={available ? undefined : "На сервере не настроен фоновый timetable_sync"}
        >
            <span className="tt-autostart-toggle__label">Автозапуск</span>
            <input
                type="checkbox"
                role="switch"
                className="tt-autostart-toggle__input"
                checked={enabled}
                disabled={busy || !available}
                onChange={(e) => onChange(e.target.checked)}
            />
            <span className="tt-autostart-switch" aria-hidden />
        </label>
    );
}
