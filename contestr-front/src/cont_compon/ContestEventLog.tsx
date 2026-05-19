import type { RegattaEvent } from "../client";
import { tourIndexFromProblemCode } from "../utils/eventLog";
import { EventLogLine } from "./EventLogLine";
import styles from "./tables.module.css";

type Props = {
    events: RegattaEvent[];
};

type LogEntry =
    | { kind: "event"; event: RegattaEvent; key: string }
    | { kind: "tourDivider"; tourIndex: number; key: string };

function groupEventsByTour(events: RegattaEvent[]): Map<number, RegattaEvent[]> {
    const groups = new Map<number, RegattaEvent[]>();

    for (const event of events) {
        const tour = tourIndexFromProblemCode(event.problem_code) || 0;
        const bucket = groups.get(tour);
        if (bucket) {
            bucket.push(event);
        } else {
            groups.set(tour, [event]);
        }
    }

    for (const bucket of groups.values()) {
        bucket.sort((a, b) => b.time_sec - a.time_sec);
    }

    return groups;
}

function buildLogEntries(events: RegattaEvent[]): LogEntry[] {
    if (events.length === 0) {
        return [];
    }

    const groups = groupEventsByTour(events);
    const tourIndices = [...groups.keys()]
        .filter((t) => t > 0)
        .sort((a, b) => b - a);
    const entries: LogEntry[] = [];

    for (const tourIndex of tourIndices) {
        const tourEvents = groups.get(tourIndex);
        if (!tourEvents?.length) {
            continue;
        }

        if (tourIndex > 0) {
            entries.push({
                kind: "tourDivider",
                tourIndex,
                key: `tour-${tourIndex}`,
            });
        }

        for (const event of tourEvents) {
            entries.push({
                kind: "event",
                event,
                key: `${tourIndex}-${event.time_sec}-${event.participant_id}-${event.problem_code}`,
            });
        }
    }

    return entries;
}

export function ContestEventLog({ events }: Props) {
    if (events.length === 0) {
        return null;
    }

    const entries = buildLogEntries(events);

    return (
        <section className={styles.eventLog} aria-label="История событий">
            <h3 className={styles.eventLogTitle}>История</h3>
            <ol className={styles.eventLogList}>
                {entries.map((entry) =>
                    entry.kind === "tourDivider" ? (
                        <li
                            key={entry.key}
                            className={styles.eventLogTourDivider}
                            aria-label={`Тур ${entry.tourIndex}`}
                        >
                            Тур {entry.tourIndex}
                        </li>
                    ) : (
                        <li key={entry.key} className={styles.eventLogItem}>
                            <EventLogLine event={entry.event} />
                        </li>
                    ),
                )}
            </ol>
        </section>
    );
}
