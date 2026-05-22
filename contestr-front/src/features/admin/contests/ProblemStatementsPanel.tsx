import {
    useEffect,
    useRef,
    useState,
    type DragEvent,
    type RefObject,
} from "react";
import { Check, RotateCw, Trash2, Upload } from "lucide-react";
import type { useAdminProblemStatements } from "./useAdminProblemStatements";
import { isPdfFile } from "./problemStatementUpload";

type Props = {
    ps: ReturnType<typeof useAdminProblemStatements>;
    selected: boolean;
};

type Row = ReturnType<typeof useAdminProblemStatements>["rows"][number];

function pickPdfFromDataTransfer(dt: DataTransfer): File | null {
    const files = Array.from(dt.files);
    return files.find(isPdfFile) ?? null;
}

function pickPdfsFromDataTransfer(dt: DataTransfer): File[] {
    return Array.from(dt.files).filter(isPdfFile);
}

type DropzoneProps = {
    className: string;
    label: string;
    hint?: string;
    busy: boolean;
    multiple?: boolean;
    fileInputRef: RefObject<HTMLInputElement | null>;
    onInvalidFile: () => void;
    onFiles: (files: File[]) => void;
};

function Dropzone({
    className,
    label,
    hint,
    busy,
    multiple = false,
    fileInputRef,
    onInvalidFile,
    onFiles,
}: DropzoneProps) {
    const [dragOver, setDragOver] = useState(false);

    const handleDragOver = (e: DragEvent<HTMLButtonElement>) => {
        e.preventDefault();
        e.stopPropagation();
        e.dataTransfer.dropEffect = busy ? "none" : "copy";
    };

    const handleDragEnter = (e: DragEvent<HTMLButtonElement>) => {
        e.preventDefault();
        e.stopPropagation();
        if (!busy && e.dataTransfer.types.includes("Files")) {
            setDragOver(true);
        }
    };

    const handleDragLeave = (e: DragEvent<HTMLButtonElement>) => {
        e.stopPropagation();
        if (!e.currentTarget.contains(e.relatedTarget as Node | null)) {
            setDragOver(false);
        }
    };

    const handleDrop = (e: DragEvent<HTMLButtonElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setDragOver(false);
        if (busy) return;

        const files = multiple
            ? pickPdfsFromDataTransfer(e.dataTransfer)
            : (() => {
                  const one = pickPdfFromDataTransfer(e.dataTransfer);
                  return one ? [one] : [];
              })();

        if (files.length === 0) {
            onInvalidFile();
            return;
        }
        onFiles(files);
    };

    const dropzoneClass = [
        className,
        dragOver ? "cf-statement-dropzone--drag-over" : "",
    ]
        .filter(Boolean)
        .join(" ");

    return (
        <button
            type="button"
            className={dropzoneClass}
            disabled={busy}
            onClick={() => fileInputRef.current?.click()}
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            onDragOver={handleDragOver}
            onDrop={handleDrop}
        >
            <Upload size={16} aria-hidden />
            <span className="cf-statement-dropzone-text">
                <span className="cf-statement-dropzone-label">{label}</span>
                {hint ? (
                    <span className="cf-statement-dropzone-hint">{hint}</span>
                ) : null}
            </span>
        </button>
    );
}

type BulkUploadZoneProps = {
    busy: boolean;
    onBulk: (files: File[]) => void;
    onInvalidFile: () => void;
};

function BulkUploadZone({ busy, onBulk, onInvalidFile }: BulkUploadZoneProps) {
    const fileInputRef = useRef<HTMLInputElement | null>(null);

    return (
        <div className="cf-statements-bulk">
            <input
                ref={fileInputRef}
                type="file"
                accept="application/pdf,.pdf"
                multiple
                className="cf-statement-file"
                disabled={busy}
                onChange={(e) => {
                    const files = Array.from(e.target.files ?? []).filter(isPdfFile);
                    if (files.length === 0 && (e.target.files?.length ?? 0) > 0) {
                        onInvalidFile();
                    } else if (files.length > 0) {
                        onBulk(files);
                    }
                    e.target.value = "";
                }}
            />
            <Dropzone
                className="cf-statement-dropzone cf-statements-bulk__zone"
                label="Массовая загрузка — перетащите PDF или выберите файлы"
                hint="Имена файлов: 1A.pdf, 1B.pdf, 2A.pdf …"
                busy={busy}
                multiple
                fileInputRef={fileInputRef}
                onInvalidFile={onInvalidFile}
                onFiles={onBulk}
            />
        </div>
    );
}

type StatementRowProps = {
    row: Row;
    busy: boolean;
    onUpload: (problemCode: string, file: File) => void | Promise<void>;
    onDelete: (problemCode: string) => void | Promise<void>;
    onInvalidFile: () => void;
};

