/** Служебный участник-заглушка; не показываем в UI контеста. */
export function isReserveDisplayName(displayName: string | undefined | null): boolean {
    return displayName?.trim().toLowerCase() === "reserve";
}
