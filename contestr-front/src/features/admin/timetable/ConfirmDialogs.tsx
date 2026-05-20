export function confirmAdvance() {
    return window.confirm(
        "Запустить следующий слот сейчас?\n\nАктивный тур или перерыв будет укорочен до текущего момента.",
    );
}
