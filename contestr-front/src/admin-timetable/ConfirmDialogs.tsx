export function confirmStartNow(tourNumber: number) {
    return window.confirm(
        `Запустить тур ${tourNumber} сейчас?\n\nВремя этого и всех следующих туров сдвинется к текущему моменту.`,
    );
}

export function confirmEditStartedDuration(tourNumber: number) {
    return window.confirm(
        `Изменить длительность запущенного тура ${tourNumber}?\n\nСдвинется конец тура и расписание всех следующих туров.`,
    );
}
