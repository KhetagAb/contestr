const REGATTA_RULES_SEEN_KEY = "contestr:regatta-rules-seen";

export function hasSeenRegattaRules(): boolean {
    try {
        return localStorage.getItem(REGATTA_RULES_SEEN_KEY) === "1";
    } catch {
        return false;
    }
}

export function markRegattaRulesSeen(): void {
    try {
        localStorage.setItem(REGATTA_RULES_SEEN_KEY, "1");
    } catch {
        /* ignore quota / private mode */
    }
}
