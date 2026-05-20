import { useEffect, useRef, useState, type FocusEvent, type FormEvent } from "react";
import { Check, Plus, RefreshCw, Save, Trash2, Upload, X } from "lucide-react";
import { useAdminContests } from "./useAdminContests";
import "./ContestsAdminPage.css";

function PanelStatus({
    message,
    kind,
}: {
    message: string;
    kind: "error" | "success" | "info";
}) {
    if (!message) {
        return null;
    }
    if (kind === "error") {
        return (
            <p
                className="admin-login-message admin-login-message--error cf-panel-status cf-panel-status--bottom"
                role="alert"
            >
                Ошибка: {message}
            </p>
        );
    }
    return (
        <p
            className={`cf-panel-status cf-panel-status--bottom cf-panel-status--${kind}`}
            role="status"
            aria-live="polite"
        >
            {message}
        </p>
    );
}

export default function ContestsAdminPage() {
    const ac = useAdminContests();
    const [showAddForm, setShowAddForm] = useState(false);
    const [newContestId, setNewContestId] = useState("");
    const [newContestName, setNewContestName] = useState("");
    const [showImportList, setShowImportList] = useState(false);
    const [importText, setImportText] = useState("");
    const [newHandle, setNewHandle] = useState("");
    const [newName, setNewName] = useState("");
    const newHandleInputRef = useRef<HTMLInputElement>(null);

    const closeAddForm = () => {
        setShowAddForm(false);
        setNewContestId("");
        setNewContestName("");
    };

    const handleAddContest = (e: FormEvent) => {
        e.preventDefault();
        const id = Number(newContestId);
        if (!Number.isInteger(id) || id <= 0) {
            return;
        }
        void ac.addContest(id, newContestName).then(() => {
            closeAddForm();
        });
    };

    const selected = ac.contests.find((c) => c.contest_id === ac.selectedContestId);

    const closeImportList = () => {
        setShowImportList(false);
        setImportText("");
    };

    useEffect(() => {
        closeImportList();
        setNewHandle("");
        setNewName("");
    }, [ac.selectedContestId]);

    const handleAddParticipant = async () => {
        if (await ac.addParticipant(newHandle, newName)) {
            setNewHandle("");
            setNewName("");
            requestAnimationFrame(() => {
                newHandleInputRef.current?.focus();
            });
        }
    };

    const handleApplyImport = async () => {
        const result = await ac.importHandlesFromList(importText);
        if (!result.ok) {
            ac.showHandlesMessage(
                result.detail || result.invalidLines[0] || "Не удалось разобрать список",
                "error",
            );
            return;
        }

        if (result.invalidLines.length > 0) {
            ac.showHandlesMessage(`Ошибок в строках: ${result.invalidLines.length}`, "info");
        }
        setImportText("");
        setShowImportList(false);
    };

    const handleHandlesTableBlur = (event: FocusEvent<HTMLDivElement>) => {
        const wrap = event.currentTarget;
        requestAnimationFrame(() => {
            if (!wrap.contains(document.activeElement)) {
                void ac.flushHandlesDraft();
            }
        });
    };

    return (
        <section className="cf-admin-page">
            <div className="cf-admin-grid">
                <div className="cf-admin-panel cf-admin-panel--contests">
                    <h3 className="cf-admin-panel-title">Зарегистрированные контесты</h3>

                    <div className="cf-contest-list-wrap">
                        {ac.loadState === "loading" && (
                            <p className="cf-admin-hint">Загрузка…</p>
                        )}

                        {ac.loadState !== "loading" && (
                            <ul className="cf-contest-list" role="list">
                                {ac.contests.length === 0 && (
                                    <li className="cf-contest-list__empty">Список пуст</li>
                                )}
                                {ac.contests.map((c) => (
                                    <li key={c.contest_id} className="cf-contest-list__row">
                                        <button
                                            type="button"
                                            className={`cf-contest-list-item${
                                                ac.selectedContestId === c.contest_id
                                                    ? " cf-contest-list-item--active"
                                                    : ""
                                            }`}
                                            onClick={() => ac.setSelectedContestId(c.contest_id)}
                                            disabled={ac.busy}
                                        >
                                            <span className="cf-contest-list-item__name">{c.name}</span>
                                            <span className="cf-contest-list-item__meta">
                                                ID {c.contest_id} · {c.system}
                                            </span>
                                        </button>
                                        <span className="cf-contest-list__actions">
                                            <button
                                                type="button"
                                                className="cf-icon-btn"
                                                title="Обновить данные контеста"
                                                aria-label="Обновить данные контеста"
                                                disabled={ac.busy}
                                                onClick={() => void ac.refreshContest(c.contest_id)}
                                            >
                                                <RefreshCw size={16} aria-hidden />
                                            </button>
                                            <button
                                                type="button"
                                                className="cf-icon-btn"
                                                title="Удалить контест"
                                                aria-label="Удалить контест"
                                                disabled={ac.busy}
                                                onClick={() => {
                                                    if (
                                                        window.confirm(
                                                            `Удалить контест «${c.name}» (ID ${c.contest_id})?`,
                                                        )
                                                    ) {
                                                        void ac.deleteContest(c.contest_id);
                                                    }
                                                }}
                                            >
                                                <Trash2 size={18} aria-hidden />
                                            </button>
                                        </span>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>

                    <footer
                        className="cf-contest-panel-footer"
                    >
                        {showAddForm ? (
                            <form className="cf-add-contest-inline" onSubmit={handleAddContest}>
                                <input
                                    type="number"
                                    className="cf-add-contest-inline__input"
                                    min={1}
                                    required
                                    autoFocus
                                    value={newContestId}
                                    onChange={(e) => setNewContestId(e.target.value)}
                                    disabled={ac.busy}
                                    placeholder="ID"
                                    aria-label="ID контеста на Codeforces"
                                />
                                <input
                                    type="text"
                                    className="cf-add-contest-inline__input cf-add-contest-inline__input--wide"
                                    value={newContestName}
                                    onChange={(e) => setNewContestName(e.target.value)}
                                    disabled={ac.busy}
                                    placeholder="Название"
                                    aria-label="Название (необязательно)"
                                />
                                <button
                                    type="submit"
                                    className="cf-add-contest-inline__btn cf-add-contest-inline__btn--submit"
                                    title="Добавить"
                                    aria-label="Добавить контест"
                                    disabled={ac.busy}
                                >
                                    <Check size={15} aria-hidden />
                                </button>
                                <button
                                    type="button"
                                    className="cf-add-contest-inline__btn"
                                    title="Отмена"
                                    aria-label="Отмена"
                                    disabled={ac.busy}
                                    onClick={closeAddForm}
                                >
                                    <X size={15} aria-hidden />
                                </button>
                            </form>
                        ) : (
                            <button
                                type="button"
                                className="cf-secondary-btn cf-secondary-btn--compact"
                                disabled={ac.busy || ac.loadState === "loading"}
                                onClick={() => setShowAddForm(true)}
                            >
                                <Plus size={16} aria-hidden />
                                Добавить контест
                            </button>
                        )}
                    </footer>

                    <PanelStatus
                        message={ac.contestsMessage}
                        kind={ac.contestsMessageKind}
                    />
                </div>

                <div className="cf-admin-panel cf-admin-panel--handles">
                    <h3 className="cf-admin-panel-title">
                        {selected
                            ? `Зарегистрированные участники: ${selected.name}`
                            : "Зарегистрированные участники"}
                    </h3>

                    {!selected && (
                        <p className="cf-admin-hint">
                            Выберите контест слева, чтобы настроить зарегистрированных на контест
                            участников
                        </p>
                    )}

                    {selected && (
                        <>
                            {ac.handlesLoadState === "loading" && (
                                <p className="cf-admin-hint">Загрузка участников…</p>
                            )}

                            <section className="cf-scoring-settings" aria-label="Настройки начисления баллов">
                                <div className="cf-scoring-settings__row">
                                    <label>
                                        <span>Режим</span>
                                        <select
                                            value={ac.draftScoringSettings.mode}
                                            onChange={(e) =>
                                                ac.updateDraftScoringSetting(
                                                    "mode",
                                                    e.target.value === "partial" ? "partial" : "binary",
                                                )
                                            }
                                            disabled={ac.busy}
                                        >
                                            <option value="binary">ICPC</option>
                                            <option value="partial">IOI</option>
                                        </select>
                                    </label>
                                    <label>
                                        <span>Обгон в режиме ICPC</span>
                                        <select
                                            value={ac.draftScoringSettings.binary_overtake_mode}
                                            onChange={(e) =>
                                                ac.updateDraftScoringSetting(
                                                    "binary_overtake_mode",
                                                    e.target.value === "during_tour_only"
                                                        ? "during_tour_only"
                                                        : "retrospective",
                                                )
                                            }
                                            disabled={ac.busy || ac.draftScoringSettings.mode === "partial"}
                                        >
                                            <option value="retrospective">Ретроспективно</option>
                                            <option value="during_tour_only">Только во время тура</option>
                                        </select>
                                    </label>
                                </div>
                                <div className="cf-scoring-settings__row">
                                    <label>
                                        <span>Полное решение</span>
                                        <input
                                            type="number"
                                            min={0}
                                            value={ac.draftScoringSettings.full_solve_bonus}
                                            onChange={(e) =>
                                                ac.updateDraftScoringSetting(
                                                    "full_solve_bonus",
                                                    Math.max(0, Number(e.target.value) || 0),
                                                )
                                            }
                                            disabled={ac.busy}
                                        />
                                    </label>
                                    <label>
                                        <span>Во время тура</span>
                                        <input
                                            type="number"
                                            min={0}
                                            value={ac.draftScoringSettings.solve_in_time_bonus}
                                            onChange={(e) =>
                                                ac.updateDraftScoringSetting(
                                                    "solve_in_time_bonus",
                                                    Math.max(0, Number(e.target.value) || 0),
                                                )
                                            }
                                            disabled={ac.busy}
                                        />
                                    </label>
                                    <label>
                                        <span>Обгон</span>
                                        <input
                                            type="number"
                                            min={0}
                                            value={ac.draftScoringSettings.overtake_bonus}
                                            onChange={(e) =>
                                                ac.updateDraftScoringSetting(
                                                    "overtake_bonus",
                                                    Math.max(0, Number(e.target.value) || 0),
                                                )
                                            }
                                            disabled={ac.busy}
                                        />
                                    </label>
                                    <button
                                        type="button"
                                        className="cf-handles-save-btn"
                                        title="Сохранить настройки"
                                        aria-label="Сохранить настройки"
                                        onClick={() => void ac.saveContestSettings()}
                                        disabled={ac.busy || !ac.settingsDirty}
                                    >
                                        <Save size={16} aria-hidden />
                                    </button>
                                </div>
                                <div className="cf-scoring-settings__row">
                                    <label>
                                        <span>Участников в группе</span>
                                        <input
                                            type="number"
                                            min={1}
                                            value={ac.draftTourSettings.group_size}
                                            onChange={(e) =>
                                                ac.updateDraftTourSetting(
                                                    "group_size",
                                                    Math.max(1, Number(e.target.value) || 1),
                                                )
                                            }
                                            disabled={ac.busy}
                                        />
                                    </label>
                                    <label>
                                        <span>Задач в туре</span>
                                        <input
                                            type="number"
                                            min={1}
                                            value={ac.draftTourSettings.problems_per_tour}
                                            onChange={(e) =>
                                                ac.updateDraftTourSetting(
                                                    "problems_per_tour",
                                                    Math.max(1, Number(e.target.value) || 1),
                                                )
                                            }
                                            disabled={ac.busy}
                                        />
                                    </label>
                                    <label>
                                        <span>Перемешивание групп, %</span>
                                        <input
                                            type="number"
                                            min={0}
                                            max={100}
                                            value={ac.draftTourSettings.group_shuffle_percent}
                                            onChange={(e) =>
                                                ac.updateDraftTourSetting(
                                                    "group_shuffle_percent",
                                                    Math.min(100, Math.max(0, Number(e.target.value) || 0)),
                                                )
                                            }
                                            disabled={ac.busy}
                                        />
                                    </label>
                                </div>
                            </section>

                            <div className="cf-handles-table-wrap" onBlur={handleHandlesTableBlur}>
                                <table className="cf-handles-table">
                                    <colgroup>
                                        <col className="cf-handles-table__col-handle" />
                                        <col className="cf-handles-table__col-name" />
                                        <col className="cf-handles-table__col-actions" />
                                    </colgroup>
                                    <thead>
                                        <tr>
                                            <th className="cf-handles-table__head-cell">Хэндл</th>
                                            <th className="cf-handles-table__head-cell">Имя</th>
                                            <th aria-label="Действия" />
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {ac.draftHandles.map((row, index) => (
                                            <tr key={index}>
                                                <td>
                                                    <input
                                                        type="text"
                                                        value={row.handle}
                                                        onChange={(e) =>
                                                            ac.updateDraftRow(
                                                                index,
                                                                "handle",
                                                                e.target.value,
                                                            )
                                                        }
                                                        disabled={ac.busy}
                                                    />
                                                </td>
                                                <td>
                                                    <input
                                                        type="text"
                                                        value={row.name}
                                                        onChange={(e) =>
                                                            ac.updateDraftRow(
                                                                index,
                                                                "name",
                                                                e.target.value,
                                                            )
                                                        }
                                                        disabled={ac.busy}
                                                    />
                                                </td>
                                                <td className="cf-handles-table__actions">
                                                    {row.handle &&
                                                    ac.handles.some((h) => h.handle === row.handle) ? (
                                                        <button
                                                            type="button"
                                                            className="cf-icon-btn cf-icon-btn--plain"
                                                            title="Удалить из БД"
                                                            disabled={ac.busy}
                                                            onClick={() => {
                                                                if (
                                                                    window.confirm(
                                                                        `Удалить маппинг для ${row.handle}?`,
                                                                    )
                                                                ) {
                                                                    void ac.deleteHandle(row.handle);
                                                                }
                                                            }}
                                                        >
                                                            <Trash2 size={16} aria-hidden />
                                                        </button>
                                                    ) : (
                                                        <button
                                                            type="button"
                                                            className="cf-icon-btn cf-icon-btn--plain"
                                                            title="Убрать строку"
                                                            disabled={ac.busy}
                                                            onClick={() => ac.removeDraftRow(index)}
                                                        >
                                                            <Trash2 size={16} aria-hidden />
                                                        </button>
                                                    )}
                                                </td>
                                            </tr>
                                        ))}
                                        <tr className="cf-handles-table__row--add">
                                            <td>
                                                <input
                                                    ref={newHandleInputRef}
                                                    type="text"
                                                    value={newHandle}
                                                    onChange={(e) => setNewHandle(e.target.value)}
                                                    disabled={ac.busy}
                                                    aria-label="Хэндл"
                                                />
                                            </td>
                                            <td>
                                                <input
                                                    type="text"
                                                    value={newName}
                                                    onChange={(e) => setNewName(e.target.value)}
                                                    onKeyDown={(e) => {
                                                        if (e.key === "Enter") {
                                                            e.preventDefault();
                                                            void handleAddParticipant();
                                                        }
                                                    }}
                                                    disabled={ac.busy}
                                                    aria-label="Имя участника"
                                                />
                                            </td>
                                            <td className="cf-handles-table__actions">
                                                <button
                                                    type="button"
                                                    className="cf-icon-btn cf-icon-btn--plain cf-icon-btn--add"
                                                    title="Добавить участника"
                                                    aria-label="Добавить участника"
                                                    disabled={ac.busy || !newHandle.trim()}
                                                    onClick={() => void handleAddParticipant()}
                                                >
                                                    <Plus size={16} aria-hidden />
                                                </button>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>

                            {showImportList && (
                                <div className="cf-import-list">
                                    <p className="cf-admin-hint cf-import-list__hint">
                                        Вставьте регистрацию участников из Codeforces
                                    </p>
                                    <textarea
                                        className="cf-import-list__textarea"
                                        value={importText}
                                        onChange={(e) => setImportText(e.target.value)}
                                        disabled={ac.busy}
                                        rows={6}
                                        spellCheck={false}
                                        placeholder={`${selected.contest_id} | student01 | xxxxx | Иванов Иван`}
                                    />
                                    <div className="cf-import-list__actions">
                                        <button
                                            type="button"
                                            className="cf-primary-btn cf-primary-btn--compact"
                                            disabled={ac.busy || !importText.trim()}
                                            onClick={handleApplyImport}
                                        >
                                            Импортировать
                                        </button>
                                        <button
                                            type="button"
                                            className="cf-secondary-btn cf-secondary-btn--compact"
                                            disabled={ac.busy}
                                            onClick={closeImportList}
                                        >
                                            Отмена
                                        </button>
                                    </div>
                                </div>
                            )}

                            {!showImportList && (
                                <footer className="cf-handles-panel-footer">
                                    <button
                                        type="button"
                                        className="cf-secondary-btn cf-secondary-btn--compact"
                                        onClick={() => setShowImportList(true)}
                                        disabled={ac.busy}
                                    >
                                        <Upload size={16} aria-hidden />
                                        Импорт
                                    </button>
                                </footer>
                            )}

                            <PanelStatus
                                message={ac.handlesMessage}
                                kind={ac.handlesMessageKind}
                            />
                        </>
                    )}
                </div>
            </div>

        </section>
    );
}
