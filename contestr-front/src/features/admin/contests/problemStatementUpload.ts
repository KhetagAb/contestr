import { isValidProblemCode } from "@/shared/utils/problemCode";

export function isPdfFile(file: File): boolean {
    return (
        file.type === "application/pdf" ||
        file.name.toLowerCase().endsWith(".pdf")
    );
}

/** Parses "1A.pdf" / "1a.pdf" → "1A", or null if not a problem code filename. */
export function codeFromPdfFilename(filename: string): string | null {
    const base = filename.replace(/\.pdf$/i, "").trim();
    const match = base.match(/^(\d+)([a-zA-Z])$/);
    if (!match) {
        return null;
    }
    const code = `${match[1]}${match[2].toUpperCase()}`;
    return isValidProblemCode(code) ? code : null;
}

export type BulkUploadItem = { code: string; file: File };

export type ClassifiedBulkFiles = {
    toUpload: BulkUploadItem[];
    skippedInvalid: string[];
    skippedUnknown: string[];
    skippedNotPdf: string[];
};

export function classifyBulkFiles(
    files: File[],
    knownCodes: Set<string>,
): ClassifiedBulkFiles {
    const toUpload: BulkUploadItem[] = [];
    const skippedInvalid: string[] = [];
    const skippedUnknown: string[] = [];
    const skippedNotPdf: string[] = [];
    const seenCodes = new Set<string>();

    for (const file of files) {
        if (!isPdfFile(file)) {
            skippedNotPdf.push(file.name);
            continue;
        }

        const code = codeFromPdfFilename(file.name);
        if (!code) {
            skippedInvalid.push(file.name);
            continue;
        }

        if (!knownCodes.has(code)) {
            skippedUnknown.push(file.name);
            continue;
        }

        if (seenCodes.has(code)) {
            skippedInvalid.push(`${file.name} (дубликат ${code})`);
            continue;
        }

        seenCodes.add(code);
        toUpload.push({ code, file });
    }

    return {
        toUpload,
        skippedInvalid,
        skippedUnknown,
        skippedNotPdf,
    };
}

export function formatBulkUploadReport(parts: {
    uploaded: number;
    total: number;
    failed: { name: string; reason: string }[];
    skippedInvalid: string[];
    skippedUnknown: string[];
    skippedNotPdf: string[];
}): string {
    const lines: string[] = [];

    if (parts.total === 0) {
        return "Нет подходящих PDF для загрузки.";
    }

    lines.push(`Загружено ${parts.uploaded} из ${parts.total}.`);

    if (parts.skippedNotPdf.length > 0) {
        lines.push(
            `Не PDF: ${parts.skippedNotPdf.join(", ")}.`,
        );
    }
    if (parts.skippedInvalid.length > 0) {
        lines.push(
            `Неверное имя (нужно 1A.pdf, 1B.pdf…): ${parts.skippedInvalid.join(", ")}.`,
        );
    }
    if (parts.skippedUnknown.length > 0) {
        lines.push(
            `Нет слота в контесте: ${parts.skippedUnknown.join(", ")}.`,
        );
    }
    if (parts.failed.length > 0) {
        lines.push(
            `Ошибки: ${parts.failed.map((f) => `${f.name} — ${f.reason}`).join("; ")}.`,
        );
    }

    return lines.join(" ");
}
