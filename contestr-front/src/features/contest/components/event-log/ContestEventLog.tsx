import type { RegattaEvent } from "@/client";
import { tourIndexFromProblemCode } from "@/shared/utils/eventLog";
import { EventLogLine } from "./EventLogLine";
import styles from "./event-log.module.css";

const PENALTY_STRIPE_LABEL = "Вне тура";

type Props = {
    events: RegattaEvent[];
};

type TourBlock = {
    tourIndex: number;
    showTourDivider: boolean;
    afterTourEvents: RegattaEvent[];
    inTimeEvents: RegattaEvent[];
};

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

function splitTourEvents(tourEvents: RegattaEvent[]): {
    afterTourEvents: RegattaEvent[];
    inTimeEvents: RegattaEvent[];
} {
    return {
        afterTourEvents: tourEvents.filter((e) => e.solved_in_time === false),
        inTimeEvents: tourEvents.filter((e) => e.solved_in_time === true),
    };
}

function buildTourBlocks(events: RegattaEvent[]): TourBlock[] {
    if (events.length === 0) {
        return [];
    }

    const groups = groupEventsByTour(events);
    const tourIndices = [...groups.keys()]
        .filter((t) => t > 0)
        .sort((a, b) => b - a);

    const blocks: TourBlock[] = [];

    for (const tourIndex of tourIndices) {
        const tourEvents = groups.get(tourIndex);
        if (!tourEvents?.length) {
            continue;
        }

        const { afterTourEvents, inTimeEvents } = splitTourEvents(tourEvents);

        blocks.push({
            tourIndex,
            showTourDivider: tourIndex > 0,
            afterTourEvents,
            inTimeEvents,
        });
    }

    return blocks;
}

function eventKey(event: RegattaEvent, tourIndex: number): string {
    return `${tourIndex}-${event.time_sec}-${event.participant_id}-${event.problem_code}`;
}

function EventLogRows({
    tourEvents,
    tourIndex,
}: {
    tourEvents: RegattaEvent[];
    tourIndex: number;
}) {
    return (
        <>
            {tourEvents.map((event) => (
                <li key={eventKey(event, tourIndex)} className={styles.eventLogItem}>
                    <EventLogLine event={event} />
                </li>
            ))}
        </>
    );
}

function TourCard({ block }: { block: TourBlock }) {
    const { afterTourEvents, inTimeEvents, tourIndex } = block;
    const showStripe =
        afterTourEvents.length > 0 && inTimeEvents.length > 0;
    const cardClass = showStripe
        ? `${styles.eventLogTourCard} ${styles.eventLogTourCardMixed}`
        : styles.eventLogTourCard;

    return (
        <div className={cardClass}>
            {afterTourEvents.length > 0 && (
                <div className={styles.eventLogOutsideBlock}>
                    <div className={styles.eventLogOutsideMain}>
                        <ol className={styles.eventLogList}>
                            <EventLogRows
                                tourEvents={afterTourEvents}
                                tourIndex={tourIndex}
                            />
                        </ol>
                        {showStripe && (
                            <div
                                className={styles.eventLogGroupDivider}
                                role="separator"
                            />
                        )}
                    </div>
                    {showStripe && (
                        <aside
                            className={styles.eventLogPenaltyStripe}
                            aria-label="Посылки вне тура"
                        >
                            {PENALTY_STRIPE_LABEL}
                        </aside>
                    )}
                </div>
            )}

            {inTimeEvents.length > 0 && (
                <ol className={styles.eventLogList}>
                    <EventLogRows tourEvents={inTimeEvents} tourIndex={tourIndex} />
                </ol>
            )}
        </div>
    );
}

export function ContestEventLog({ events }: Props) {
    if (events.length === 0) {
        return null;
    }

    const blocks = buildTourBlocks(events);

    return (
        <section className={styles.eventLog} aria-label="История событий">
            <h3 className={styles.eventLogTitle}>История</h3>
            <div className={styles.eventLogPanel}>
                {blocks.map((block) => (
                    <div key={`block-${block.tourIndex}`} className={styles.eventLogTourBlock}>
                        {block.showTourDivider && (
                            <div
                                className={styles.eventLogTourDivider}
                                aria-label={`Тур ${block.tourIndex}`}
                            >
                                Тур {block.tourIndex}
                            </div>
                        )}
                        <TourCard block={block} />
                    </div>
                ))}
            </div>
        </section>
    );
}
