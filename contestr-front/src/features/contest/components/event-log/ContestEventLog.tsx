import { useMemo, useState } from "react";
import type { RegattaEvent } from "@/client";
import { tourIndexFromProblemCode } from "@/shared/utils/eventLog";
import { EventLogLine } from "./EventLogLine";
import { eventEnterDelayMs, regattaEventKey } from "./eventLogKeys";
import styles from "./event-log.module.css";

function filterVisibleEvents(events: RegattaEvent[], showRejected: boolean) {
    if (showRejected) {
        return events;
    }
    return events.filter((e) => e.type !== "problem_rejected");
}

type Props = {
    events: RegattaEvent[];
};

type TourBlock = {
    tourIndex: number;
    showTourDivider: boolean;
    tourEvents: RegattaEvent[];
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

        blocks.push({
            tourIndex,
            showTourDivider: tourIndex > 0,
            tourEvents,
        });
    }

    return blocks;
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
            {tourEvents.map((event, index) => (
                <li
                    key={regattaEventKey(event, tourIndex)}
                    className={styles.eventLogItem}
                >
                    <EventLogLine
                        event={event}
                        enterDelayMs={eventEnterDelayMs(index)}
                    />
                </li>
            ))}
        </>
    );
}

function TourCard({ tourEvents, tourIndex }: { tourEvents: RegattaEvent[]; tourIndex: number }) {
    return (
        <div className={styles.eventLogTourCard}>
            <ol className={styles.eventLogList}>
                <EventLogRows tourEvents={tourEvents} tourIndex={tourIndex} />
            </ol>
        </div>
    );
}

export function ContestEventLog({ events }: Props) {
    const [showRejected, setShowRejected] = useState(false);

    const visibleEvents = useMemo(
        () => filterVisibleEvents(events, showRejected),
        [events, showRejected]
    );

    const blocks = useMemo(() => buildTourBlocks(visibleEvents), [visibleEvents]);

    if (events.length === 0) {
        return null;
    }

    return (
        <section className={styles.eventLog} aria-label="История событий">
            <div className={styles.eventLogPanel}>
                {blocks.length === 0 ? (
                    <p className={styles.eventLogEmpty}>Нет удачных событий для отображения</p>
                ) : (
                    blocks.map((block) => (
                        <div
                            key={`block-${block.tourIndex}`}
                            className={styles.eventLogTourBlock}
                        >
                            {block.showTourDivider && (
                                <div
                                    className={styles.eventLogTourDivider}
                                    aria-label={`Тур ${block.tourIndex}`}
                                >
                                    Тур {block.tourIndex}
                                </div>
                            )}
                            <TourCard
                                tourEvents={block.tourEvents}
                                tourIndex={block.tourIndex}
                            />
                        </div>
                    ))
                )}
            </div>
            <div className={styles.eventLogFooter}>
                <label className={styles.eventLogFilterToggle}>
                    <span className={styles.eventLogFilterLabel}>Неудачные посылки</span>
                    <input
                        type="checkbox"
                        role="switch"
                        className={styles.eventLogFilterInput}
                        checked={showRejected}
                        onChange={(e) => setShowRejected(e.target.checked)}
                    />
                    <span className={styles.eventLogFilterSwitch} aria-hidden />
                </label>
            </div>
        </section>
    );
}
