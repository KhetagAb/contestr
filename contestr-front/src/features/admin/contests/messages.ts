export function formatImportSuccessMessage(
    addedCount: number,
    skippedOtherContest: number,
): string {
    const parts: string[] = [formatAddedParticipants(addedCount)];
    if (skippedOtherContest > 0) {
        parts.push(formatSkippedOtherContest(skippedOtherContest));
    }
    return parts.join(". ");
}

export function formatImportInvalidLinesReport(invalidLines: string[]): string {
    if (invalidLines.length === 0) {
        return "";
    }
    if (invalidLines.length === 1) {
        return invalidLines[0];
    }
    return `Ошибки в ${invalidLines.length} строках: ${invalidLines.join("; ")}.`;
}

export function formatImportResultMessage(
    addedCount: number,
    skippedOtherContest: number,
    invalidLines: string[],
): string {
    const parts = [formatImportSuccessMessage(addedCount, skippedOtherContest)];
    const invalidReport = formatImportInvalidLinesReport(invalidLines);
    if (invalidReport) {
        parts.push(invalidReport);
    }
    return parts.join(" ");
}

function formatAddedParticipants(count: number): string {
    if (count === 1) {
        return "Добавлен 1 участник";
    }
    const mod10 = count % 10;
    const mod100 = count % 100;
    if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) {
        return `Добавлено ${count} участника`;
    }
    return `Добавлено ${count} участников`;
}

function formatSkippedOtherContest(count: number): string {
    if (count === 1) {
        return "Пропущен 1 участник из других контестов";
    }
    const mod10 = count % 10;
    const mod100 = count % 100;
    if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) {
        return `Пропущено ${count} участника из других контестов`;
    }
    return `Пропущено ${count} участников из других контестов`;
}