function StatementChip({
    row,
    busy,
    onUpload,
    onDelete,
    onInvalidFile,
}: StatementRowProps) {
    const fileInputRef = useRef<HTMLInputElement | null>(null);
    const [dragOver, setDragOver] = useState(false);
    const uploaded = row.status === "uploaded";

    const handleFile = (file: File | undefined) => {
        if (!file) return;
        if (!isPdfFile(file)) {
            onInvalidFile();
            return;
        }
        void onUpload(row.problem_code, file);
    };

    const fileInput = (
        <input
            ref={fileInputRef}
            type="file"
            accept="application/pdf,.pdf"
            className="cf-statement-file"
            disabled={busy}
            onChange={(e) => {
                handleFile(e.target.files?.[0]);
                e.target.value = "";
            }}
        />
    );

    if (!uploaded) {
        const chipClass = [
            "cf-statement-chip",
            "cf-statement-chip--empty",
            dragOver ? "cf-statement-chip--drag-over" : "",
        ]
            .filter(Boolean)
            .join(" ");

        return (
            <button
                type="button"
                className={chipClass}
                disabled={busy}
                title={`Загрузить PDF для ${row.problem_code}`}
                onClick={() => fileInputRef.current?.click()}
                onDragEnter={(e) => {
                    e.preventDefault();
                    if (!busy) setDragOver(true);
                }}
                onDragLeave={(e) => {
                    if (!e.currentTarget.contains(e.relatedTarget as Node | null)) {
                        setDragOver(false);
                    }
                }}
                onDragOver={(e) => {
                    e.preventDefault();
                    e.dataTransfer.dropEffect = busy ? "none" : "copy";
                }}
                onDrop={(e) => {
                    e.preventDefault();
                    setDragOver(false);
                    if (busy) return;
                    const file = pickPdfFromDataTransfer(e.dataTransfer);
                    if (!file) {
                        onInvalidFile();
                        return;
                    }
                    handleFile(file);
                }}
            >
                {fileInput}
                <span className="cf-statement-chip__code">{row.problem_code}</span>
                <Upload size={12} aria-hidden className="cf-statement-chip__icon" />
            </button>
        );
    }

    return (
        <div className="cf-statement-chip cf-statement-chip--ok">
            {fileInput}
            {row.public_url ? (
                <a
                    href={row.public_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    download={`${row.problem_code}.pdf`}
                    className="cf-statement-chip__code cf-statement-chip__code--linked"
                    title="Скачать PDF"
                >
                    {row.problem_code}
                </a>
            ) : (
                <span className="cf-statement-chip__code cf-statement-chip__code--linked">
                    {row.problem_code}
                </span>
            )}
            <span className="cf-statement-chip__ok" aria-hidden>
                <Check size={12} />
            </span>
            <button
                type="button"
                className="cf-statement-chip__action"
                title="Заменить PDF"
                aria-label={`Заменить PDF для ${row.problem_code}`}
                disabled={busy}
                onClick={() => fileInputRef.current?.click()}
            >
                <RotateCw size={12} aria-hidden />
            </button>
            <button
                type="button"
                className="cf-statement-chip__action"
                title="Удалить PDF"
                aria-label={`Удалить PDF для ${row.problem_code}`}
                disabled={busy}
                onClick={() => {
                    if (
                        window.confirm(`Удалить условие ${row.problem_code}?`)
                    ) {
                        void onDelete(row.problem_code);
                    }
                }}
            >
                <Trash2 size={12} aria-hidden />
            </button>
        </div>
    );
}

export function ProblemStatementsPanel({ ps, selected }: Props) {
    const [validationError, setValidationError] = useState("");

    useEffect(() => {
        setValidationError("");
    }, [ps.message, selected]);

    const statusMessage = validationError || ps.message;
    const statusKind = validationError ? "error" : ps.messageKind;

    const showBulk = selected && ps.rows.length > 0 && ps.loadState !== "loading";

    return (
        <div className="cf-admin-panel cf-admin-panel--statements">
            <h3 className="cf-admin-panel-title">Условия задач (PDF)</h3>

            {!selected && (
                <p className="cf-admin-hint">
                    Выберите контест, чтобы загрузить PDF условий в Object Storage
                </p>
            )}

            {selected && (
                <>
                    {ps.loadState === "loading" && (
                        <p className="cf-admin-hint">Загрузка…</p>
                    )}
                    {ps.loadState === "error" && (
                        <p className="cf-admin-hint" role="alert">
                            Не удалось загрузить список условий
                        </p>
                    )}

                    <div className="cf-statements-layout">
                        {showBulk && (
                            <BulkUploadZone
                                busy={ps.busy}
                                onBulk={(files) => {
                                    setValidationError("");
                                    void ps.uploadBulk(files);
                                }}
                                onInvalidFile={() =>
                                    setValidationError("Нужен файл PDF")
                                }
                            />
                        )}

                        <div className="cf-statements-grid">
                            {ps.rows.map((row) => (
                                <StatementChip
                                    key={row.problem_code}
                                    row={row}
                                    busy={ps.busy}
                                    onUpload={(code, file) => {
                                        setValidationError("");
                                        return ps.uploadPdf(code, file);
                                    }}
                                    onDelete={(code) => ps.deletePdf(code)}
                                    onInvalidFile={() =>
                                        setValidationError("Нужен файл PDF")
                                    }
                                />
                            ))}
                        </div>
                    </div>

                    {ps.rows.length === 0 && ps.loadState !== "loading" && (
                        <p className="cf-admin-hint">
                            Добавьте туры в расписание, чтобы появились слоты 1A, 1B…
                        </p>
                    )}

                    {statusMessage ? (
                        <p
                            className={
                                statusKind === "error"
                                    ? "admin-login-message admin-login-message--error cf-panel-status cf-panel-status--bottom"
                                    : `cf-panel-status cf-panel-status--bottom cf-panel-status--${statusKind}`
                            }
                            role={statusKind === "error" ? "alert" : "status"}
                        >
                            {statusKind === "error"
                                ? `Ошибка: ${statusMessage}`
                                : statusMessage}
                        </p>
                    ) : null}
                </>
            )}
        </div>
    );
}
