const RESERVE_PREFIX = "reserve";

/** Служебный участник-заглушка; не показываем в UI контеста. */
export function isReserveDisplayName(displayName: string | undefined | null): boolean {
    const normalized = displayName?.trim().toLowerCase();
    return !!normalized && normalized.startsWith(RESERVE_PREFIX);
}
