import blobs from "@/assets/icons/circles.svg";
import crosses from "@/assets/icons/Kresty.svg";
import emblem from "@/assets/images/landing-emblem.svg";
import { useContests } from "@/shared/hooks/useContests";
import styles from "./ContestHomePage.module.css";

export default function ContestHomePage() {
    const { contests, isLoading, isError } = useContests();

    return (
        <section className={styles.page} aria-labelledby="contest-home-title">
            <div className={styles.backdrop} aria-hidden>
                <img className={styles.blobs} src={blobs} alt="" />
            </div>

            <div className={styles.decorTop} aria-hidden>
                <img className={styles.crosses} src={crosses} alt="" />
                <span className={`${styles.pixel} ${styles.pixelTop}`} />
                <span className={`${styles.pixel} ${styles.pixelMid}`} />
            </div>

            <div className={styles.decorBottom} aria-hidden>
                <img className={styles.emblem} src={emblem} alt="" />
                <span className={`${styles.pixel} ${styles.pixelLow}`} />
            </div>

            <div className={styles.content}>
                <h1 id="contest-home-title" className={styles.title}>
                    Контестер
                </h1>
                <p className={styles.lead}>
                    Контестер — образовательная платформа по олимпиадному программированию для
                    проведения соревнований.
                </p>

                <h2 className={styles.regattaHeading}>Соревнования регаты:</h2>

                {isLoading && (
                    <p className={styles.contestStatus} role="status">
                        Загрузка…
                    </p>
                )}
                {isError && (
                    <p className={styles.contestStatus} role="alert">
                        Не удалось загрузить список контестов
                    </p>
                )}
                {!isLoading && !isError && contests.length === 0 && (
                    <p className={styles.contestStatus}>Контесты не настроены</p>
                )}
                {!isLoading && !isError && contests.length > 0 && (
                    <ul className={styles.contestList}>
                        {contests.map((item) => (
                            <li key={item.contest_id}>
                                <a
                                    className={styles.contestLink}
                                    href={`/?contestId=${item.contest_id}`}
                                >
                                    {item.name}
                                </a>
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </section>
    );
}
